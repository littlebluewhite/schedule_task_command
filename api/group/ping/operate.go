package ping

import (
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
	"schedule_task_command/app/dbs"
)

type Operate struct {
	db    *gorm.DB
	cache *cache.Cache
}

func NewOperate(dbs dbs.Dbs) *Operate {
	return &Operate{
		db:    dbs.GetSql(),
		cache: dbs.GetCache(),
	}
}
