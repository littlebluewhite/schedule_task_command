package command_server

import (
	"context"
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/littlebluewhite/schedule_task_command/api"
	"github.com/littlebluewhite/schedule_task_command/dal/model"
	"github.com/littlebluewhite/schedule_task_command/dal/query"
	"github.com/littlebluewhite/schedule_task_command/entry/e_command"
	"github.com/littlebluewhite/schedule_task_command/entry/e_module"
	"github.com/littlebluewhite/schedule_task_command/util/my_log"
	"github.com/littlebluewhite/schedule_task_command/util/redis_stream"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type CommandServer struct {
	dbs          api.Dbs
	hm           HubManager
	l            api.Logger
	c            sync.Map
	streamComMap map[string]func(rsc map[string]interface{}) (string, error)
	count        atomic.Uint64
	chs          chs
	httpClient   *http.Client
}

func NewCommandServer(dbs api.Dbs, wh HubManager) *CommandServer {
	l := my_log.NewLog("app/command_server")
	return &CommandServer{
		dbs: dbs,
		hm:  wh,
		l:   l,
		c:   sync.Map{},
		chs: chs{
			mu: new(sync.RWMutex),
			wg: new(sync.WaitGroup),
		},
		httpClient: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:        100,              // 全局最大閒置連接數
				MaxIdleConnsPerHost: 10,               // 每個主機的最大閒置連接數
				IdleConnTimeout:     60 * time.Second, // 閒置連接的超時時間
			},
		},
	}
}

func (c *CommandServer) Start(ctx context.Context, removeTime time.Duration) {
	c.initialCounter(ctx)
	// stream command initial
	c.initStreamComMap()
	c.l.Infoln("Command server started")

	c.chs.wg.Add(1)
	go func() {
		defer c.chs.wg.Done()
		c.removeFinishedCommand(ctx, removeTime)
	}()

	c.chs.wg.Add(1)
	go func() {
		defer c.chs.wg.Done()
		c.rdbSub(ctx)
	}()
	c.chs.wg.Add(1)
	go func() {
		defer c.chs.wg.Done()
		c.receiveStream(ctx)
	}()
	c.chs.wg.Add(1)
	go func() {
		ticker := time.NewTicker(time.Second * 10)
		defer ticker.Stop()
		defer c.chs.wg.Done()
		for {
			select {
			case <-ctx.Done():
				c.counterWrite()
				return
			case <-ticker.C:
				c.counterWrite()
			}
		}
	}()
}

func (c *CommandServer) Close() {
	c.chs.wg.Wait()
	c.l.Infoln("command server stop gracefully")
}

func (c *CommandServer) initialCounter(ctx context.Context) {
	qc := query.Use(c.dbs.GetSql()).Counter
	cc, err := qc.WithContext(ctx).Where(qc.Name.In("command")).First()
	if err != nil {
		cc = &model.Counter{Name: "command", Value: 0}
		e := qc.WithContext(ctx).Create(cc)
		if e != nil {
			c.l.Errorln(e)
		}
	}
	c.count.Store(uint64(cc.Value))
}

func (c *CommandServer) initStreamComMap() {
	c.streamComMap = map[string]func(rsc map[string]interface{}) (string, error){}
}

func (c *CommandServer) counterWrite() {
	ctx := context.Background()
	qc := query.Use(c.dbs.GetSql()).Counter
	_, err := qc.WithContext(ctx).Where(qc.Name.Eq("command")).Update(qc.Value, c.count.Load())
	if err != nil {
		c.l.Errorln(err)
	}
}

func (c *CommandServer) rdbSub(ctx context.Context) {
	pubsub := c.dbs.GetRdb().Subscribe(ctx, "sendCommand")
	defer func(pubsub *redis.PubSub) {
		err := pubsub.Close()
		if err != nil {
			c.l.Errorln(err)
		}
	}(pubsub)
	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				c.l.Infoln("rdbSub receive canceled")
				return
			}
			c.l.Errorln(err)
			if errors.Is(err, redis.ErrClosed) {
				return
			}
			continue
		}

		// deal with message
		go func(payload string) {
			var com e_command.Command
			if e := json.Unmarshal([]byte(payload), &com); e != nil {
				c.l.Errorln(SendToRedisErr)
				return
			}
			if _, e := c.ExecuteReturnId(ctx, com); e != nil {
				c.l.Errorln("Error executing Task:", e)
			}
		}(msg.Payload)
	}
}

