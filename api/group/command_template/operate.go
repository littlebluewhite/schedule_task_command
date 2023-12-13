package command_template

import (
	"context"
	"errors"
	"fmt"
	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"schedule_task_command/api"
	"schedule_task_command/app/dbs"
	"schedule_task_command/dal/model"
	"schedule_task_command/dal/query"
	"schedule_task_command/entry/e_command"
	"schedule_task_command/entry/e_command_template"
	"schedule_task_command/util"
)

type Operate struct {
	db            *gorm.DB
	cache         *cache.Cache
	rdb           *redis.Client
	commandS      api.CommandServer
	taskTemplateO taskTemplateOperate
}

type taskTemplateOperate interface {
	ReloadCache(ctx context.Context, q *query.Query, ids []int32) error
}

func NewOperate(dbs dbs.Dbs, commandS api.CommandServer, taskTemplateO taskTemplateOperate) *Operate {
	o := &Operate{
		db:            dbs.GetSql(),
		cache:         dbs.GetCache(),
		rdb:           dbs.GetRdb(),
		commandS:      commandS,
		taskTemplateO: taskTemplateO,
	}
	err := o.ReloadCache()
	if err != nil {
		panic("initial command template Operate error")
	}
	return o
}

func (o *Operate) getCacheMap() map[int]model.CommandTemplate {
	var cacheMap map[int]model.CommandTemplate
	if x, found := o.cache.Get("commandTemplates"); found {
		cacheMap = x.(map[int]model.CommandTemplate)
	} else {
		return make(map[int]model.CommandTemplate)
	}
	return cacheMap
}

func (o *Operate) setCacheMap(cacheMap map[int]model.CommandTemplate) {
	o.cache.Set("commandTemplates", cacheMap, cache.NoExpiration)
}

func (o *Operate) listDB() ([]*model.CommandTemplate, error) {
	c := query.Use(o.db).CommandTemplate
	ctx := context.Background()
	ct, err := c.WithContext(ctx).Preload(field.Associations).Preload(c.Monitor.MConditions).Find()
	if err != nil {
		return nil, err
	}
	return ct, nil
}

func (o *Operate) listCache() ([]model.CommandTemplate, error) {
	var tt []model.CommandTemplate
	cacheMap := o.getCacheMap()
	for _, value := range cacheMap {
		tt = append(tt, value)
	}
	return tt, nil
}

func (o *Operate) List() ([]model.CommandTemplate, error) {
	return o.listCache()
}

func (o *Operate) ReloadCache() (e error) {
	tt, err := o.listDB()
	if err != nil {
		e = err
		return
	}
	cacheMap := make(map[int]model.CommandTemplate)
	for i := 0; i < len(tt); i++ {
		entry := tt[i]
		cacheMap[int(entry.ID)] = *entry
	}
	o.setCacheMap(cacheMap)
	return
}
func (o *Operate) findDB(ctx context.Context, q *query.Query, ids []int32) ([]*model.CommandTemplate, error) {
	c := q.CommandTemplate
	CommandTemplates, err := c.WithContext(ctx).Preload(field.Associations).Preload(c.Monitor.MConditions).Where(c.ID.In(ids...)).Find()
	if err != nil {
		return nil, err
	}
	return CommandTemplates, nil
}

func (o *Operate) findCache(ids []int32) ([]model.CommandTemplate, error) {
	tt := make([]model.CommandTemplate, 0, len(ids))
	var cacheMap map[int]model.CommandTemplate
	if x, found := o.cache.Get("commandTemplates"); found {
		cacheMap = x.(map[int]model.CommandTemplate)
	} else {
		return nil, errors.New("cache error")
	}
	for _, id := range ids {
		t, ok := cacheMap[int(id)]
		if !ok {
			return nil, fmt.Errorf("id: %v not found", id)
		}
		tt = append(tt, t)
	}
	return tt, nil
}

func (o *Operate) Find(ids []int32) ([]model.CommandTemplate, error) {
	return o.findCache(ids)
}

