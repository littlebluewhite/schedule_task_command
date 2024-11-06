package redis_stream

import (
	"context"
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/littlebluewhite/schedule_task_command/api"
	"github.com/redis/go-redis/v9"
)

type RedisStream struct {
	rdb        redis.UniversalClient
	streamName string
	groupName  string
	l          api.Logger
}

func NewStreamRead(rdb redis.UniversalClient, streamName string, groupName string, l api.Logger) *RedisStream {
	return &RedisStream{
		rdb:        rdb,
		streamName: streamName,
		groupName:  groupName,
		l:          l,
	}
}

func (rs *RedisStream) Start(ctx context.Context, streamComMap map[string]func(map[string]interface{}) (string, error)) {
	err := rs.streamInit(ctx)
	if err != nil {
		rs.l.Errorln("receiveStream error: ", err)
		return
	}
	for {
		rsr, err := rs.ReadGroup(ctx)
		rs.l.Infoln("TimeTemplate get stream")
		if err != nil {
			rs.l.Errorln("receiveStream error: ", err)
			if errors.Is(err, redis.ErrClosed) {
				return
			}
			continue
		}
		go func(rsr map[string]interface{}) {
			streamCom := streamComMap[rsr["command"].(string)]
			result, err := streamCom(rsr)
			if err != nil {
				rs.l.Errorln("deal stream error: ", err)
			}
			if rsr["is_wait_call_back"].(string) == "1" {
				err = rs.CallBack(ctx, rsr, result, err)
				if err != nil {
					rs.l.Errorln("call back publish error: ", err)
				}
				rs.l.Infoln("return callback success")
			}
		}(rsr)
	}
}

func (rs *RedisStream) streamInit(ctx context.Context) (err error) {
	r, e := rs.rdb.XInfoGroups(ctx, rs.streamName).Result()
	if e != nil || len(r) == 0 {
		err = rs.rdb.XGroupCreateMkStream(ctx, rs.streamName, rs.groupName, "0").Err()
		if err != nil {
			return
		}
	}
	return nil
}

func (rs *RedisStream) ReadGroup(ctx context.Context) (
	rsr map[string]interface{}, err error) {
	re, err := rs.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    rs.groupName,
		Consumer: "s1",
		Streams:  []string{rs.streamName, ">"},
		Count:    1,
		Block:    0,
		NoAck:    true,
	}).Result()
	if err != nil {
		return
	}
	if len(re) == 0 || len(re[0].Messages) == 0 {
		err = fmt.Errorf("no messages found")
		return
	}

	message := re[0].Messages[0]
	rsr = message.Values

	err = rs.rdb.XDel(ctx, rs.streamName, message.ID).Err()
	if err != nil {
		return
	}
	return
}

func (rs *RedisStream) CallBack(ctx context.Context, rsr map[string]interface{}, result string, err error) error {
	if err != nil {
		rsr["data"] = err.Error()
		rsr["status_code"] = "422"
	} else {
		rsr["data"] = result
	}
	channel := rsr["callback_channel"].(string)
	cb, _ := json.Marshal(rsr)
	err = rs.rdb.Publish(ctx, channel, cb).Err()
	return err
}
