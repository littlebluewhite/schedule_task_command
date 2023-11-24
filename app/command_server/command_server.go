package command_server

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"schedule_task_command/app/dbs"
	"schedule_task_command/dal/model"
	"schedule_task_command/dal/query"
	"schedule_task_command/entry/e_command"
	"schedule_task_command/util/logFile"
	"schedule_task_command/util/redis_stream"
	"sync"
	"sync/atomic"
	"time"
)

type CommandServer struct {
	dbs          dbs.Dbs
	wm           websocketManager
	l            logFile.LogFile
	c            map[uint64]e_command.Command
	streamComMap map[string]func(rsc map[string]interface{}) (string, error)
	count        atomic.Uint64
	chs          chs
}

func NewCommandServer(dbs dbs.Dbs, wm websocketManager) *CommandServer {
	l := logFile.NewLogFile("app", "command_server")
	c := make(map[uint64]e_command.Command)
	mu := new(sync.RWMutex)
	return &CommandServer{
		dbs: dbs,
		wm:  wm,
		l:   l,
		c:   c,
		chs: chs{
			mu: mu,
		},
	}
}

func (c *CommandServer) Start(ctx context.Context, removeTime time.Duration) {
	c.initialCounter(ctx)
	// stream command initial
	c.initStreamComMap()
	c.l.Info().Println("Command server started")
	go func() {
		c.removeFinishedCommand(ctx, removeTime)
	}()
	go func() {
		c.rdbSub(ctx)
	}()
	go func() {
		c.receiveStream(ctx)
	}()
	go func() {
		_ = <-ctx.Done()
		c.stopCounter()
		c.l.Info().Println("command server stop gracefully")
	}()
}

func (c *CommandServer) initialCounter(ctx context.Context) {
	qc := query.Use(c.dbs.GetSql()).Counter
	cc, err := qc.WithContext(ctx).Where(qc.Name.In("command")).First()
	if err != nil {
		cc = &model.Counter{Name: "command", Value: 0}
		e := qc.WithContext(ctx).Create(cc)
		if e != nil {
			c.l.Error().Println(e)
		}
	}
	c.count.Store(uint64(cc.Value))
}

func (c *CommandServer) initStreamComMap() {
	c.streamComMap = map[string]func(rsc map[string]interface{}) (string, error){}
}

func (c *CommandServer) stopCounter() {
	ctx := context.Background()
	qc := query.Use(c.dbs.GetSql()).Counter
	_, err := qc.WithContext(ctx).Where(qc.Name.Eq("command")).Update(qc.Value, c.count.Load())
	if err != nil {
		c.l.Error().Println(err)
	}
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

func (c *CommandServer) receiveStream(ctx context.Context) {
	c.l.Info().Println("----------------------------------- start command receiveStream --------------------------------")
	rs := redis_stream.NewStreamRead(c.dbs.GetRdb(), "Command", "server", c.l)
	rs.Start(ctx, c.streamComMap)
}

func (c *CommandServer) ReadMap() map[uint64]e_command.Command {
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

func (c *CommandServer) ExecuteReturnId(ctx context.Context, com e_command.Command) (id uint64, err error) {
	// pass the variables
	com = c.getVariables(com)
	if com.Message != nil {
		err = com.Message
		c.l.Error().Println(err)
		return
	}
	from := time.Now()
	com.From = from
	com.ID = c.count.Add(1)
	id = com.ID
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
	if com.Message != nil {
		c.l.Error().Println(com.Message)
		return com
	}
	from := time.Now()
	com.From = from
	com.ID = c.count.Add(1)
	ch := make(chan e_command.Command)
	go func() {
		ch <- c.doCommand(ctx, com)
	}()
	com = <-ch
	return com
}

func (c *CommandServer) doCommand(ctx context.Context, com e_command.Command) e_command.Command {
	ctx, cancel := context.WithTimeout(ctx,
		time.Duration(com.CommandData.Timeout)*time.Millisecond)

	com.Status = e_command.Process
	com.CancelFunc = cancel
	// write command
	c.writeCommand(com)

	com = c.requestProtocol(ctx, com)
	now := time.Now()
	com.To = &now

	// write client message
	com.ClientMessage = c.ReadMap()[com.ID].ClientMessage

	// write command
	c.writeCommand(com)

	// write to history in influxdb
	c.writeToHistory(com)
	// publish to all channel
	c.publishContainer(context.Background(), com)
	return com
}

func (c *CommandServer) CancelCommand(id uint64, message string) error {
	m := c.ReadMap()
	com, ok := m[id]
	if !ok {
		return CommandNotFind
	}
	if com.Status != e_command.Process {
		return CommandCannotCancel
	}
	com.ClientMessage = message
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
				if item.Status != e_command.Process && item.To.Add(s).Before(now) {
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
		comTemplateValue = fmt.Sprintf(`|> filter(fn: (r) => r.command_template_id == "%s")`, comTemplateId)
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

func (c *CommandServer) publishContainer(ctx context.Context, com e_command.Command) {
	go func() {
		_ = c.rdbPub(ctx, com)
	}()
	go func() {
		c.sendWebsocket(com)
	}()
}

func (c *CommandServer) rdbPub(ctx context.Context, com e_command.Command) (err error) {
	cb, _ := json.Marshal(e_command.ToPub(com))
	err = c.dbs.GetRdb().Publish(ctx, "CommandRec", cb).Err()
	if err != nil {
		c.l.Error().Println("redis publish error: ", err)
		return
	}
	return
}

func (c *CommandServer) sendWebsocket(com e_command.Command) {
	cb, _ := json.Marshal(e_command.ToPub(com))
	c.wm.Broadcast(1, cb)
}

func (c *CommandServer) writeCommand(com e_command.Command) {
	c.chs.mu.Lock()
	defer c.chs.mu.Unlock()
	c.c[com.ID] = com
}

func (c *CommandServer) getVariables(com e_command.Command) e_command.Command {
	v := make(map[string]string)
	// template 有變數
	if com.CommandData.Variable != nil {
		e := json.Unmarshal(com.CommandData.Variable, &v)
		if e != nil {
			com.Message = &CommandTemplateVariable
			return com
		}
	}
	// 傳進來有變數
	if com.Variables != nil {
		for key, value := range com.Variables {
			v[key] = value
		}
	}
	com.Variables = v
	return com
}
