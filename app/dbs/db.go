package dbs

import (
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"schedule_task_command/app/dbs/influxdb"
	"schedule_task_command/app/dbs/rdb"
	"schedule_task_command/app/dbs/sql"
	"schedule_task_command/util/logFile"
	"time"
)

type Dbs interface {
	initSql(log logFile.LogFile)
	initCache()
	initRdb(log logFile.LogFile)
	initIdb(log logFile.LogFile)
	GetSql() *gorm.DB
	GetCache() *cache.Cache
	GetRdb() *redis.Client
	GetIdb() HistoryDB
}

type HistoryDB interface {
	Close()
	Writer() api.WriteAPIBlocking
	Querier() api.QueryAPI
}

type dbs struct {
	Sql   *gorm.DB
	Cache *cache.Cache
	Rdb   *redis.Client
	Idb   HistoryDB
}

func NewDbs(log logFile.LogFile, IsTest bool) Dbs {
	d := &dbs{}
	if IsTest {
		d.initTestSql(log)
	} else {
		d.initSql(log)
	}
	d.initCache()
	d.initRdb(log)
	d.initIdb(log)
	return d
}

// DB start
func (d *dbs) initTestSql(log logFile.LogFile) {
	s, err := sql.NewDB("mySQL", "DB_test.log", "db_test")
	if err != nil {
		log.Error().Println("DB Connection failed")
		panic(err)
	} else {
		log.Info().Println("DB Connection successful")
	}
	d.Sql = s
}

// DB start
func (d *dbs) initSql(log logFile.LogFile) {
	s, err := sql.NewDB("mySQL", "DB.log", "db")
	if err != nil {
		log.Error().Println("DB Connection failed")
		panic(err)
	} else {
		log.Info().Println("DB Connection successful")
	}
	d.Sql = s
}

func (d *dbs) initCache() {
	d.Cache = cache.New(5*time.Minute, 10*time.Minute)
}

func (d *dbs) initRdb(log logFile.LogFile) {
	d.Rdb = rdb.NewRedis("redis")
	log.Info().Println("Redis Connection successful")
}

func (d *dbs) initIdb(log logFile.LogFile) {
	d.Idb = influxdb.NewInfluxdb("influxdb")
	log.Info().Println("InfluxDB Connection successful")
}

func (d *dbs) GetSql() *gorm.DB {
	return d.Sql
}

func (d *dbs) GetCache() *cache.Cache {
	return d.Cache
}

func (d *dbs) GetRdb() *redis.Client {
	return d.Rdb
}

func (d *dbs) GetIdb() HistoryDB {
	return d.Idb
}
