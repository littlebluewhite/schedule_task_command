package influxdb

import (
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/littlebluewhite/schedule_task_command/util/config"
)

type Influx struct {
	client  influxdb2.Client
	writer  api.WriteAPI
	querier api.QueryAPI
}

type Logger interface {
	Infoln(args ...interface{})
	Infof(s string, args ...interface{})
	Errorln(args ...interface{})
	Errorf(s string, args ...interface{})
	Warnln(args ...interface{})
	Warnf(s string, args ...interface{})
}

func NewInfluxdb(influxConfig config.InfluxdbConfig, log Logger) *Influx {
	dsn := fmt.Sprintf("http://%s:%s", influxConfig.Host, influxConfig.Port)
	writeOptions := influxdb2.DefaultOptions().SetBatchSize(500).SetFlushInterval(10000)
	client := influxdb2.NewClientWithOptions(dsn, influxConfig.Token, writeOptions)
	writeAPI := client.WriteAPI(influxConfig.Org, influxConfig.Bucket)
	queryAPI := client.QueryAPI(influxConfig.Org)

	// handle error
	go func() {
		for err := range writeAPI.Errors() {
			log.Errorln(fmt.Printf("Write error: %s", err.Error()))
		}
	}()

	return &Influx{
		client,
		writeAPI,
		queryAPI,
	}
}

func (i *Influx) Close() {
	i.client.Close()
}

func (i *Influx) Writer() api.WriteAPI {
	return i.writer
}

func (i *Influx) Querier() api.QueryAPI {
	return i.querier
}