func (c *CommandServer) receiveStream(ctx context.Context) {
	c.l.Infoln("----------------------------------- start command receiveStream --------------------------------")
	rs := redis_stream.NewStreamRead(c.dbs.GetRdb(), "Command", "server", c.l)
	rs.Start(ctx, c.streamComMap)
}

func (c *CommandServer) ReadMap() map[uint64]e_command.Command {
	result := make(map[uint64]e_command.Command)
	c.c.Range(func(key, value interface{}) bool {
		id, ok1 := key.(uint64)
		com, ok2 := value.(e_command.Command)
		if ok1 && ok2 {
			result[id] = com
		}
		return true
	})
	return result
}

func (c *CommandServer) ReadOne(id uint64) (e_command.Command, error) {
	if com, ok := c.c.Load(id); ok {
		return com.(e_command.Command), nil
	} else {
		return e_command.Command{}, errors.New(fmt.Sprintf("cannot find id: %d", id))
	}
}

func (c *CommandServer) DeleteOne(id uint64) {
	c.c.Delete(id)
}

func (c *CommandServer) GetList() []e_command.Command {
	result := make([]e_command.Command, 0, 1000)
	c.c.Range(func(key, value interface{}) bool {
		_, ok1 := key.(uint64)
		com, ok2 := value.(e_command.Command)
		if ok1 && ok2 {
			result = append(result, com)
		}
		return true
	})
	return result
}

