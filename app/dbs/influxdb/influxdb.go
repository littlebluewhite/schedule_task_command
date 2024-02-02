package influxdb

import (
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"schedule_task_command/util/config"
)

type Influx struct {
	client  influxdb2.Client
	writer  api.WriteAPIBlocking
	querier api.QueryAPI
}

func NewInfluxdb(influxConfig config.InfluxdbConfig) *Influx {
	dsn := fmt.Sprintf("http://%s:%s", influxConfig.Host, influxConfig.Port)
	client := influxdb2.NewClient(dsn, influxConfig.Token)
	writeAPI := client.WriteAPIBlocking(influxConfig.Org, influxConfig.Bucket)
	queryAPI := client.QueryAPI(influxConfig.Org)
	return &Influx{
		client,
		writeAPI,
		queryAPI,
	}
}

func (i *Influx) Close() {
	i.client.Close()
}

func (i *Influx) Writer() api.WriteAPIBlocking {
	return i.writer
}

func (i *Influx) Querier() api.QueryAPI {
	return i.querier
}
