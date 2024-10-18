package ping

import (
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
	"schedule_task_command/api"
)

type Operate struct {
	db    *gorm.DB
	cache *cache.Cache
}

func NewOperate(dbs api.Dbs) *Operate {
	return &Operate{
		db:    dbs.GetSql(),
		cache: dbs.GetCache(),
	}
}