func (o *Operate) Create(c []*e_command_template.CommandTemplateCreate) ([]model.CommandTemplate, error) {
	q := query.Use(o.db)
	ctx := context.Background()
	cacheMap := o.getCacheMap()
	commandTemplates := e_command_template.CreateConvert(c)
	result := make([]model.CommandTemplate, 0, len(commandTemplates))
	err := q.Transaction(func(tx *query.Query) error {
		if err := tx.CommandTemplate.WithContext(ctx).CreateInBatches(commandTemplates, 100); err != nil {
			return err
		}
		for _, t := range commandTemplates {
			cacheMap[int(t.ID)] = *t
			result = append(result, *t)
		}
		o.setCacheMap(cacheMap)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (o *Operate) Update(u []*e_command_template.CommandTemplateUpdate) error {
	cacheMap := o.getCacheMap()
	ct, e := e_command_template.UpdateConvert(cacheMap, u)
	if e != nil {
		return e
	}
	ids := make([]int32, 0, len(ct))
	q := query.Use(o.db)
	ctx := context.Background()
	err := q.Transaction(func(tx *query.Query) error {
		for _, item := range ct {
			ids = append(ids, item.ID)
			if item.Monitor != nil {
				mcUpdate := make([]map[string]interface{}, 0, 10)
				mcCreate := make([]*model.MCondition, 0, 10)
				mcDelete := make([]int32, 0, 10)
				MonitorId := item.Monitor.ID
				for _, mCondition := range item.Monitor.MConditions {
					mc := mCondition
					switch {
					case mc.ID < 0:
						mcDelete = append(mcDelete, -mc.ID)
					case mc.ID == 0:
						mc.MonitorID = &MonitorId
						mcCreate = append(mcCreate, &mc)
					case mc.ID > 0:
						mc.MonitorID = &MonitorId
						mcUpdate = append(mcUpdate, util.StructToMap(mc))
					}
				}
				m := util.StructToMap(item.Monitor)
				delete(m, "m_conditions")
				if _, err := tx.Monitor.WithContext(ctx).Where(tx.Monitor.ID.Eq(
					item.Monitor.ID)).Updates(m); err != nil {
					return err
				}
				for _, mci := range mcUpdate {
					if _, err := tx.MCondition.WithContext(ctx).Where(tx.MCondition.ID.Eq(
						(mci["id"]).(int32))).Updates(mci); err != nil {
						return err
					}
				}
				if err := tx.MCondition.WithContext(ctx).CreateInBatches(mcCreate, 100); err != nil {
					return err
				}
				if _, err := tx.MCondition.WithContext(ctx).Where(tx.MCondition.ID.In(mcDelete...)).Delete(); err != nil {
					return err
				}
			}
			prUpdate := make([]map[string]interface{}, 0, 10)
			prCreate := make([]*model.ParserReturn, 0, 10)
			prDelete := make([]int32, 0, 10)
			for _, parserReturn := range item.ParserReturn {
				pr := parserReturn
				switch {
				case pr.ID < 0:
					prDelete = append(prDelete, -pr.ID)
				case pr.ID == 0:
					pr.CommandTemplateID = item.ID
					prCreate = append(prCreate, &pr)
				case pr.ID > 0:
					pr.CommandTemplateID = item.ID
					prUpdate = append(prUpdate, util.StructToMap(pr))
				}
			}
			for _, prs := range prUpdate {
				if _, err := tx.ParserReturn.WithContext(ctx).Where(tx.ParserReturn.ID.Eq(
					(prs["id"]).(int32))).Updates(prs); err != nil {
					return err
				}
			}
			if err := tx.ParserReturn.WithContext(ctx).CreateInBatches(prCreate, 100); err != nil {
				return err
			}
			if _, err := tx.ParserReturn.WithContext(ctx).Where(tx.ParserReturn.ID.In(prDelete...)).Delete(); err != nil {
				return err
			}

			if item.Http != nil {
				h := util.StructToMap(item.Http)
				if _, err := tx.HTTPSCommand.WithContext(ctx).Where(tx.HTTPSCommand.ID.Eq(
					item.Http.ID)).Updates(h); err != nil {
					return err
				}
			}
			if item.Mqtt != nil {
				mq := util.StructToMap(item.Monitor)
				if _, err := tx.MqttCommand.WithContext(ctx).Where(tx.MqttCommand.ID.Eq(
					item.Mqtt.ID)).Updates(mq); err != nil {
					return err
				}
			}
			if item.Websocket != nil {
				w := util.StructToMap(item.Websocket)
				if _, err := tx.WebsocketCommand.WithContext(ctx).Where(tx.WebsocketCommand.ID.Eq(
					item.Websocket.ID)).Updates(w); err != nil {
					return err
				}
			}
			if item.Redis != nil {
				r := util.StructToMap(item.Redis)
				if _, err := tx.RedisCommand.WithContext(ctx).Where(tx.RedisCommand.ID.Eq(
					item.Redis.ID)).Updates(r); err != nil {
					return err
				}
			}
			t := util.StructToMap(item)
			delete(t, "http")
			delete(t, "mqtt")
			delete(t, "websocket")
			delete(t, "redis")
			delete(t, "monitor")
			delete(t, "parser_return")
			delete(t, "created_at")
			delete(t, "updated_at")
			if _, err := tx.CommandTemplate.WithContext(ctx).Where(tx.CommandTemplate.ID.Eq(
				item.ID)).Updates(t); err != nil {
				return err
			}
		}
		newCommandTemplate, err := o.findDB(ctx, tx, ids)
		if err != nil {
			return err
		}
		for _, c := range newCommandTemplate {
			cacheMap[int(c.ID)] = *c
		}
		o.setCacheMap(cacheMap)
		// reload task template
		err = o.reloadTaskTemplate(ctx, tx, ids)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (o *Operate) Delete(ids []int32) error {
	cacheMap := o.getCacheMap()
	q := query.Use(o.db)
	ctx := context.Background()
	err := q.Transaction(func(tx *query.Query) error {
		if _, err := tx.CommandTemplate.WithContext(ctx).Where(
			tx.CommandTemplate.ID.In(ids...)).Delete(); err != nil {
			return err
		}
		for _, id := range ids {
			delete(cacheMap, int(id))
		}
		o.setCacheMap(cacheMap)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (o *Operate) Execute(ctx context.Context, sc e_command_template.SendCommandTemplate) (id uint64, err error) {
	c := o.generateCommand(sc)
	id, err = o.commandS.ExecuteReturnId(ctx, c)
	return
}

func (o *Operate) generateCommand(sc e_command_template.SendCommandTemplate) (c e_command.Command) {
	c = e_command.Command{
		TemplateId:     sc.TemplateId,
		TriggerFrom:    sc.TriggerFrom,
		TriggerAccount: sc.TriggerAccount,
		Token:          sc.Token,
		Variables:      sc.Variables,
	}
	cList, err := o.findCache([]int32{sc.TemplateId})
	if err != nil {
		c.Status = e_command.Failure
		c.Message = &e_command_template.CannotFindTemplate
		return
	}
	ct := e_command_template.Format(cList)[0]
	c.CommandData = ct

	return
}

func (o *Operate) reloadTaskTemplate(ctx context.Context, q *query.Query, commandTemplateIds []int32) error {
	qts := q.StageItem
	qt := q.TaskTemplateStage
	stages, err := qts.WithContext(ctx).Select(qts.ID).Where(qts.CommandTemplateID.In(commandTemplateIds...)).Find()
	if err != nil {
		return err
	}
	stageIds := make([]int32, 0, len(stages))
	for _, stage := range stages {
		stageIds = append(stageIds, stage.ID)
	}
	tts, err := qt.WithContext(ctx).Select(qt.TaskTemplateID).Where(qt.StageItemID.In(stageIds...)).Find()
	if err != nil {
		return err
	}
	templateIds := make([]int32, 0, len(tts))
	for _, tt := range tts {
		templateIds = append(templateIds, tt.TaskTemplateID)
	}
	err = o.taskTemplateO.ReloadCache(ctx, q, templateIds)
	if err != nil {
		return err
	}
	return nil
}
