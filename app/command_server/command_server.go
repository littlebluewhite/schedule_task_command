package command_server

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"schedule_task_command/app/dbs"
	"schedule_task_command/dal/model"
	"schedule_task_command/entry/e_command"
	"schedule_task_command/entry/e_command_template"
	"schedule_task_command/util/logFile"
	"sync"
	"time"
)

type CommandServer struct {
	dbs dbs.Dbs
	l   logFile.LogFile
	c   map[string]e_command.Command
	chs chs
}

func NewCommandServer(dbs dbs.Dbs) *CommandServer {
	l := logFile.NewLogFile("app", "command_server")
	c := make(map[string]e_command.Command)
	rec := make(chan e_command.Command)
	mu := new(sync.RWMutex)
	return &CommandServer{
		dbs: dbs,
		l:   l,
		c:   c,
		chs: chs{
			rec: rec,
			mu:  mu,
		},
	}
}

func (c *CommandServer) Start(ctx context.Context, removeTime time.Duration) {
	c.l.Info().Println("Command server started")
	defer c.l.Info().Println("Command server stopped")
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func(wg *sync.WaitGroup) {
		c.removeFinishedCommand(ctx, removeTime)
		wg.Done()
	}(wg)
	go func(wg *sync.WaitGroup) {
		c.rdbSub(ctx)
		wg.Done()
	}(wg)
	wg.Wait()
}

func (c *CommandServer) rdbSub(ctx context.Context) {
	pubsub := c.dbs.GetRdb().Subscribe(ctx, "sendCommand")
	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}
		b := []byte(msg.Payload)
		var s e_command.SendCommand
		err = json.Unmarshal(b, &s)
		if err != nil {
			c.l.Error().Println(SendToRedisErr)
			continue
		}
		s.TriggerFrom = append(s.TriggerFrom, "redis channel")
		_, err = c.execute(s)
		if err != nil {
			c.l.Error().Println("Error executing Command")
		}
	}
}

func (c *CommandServer) execute(sc e_command.SendCommand) (commandId string, err error) {
	ctx := context.Background()
	com := c.generateCommand(sc)
	// publish to redis
	_ = c.rdbPub(com)
	if err = com.Message; err != nil {
		c.l.Error().Println(err)
		return
	}
	go func() {
		c.doCommand(ctx, com)
	}()
	commandId = com.CommandId
	return
}

func (c *CommandServer) Execute(ctx context.Context, sc e_command.SendCommand) (com e_command.Command) {
	com = c.generateCommand(sc)
	// publish to redis
	_ = c.rdbPub(com)
	if com.Message != nil {
		c.l.Error().Println(com.Message)
		return
	}
	ch := make(chan e_command.Command)
	go func() {
		ch <- c.doCommand(ctx, com)
	}()
	com = <-ch
	return
}

func (c *CommandServer) doCommand(ctx context.Context, com e_command.Command) e_command.Command {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(com.Template.Timeout)*time.Millisecond)
	defer cancel()

	com.Status = e_command.Process
	com.CancelFunc = cancel

	// write command
	c.writeCommand(com)

	com = c.requestProtocol(ctx, com)
	now := time.Now()
	com.To = &now

	// write command
	c.writeCommand(com)

	// write to history in influxdb
	c.writeToHistory(com)
	// send to redis channel
	if e := c.rdbPub(com); e != nil {
		panic(e)
	}
	return com
}

func (c *CommandServer) generateCommand(sc e_command.SendCommand) (com e_command.Command) {
	cache := c.dbs.GetCache()
	var cacheMap map[int]model.CommandTemplate
	if x, found := cache.Get("commandTemplates"); found {
		cacheMap = x.(map[int]model.CommandTemplate)
	}
	ct, ok := cacheMap[sc.TemplateId]
	if !ok {
		com = e_command.Command{Token: sc.Token, Message: &CannotFindTemplate, Status: e_command.Failure}
		return
	}
	from := time.Now()
	commandId := fmt.Sprintf("%v_%v_%v_%v", sc.TemplateId, ct.Name, ct.Protocol, from.UnixMicro())
	com = e_command.Command{
		CommandId:      commandId,
		Token:          sc.Token,
		From:           from,
		TriggerFrom:    sc.TriggerFrom,
		TriggerAccount: sc.TriggerAccount,
		TemplateId:     sc.TemplateId,
		Template:       e_command_template.Format([]model.CommandTemplate{ct})[0],
	}
	return
}

func (c *CommandServer) CancelCommand(commandId string) error {
	c.chs.mu.RLock()
	com, ok := c.c[commandId]
	c.chs.mu.RUnlock()
	if !ok {
		return CommandNotFind
	}
	if com.Status != e_command.Process {
		return CommandCannotCancel
	} else {
		com.CancelFunc()
	}
	return nil
}

func (c *CommandServer) ShowCommandList() (cs []e_command.Command) {
	c.chs.mu.RLock()
	defer c.chs.mu.RUnlock()
	for _, item := range c.c {
		cs = append(cs, item)
	}
	return
}

func (c *CommandServer) removeFinishedCommand(ctx context.Context, s time.Duration) {
Loop1:
	for {
		select {
		case <-ctx.Done():
			break Loop1
		default:
			c.chs.mu.Lock()
			now := time.Now()
			for cId, item := range c.c {
				if item.Status != e_command.Process && item.To.Add(s).After(now) {
					delete(c.c, cId)
				}
			}
			time.Sleep(s)
		}
	}
}

func (c *CommandServer) writeToHistory(com e_command.Command) {
	ctx := context.Background()
	p := influxdb2.NewPoint("command_history",
		map[string]string{"command_id": com.CommandId, "status": com.Status.String()},
		map[string]interface{}{"data": com},
		com.From,
	)
	if err := c.dbs.GetIdb().Writer().WritePoint(ctx, p); err != nil {
		panic(err)
	}
}

func (c *CommandServer) ReadFromHistory(commandId, start, stop, status string) (hc []e_command.Command) {
	ctx := context.Background()
	stopValue := ""
	if stop != "" {
		stopValue = fmt.Sprintf(", stop: %s", stop)
	}
	statusValue := ""
	if status != "" {
		statusValue = fmt.Sprintf(`|> filter(fn: (r) => r.status == "%s")`, status)
	}
	stmt := fmt.Sprintf(`from(bucket:"schedule")
|> range(start: %s%s)
|> filter(fn: (r) => r._measurement == "command_history")
|> filter(fn: (r) => r.command_id == "%s")
|> filter(fn: (r) => r._field == "data")
%s
`, start, stopValue, commandId, statusValue)
	result, err := c.dbs.GetIdb().Querier().Query(ctx, stmt)
	if err == nil {
		for result.Next() {
			var com e_command.Command
			v := result.Record().Value()
			if e := json.Unmarshal([]byte(v.(string)), &com); e != nil {
				panic(e)
			}
			hc = append(hc, com)
		}
	} else {
		panic(err)
	}
	return
}

func (c *CommandServer) rdbPub(com e_command.Command) (e error) {
	ctx := context.Background()
	cb, _ := json.Marshal(e_command.ToPub(com))
	e = c.dbs.GetRdb().Publish(ctx, "CommandRec", cb).Err()
	if e != nil {
		c.l.Error().Println("redis publish error")
		return
	}
	return
}

func (c *CommandServer) writeCommand(com e_command.Command) {
	c.chs.mu.Lock()
	defer c.chs.mu.Unlock()
	c.c[com.CommandId] = com
}
