package redis_stream

import (
	"context"
	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	"schedule_task_command/util/logFile"
)

type RedisStream struct {
	rdb        *redis.Client
	streamName string
	groupName  string
	l          logFile.LogFile
}

func NewStreamRead(rdb *redis.Client, streamName string, groupName string, l logFile.LogFile) *RedisStream {
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
		rs.l.Error().Println("receiveStream error: ", err)
		return
	}
	for {
		rsr, err := rs.ReadGroup(ctx)
		rs.l.Info().Println("TimeTemplate get stream")
		if err != nil {
			rs.l.Error().Println("receiveStream error: ", err)
			continue
		}
		go func(rsr map[string]interface{}) {
			streamCom := streamComMap[rsr["command"].(string)]
			result, err := streamCom(rsr)
			if err != nil {
				rs.l.Error().Println("deal stream error: ", err)
			}
			if rsr["is_wait_call_back"].(string) == "1" {
				err = rs.CallBack(ctx, rsr, result, err)
				if err != nil {
					rs.l.Error().Println("call back publish error: ", err)
				}
				rs.l.Info().Println("return callback success")
			}
		}(rsr)
	}
}

func (rs *RedisStream) streamInit(ctx context.Context) (err error) {
	if e := rs.rdb.XInfoGroups(ctx, rs.streamName).Err(); e != nil {
		r, _ := rs.rdb.XAdd(ctx, &redis.XAddArgs{
			Stream: rs.streamName,
			Values: map[string]interface{}{
				"message": "initial",
			},
		}).Result()
		err = rs.rdb.XGroupCreate(ctx, rs.streamName, rs.groupName, r).Err()
		if err != nil {
			return
		}
		err = rs.rdb.XDel(ctx, rs.streamName, r).Err()
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
	rsr = re[0].Messages[0].Values
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
