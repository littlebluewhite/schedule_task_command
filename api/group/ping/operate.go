package ping

import (
	"github.com/littlebluewhite/schedule_task_command/api"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
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
