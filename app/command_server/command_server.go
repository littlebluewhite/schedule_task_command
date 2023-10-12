package command_server

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"schedule_task_command/app/dbs"
	"schedule_task_command/entry/e_command"
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
		var com e_command.Command
		err = json.Unmarshal(b, &com)
		if err != nil {
			c.l.Error().Println(SendToRedisErr)
			continue
		}
		com.TriggerFrom = append(com.TriggerFrom, "redis channel")
		_, err = c.ExecuteReturnId(ctx, com)
		if err != nil {
			c.l.Error().Println("Error executing Command")
		}
	}
}

func (c *CommandServer) ReadMap() map[string]e_command.Command {
	c.chs.mu.RLock()
	defer c.chs.mu.RUnlock()
	return c.c
}

func (c *CommandServer) GetList() []e_command.Command {
	cl := make([]e_command.Command, 0, len(c.c))
	m := c.ReadMap()
	for _, v := range m {
		cl = append(cl, v)
	}
	return cl
}

func (c *CommandServer) ExecuteReturnId(ctx context.Context, com e_command.Command) (commandId string, err error) {
	// pass the variables
	com = c.getVariables(com)
	// publish to redis
	_ = c.rdbPub(com)
	if com.Message != nil {
		err = com.Message
		c.l.Error().Println(err)
		return
	}
	from := time.Now()
	com.From = from
	commandId = fmt.Sprintf("%v_%v_%v_%v",
		com.TemplateId, com.Template.Name, com.Template.Protocol, from.UnixMicro())
	com.CommandId = commandId
	go func() {
		c.doCommand(ctx, com)
	}()
	return
}

func (c *CommandServer) ExecuteWait(ctx context.Context, com e_command.Command) e_command.Command {
	// pass the variables
	com = c.getVariables(com)
	// add initial variables
	if com.Variables == nil {
		com.Variables = make(map[string]string)
	}
	// publish to redis
	_ = c.rdbPub(com)
	if com.Message != nil {
		c.l.Error().Println(com.Message)
		return com
	}
	from := time.Now()
	com.From = from
	com.CommandId = fmt.Sprintf("%v_%v_%v_%v",
		com.TemplateId, com.Template.Name, com.Template.Protocol, from.UnixMicro())
	ch := make(chan e_command.Command)
	go func() {
		ch <- c.doCommand(ctx, com)
	}()
	com = <-ch
	return com
}

func (c *CommandServer) doCommand(ctx context.Context, com e_command.Command) e_command.Command {
	ctx, cancel := context.WithTimeout(ctx,
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

func (c *CommandServer) CancelCommand(commandId, message string) error {
	m := c.ReadMap()
	com, ok := m[commandId]
	if !ok {
		return CommandNotFind
	}
	if com.Status != e_command.Process {
		return CommandCannotCancel
	}
	com.AccountMessage = message
	c.writeCommand(com)
	com.CancelFunc()
	return nil
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
			c.chs.mu.Unlock()
			time.Sleep(s)
		}
	}
}

func (c *CommandServer) writeToHistory(com e_command.Command) {
	ctx := context.Background()
	tp := e_command.ToPub(com)
	fmt.Println(string(tp.Template.Http.Body))
	jCom, err := json.Marshal(tp)
	if err != nil {
		panic(err)
	}
	templateId := fmt.Sprintf("%d", com.TemplateId)
	p := influxdb2.NewPoint("command_history",
		map[string]string{"command_template_id": templateId, "status": com.Status.String()},
		map[string]interface{}{"data": jCom},
		com.From,
	)
	if err := c.dbs.GetIdb().Writer().WritePoint(ctx, p); err != nil {
		panic(err)
	}
}

func (c *CommandServer) ReadFromHistory(comTemplateId, start, stop, status string) (hc []e_command.CommandPub, err error) {
	ctx := context.Background()
	stopValue := ""
	if stop != "" {
		stopValue = fmt.Sprintf(", stop: %s", stop)
	}
	statusValue := ""
	if status != "" {
		statusValue = fmt.Sprintf(`|> filter(fn: (r) => r.status == "%s")`, status)
	}
	comTemplateValue := ""
	if comTemplateId != "" {
		comTemplateValue = fmt.Sprintf(`|> filter(fn: (r) => r.task_template_id == "%s")`, comTemplateId)
	}
	stmt := fmt.Sprintf(`from(bucket:"schedule")
|> range(start: %s%s)
|> filter(fn: (r) => r._measurement == "command_history")
|> filter(fn: (r) => r._field == "data")
%s
%s
`, start, stopValue, comTemplateValue, statusValue)
	result, err := c.dbs.GetIdb().Querier().Query(ctx, stmt)
	if err == nil {
		for result.Next() {
			var com e_command.CommandPub
			v := result.Record().Value()
			if e := json.Unmarshal([]byte(v.(string)), &com); e != nil {
				panic(e)
			}
			hc = append(hc, com)
		}
	} else {
		return
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

func (c *CommandServer) getVariables(com e_command.Command) e_command.Command {
	if com.Variables == nil {
		v := make(map[string]string)
		com.Variables = v
		// template have variables
		if com.Template.Variable != nil {
			e := json.Unmarshal(com.Template.Variable, &v)
			if e != nil {
				com.Message = &CommandTemplateVariable
				return com
			}
		}
	}
	return com
}