func (c *CommandServer) ExecuteReturnId(ctx context.Context, com e_command.Command) (id uint64, err error) {
	// pass the variables
	com = c.getVariables(com)
	if com.Message != nil {
		err = com.Message
		c.l.Errorln(err)
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
		c.l.Errorln(com.Message)
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

func (c *CommandServer) TestExecute(ctx context.Context, com e_command.Command) e_command.Command {
	// pass the variables
	com = c.getVariables(com)
	// add initial variables
	if com.Variables == nil {
		com.Variables = make(map[string]string)
	}
	if com.Message != nil {
		c.l.Errorln(com.Message)
		return com
	}
	from := time.Now()
	com.From = from
	ch := make(chan e_command.Command)
	go func() {
		ch <- c.doCommandNoRecord(ctx, com)
	}()
	com = <-ch
	return com
}

func (c *CommandServer) doCommand(ctx context.Context, com e_command.Command) e_command.Command {
	ctx, cancel := context.WithTimeout(ctx,
		time.Duration(com.CommandData.Timeout)*time.Millisecond)
	defer cancel()

	// add default variables to command
	com = addDefaultVariables(com)

	com.Status = e_command.Process
	com.CancelFunc = cancel
	// write command
	c.writeCommand(com)

	com = c.requestProtocol(ctx, com)
	now := time.Now()
	com.To = &now

	// write client message
	com2, err := c.ReadOne(com.ID)
	if err != nil {
		c.l.Errorln(err)
	}
	com.ClientMessage = com2.ClientMessage

	// write command
	c.writeCommand(com)

	// write to history in influxdb
	go func() {
		c.writeToHistory(com)
	}()

	// publish to all channel
	c.publishContainer(context.Background(), com)
	return com
}

func (c *CommandServer) doCommandNoRecord(ctx context.Context, com e_command.Command) e_command.Command {
	ctx, cancel := context.WithTimeout(ctx,
		time.Duration(com.CommandData.Timeout)*time.Millisecond)

	// add default variables to command
	com = addDefaultVariables(com)

	com.Status = e_command.Process
	com.CancelFunc = cancel

	com = c.requestProtocol(ctx, com)
	now := time.Now()
	com.To = &now

	return com
}

func (c *CommandServer) CancelCommand(id uint64, message string) error {
	com, err := c.ReadOne(id)
	if err != nil {
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

func (c *CommandServer) removeFinishedCommand(ctx context.Context, removeTime time.Duration) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			c.l.Infoln("removeFinishedCommand goroutine exiting")
			return
		case <-ticker.C:
			c.cleanupCommand(removeTime)
		}
	}
}

func (c *CommandServer) cleanupCommand(removeTime time.Duration) {
	now := time.Now()
	c.c.Range(func(key, value interface{}) bool {
		id, ok1 := key.(uint64)
		com, ok2 := value.(e_command.Command)
		if ok1 && ok2 {
			if com.To != nil && com.Status != e_command.Process && com.To.Add(removeTime).Before(now) {
				c.c.Delete(id)
			}
		}
		return true
	})
}

func (c *CommandServer) writeToHistory(com e_command.Command) {
	tp := e_command.ToPub(com)
	jCom, err := json.Marshal(tp)
	if err != nil {
		c.l.Errorf("err: %v, ToPub: %+v", err, tp)
	}
	templateId := fmt.Sprintf("%d", com.TemplateId)
	p := influxdb2.NewPoint("command_history",
		map[string]string{
			"id":                  strconv.FormatUint(com.ID, 10),
			"command_template_id": templateId,
			"status":              com.Status.String()},
		map[string]interface{}{"data": jCom},
		com.From,
	)

	// write to influxdb
	c.dbs.GetIdb().Writer().WritePoint(p)
}

func (c *CommandServer) ReadFromHistory(id, comTemplateId, start, stop, status string) (hc []e_command.CommandPub, err error) {
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
	comIDValue := ""
	if id != "" {
		comIDValue = fmt.Sprintf(`|> filter(fn: (r) => r.id == "%s")`, id)
	}
	stmt := fmt.Sprintf(`from(bucket:"schedule")
|> range(start: %s%s)
|> filter(fn: (r) => r._measurement == "command_history")
|> filter(fn: (r) => r._field == "data")
%s
%s
%s
`, start, stopValue, comTemplateValue, statusValue, comIDValue)
	result, err := c.dbs.GetIdb().Querier().Query(ctx, stmt)
	if err == nil {
		for result.Next() {
			var com e_command.CommandPub
			v := result.Record().Value()
			vString, ok := v.(string)
			if !ok {
				c.l.Errorf("value: %v is not string", v)
				continue
			}
			if e := json.Unmarshal([]byte(vString), &com); e != nil {
				c.l.Errorln(e)
				continue
			}
			hc = append(hc, com)
		}
	} else {
		return
	}
	// send empty []
	if hc == nil {
		hc = make([]e_command.CommandPub, 0)
	}
	return
}

func (c *CommandServer) publishContainer(
	ctx context.Context,
	com e_command.Command) {
	//go func() {
	//	_ = c.rdbPub(ctx, com)
	//}()
	go func() {
		c.sendWebsocket(com)
	}()
}

func (c *CommandServer) rdbPub(ctx context.Context, com e_command.Command) (err error) {
	cb, _ := json.Marshal(e_command.ToPub(com))
	err = c.dbs.GetRdb().Publish(ctx, "CommandRec", cb).Err()
	if err != nil {
		c.l.Errorln("redis publish error: ", err)
		return
	}
	return
}

func (c *CommandServer) sendWebsocket(com e_command.Command) {
	cb, _ := json.Marshal(e_command.ToPub(com))
	if c.hm != nil {
		c.hm.Broadcast(e_module.Command, cb)
	}
}

func (c *CommandServer) writeCommand(com e_command.Command) {
	c.c.Store(com.ID, com)
}

func (c *CommandServer) getVariables(com e_command.Command) e_command.Command {
	v := make(map[string]string)
	// template 有變數
	if com.CommandData.Variable != nil && string(com.CommandData.Variable) != "null" {
		if err := json.Unmarshal(com.CommandData.Variable, &v); err != nil {
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
