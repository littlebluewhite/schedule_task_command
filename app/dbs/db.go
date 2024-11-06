package dbs

import (
	"github.com/littlebluewhite/schedule_task_command/api"
	"github.com/littlebluewhite/schedule_task_command/app/dbs/influxdb"
	"github.com/littlebluewhite/schedule_task_command/app/dbs/rdb"
	"github.com/littlebluewhite/schedule_task_command/app/dbs/sql"
	"github.com/littlebluewhite/schedule_task_command/util/config"
	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"time"
)

type Dbs struct {
	Sql   *gorm.DB
	Cache *cache.Cache
	Rdb   redis.UniversalClient
	Idb   *influxdb.Influx
}

func NewDbs(log api.Logger, IsTest bool, config config.Config) *Dbs {
	d := &Dbs{}
	if IsTest {
		d.initTestSql(log, config.TestSQL)
	} else {
		d.initSql(log, config.SQL)
	}
	d.initCache()
	d.initRdb(log, config.Redis)
	d.initIdb(log, config.Influxdb)
	return d
}

// DB start
func (d *Dbs) initTestSql(log api.Logger, Config config.SQLConfig) {
	s, err := sql.NewDB("mySQL", "DB_test.my_log", Config)
	if err != nil {
		log.Errorln("DB Connection failed")
		panic(err)
	} else {
		log.Infoln("DB Connection successful")
	}
	d.Sql = s
}

// DB start
func (d *Dbs) initSql(log api.Logger, Config config.SQLConfig) {
	s, err := sql.NewDB("mySQL", "DB.my_log", Config)
	if err != nil {
		log.Errorln("DB Connection failed")
		panic(err)
	} else {
		log.Infoln("DB Connection successful")
	}
	d.Sql = s
}

func (d *Dbs) initCache() {
	d.Cache = cache.New(5*time.Minute, 10*time.Minute)
}

func (d *Dbs) initRdb(log api.Logger, Config config.RedisConfig) {
	d.Rdb = rdb.NewClient(Config)
	log.Infoln("Redis Connection successful")
}

func (d *Dbs) initIdb(log api.Logger, Config config.InfluxdbConfig) {
	d.Idb = influxdb.NewInfluxdb(Config, log)
	log.Infoln("InfluxDB Connection successful")
}

func (d *Dbs) GetSql() *gorm.DB {
	return d.Sql
}

func (d *Dbs) GetCache() *cache.Cache {
	return d.Cache
}

func (d *Dbs) GetRdb() redis.UniversalClient {
	return d.Rdb
}

func (d *Dbs) GetIdb() *influxdb.Influx {
	return d.Idb
}

func (d *Dbs) Close() {
	_ = d.Rdb.Close()
	d.Idb.Close()
}
