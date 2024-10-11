package influxdb

import (
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	api2 "schedule_task_command/api"
	"schedule_task_command/util/config"
)

type Influx struct {
	client  influxdb2.Client
	writer  api.WriteAPI
	querier api.QueryAPI
}

func NewInfluxdb(influxConfig config.InfluxdbConfig, log api2.Logger) *Influx {
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
