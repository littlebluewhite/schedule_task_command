package task_template

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
	"schedule_task_command/entry/e_task"
	"schedule_task_command/entry/e_task_template"
	"schedule_task_command/util"
	"time"
)

type Operate struct {
	db    *gorm.DB
	cache *cache.Cache
	rdb   *redis.Client
	taskS api.TaskServer
}

func NewOperate(dbs dbs.Dbs, taskS api.TaskServer) *Operate {
	o := &Operate{
		db:    dbs.GetSql(),
		cache: dbs.GetCache(),
		rdb:   dbs.GetRdb(),
		taskS: taskS,
	}
	err := o.ReloadCache()
	if err != nil {
		panic("initial time template Operate error")
	}
	return o
}

func (o *Operate) getCacheMap() map[int]model.TaskTemplate {
	var cacheMap map[int]model.TaskTemplate
	if x, found := o.cache.Get("taskTemplates"); found {
		cacheMap = x.(map[int]model.TaskTemplate)
	} else {
		return make(map[int]model.TaskTemplate)
	}
	return cacheMap
}

func (o *Operate) setCacheMap(cacheMap map[int]model.TaskTemplate) {
	o.cache.Set("taskTemplates", cacheMap, cache.NoExpiration)
}

func (o *Operate) listDB() ([]*model.TaskTemplate, error) {
	t := query.Use(o.db).TaskTemplate
	ctx := context.Background()
	tt, err := t.WithContext(ctx).Preload(field.Associations).Preload(t.Stages.CommandTemplate).Preload(
		t.Stages.CommandTemplate.Monitor).Preload(t.Stages.CommandTemplate.Monitor.MConditions).Find()
	if err != nil {
		return nil, err
	}
	return tt, nil
}

func (o *Operate) listCache() ([]model.TaskTemplate, error) {
	var tt []model.TaskTemplate
	cacheMap := o.getCacheMap()
	fmt.Println(cacheMap)
	for _, value := range cacheMap {
		tt = append(tt, value)
	}
	return tt, nil
}

func (o *Operate) List() ([]model.TaskTemplate, error) {
	return o.listCache()
}

func (o *Operate) ReloadCache() (e error) {
	tt, err := o.listDB()
	if err != nil {
		e = err
		return
	}
	cacheMap := make(map[int]model.TaskTemplate)
	for i := 0; i < len(tt); i++ {
		entry := tt[i]
		cacheMap[int(entry.ID)] = *entry
	}
	o.setCacheMap(cacheMap)
	return
}

func (o *Operate) findDB(ctx context.Context, q *query.Query, ids []int32) ([]*model.TaskTemplate, error) {
	t := q.TaskTemplate
	TaskTemplates, err := t.WithContext(ctx).Preload(field.Associations).Preload(t.Stages.CommandTemplate).Where(t.ID.In(ids...)).Find()
	if err != nil {
		return nil, err
	}
	return TaskTemplates, nil
}

func (o *Operate) findCache(ids []int32) ([]model.TaskTemplate, error) {
	tt := make([]model.TaskTemplate, 0, len(ids))
	var cacheMap map[int]model.TaskTemplate
	if x, found := o.cache.Get("taskTemplates"); found {
		cacheMap = x.(map[int]model.TaskTemplate)
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

func (o *Operate) Find(ids []int32) ([]model.TaskTemplate, error) {
	return o.findCache(ids)
}

func (o *Operate) Create(c []*e_task_template.TaskTemplateCreate) ([]model.TaskTemplate, error) {
	q := query.Use(o.db)
	ctx := context.Background()
	cacheMap := o.getCacheMap()
	taskTemplates := e_task_template.CreateConvert(c)
	result := make([]model.TaskTemplate, 0, len(taskTemplates))
	err := q.Transaction(func(tx *query.Query) error {
		if err := tx.TaskTemplate.WithContext(ctx).CreateInBatches(taskTemplates, 100); err != nil {
			return err
		}
		for _, t := range taskTemplates {
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

func (o *Operate) Update(u []*e_task_template.TaskTemplateUpdate) error {
	cacheMap := o.getCacheMap()
	tt, e := e_task_template.UpdateConvert(cacheMap, u)
	if e != nil {
		return e
	}
	ids := make([]int32, 0, len(tt))
	q := query.Use(o.db)
	ctx := context.Background()
	err := q.Transaction(func(tx *query.Query) error {
		for _, item := range tt {
			ids = append(ids, item.ID)
			t := util.StructToMap(item)
			sUpdate := make([]map[string]interface{}, 0, 10)
			sCreate := make([]*model.TaskStage, 0, 10)
			sDelete := make([]int32, 0, 10)
			for _, stage := range item.Stages {
				s := stage
				switch {
				case stage.ID < 0:
					sDelete = append(sDelete, -s.ID)
				case stage.ID == 0:
					sCreate = append(sCreate, &s)
				case stage.ID > 0:
					sUpdate = append(sUpdate, util.StructToMap(s))
				}
			}
			t["stages"] = sUpdate
			delete(t, "stages")
			delete(t, "updated_at")
			delete(t, "created_at")
			if _, err := tx.TaskTemplate.WithContext(ctx).Where(tx.TaskTemplate.ID.Eq(
				item.ID)).Updates(t); err != nil {
				return err
			}
			for _, si := range sUpdate {
				delete(si, "command_template")
				if _, err := tx.TaskStage.WithContext(ctx).Where(tx.TaskStage.ID.Eq((si["id"]).(int32))).Updates(si); err != nil {
					return err
				}
			}
			if err := tx.TaskStage.WithContext(ctx).Create(sCreate...); err != nil {
				return err
			}
			tts := make([]*model.TaskTemplateStage, 0, len(sCreate))
			for _, ts := range sCreate {
				tts = append(tts, &model.TaskTemplateStage{
					TaskStageID: ts.ID, TaskTemplateID: item.ID})
			}
			if err := tx.TaskTemplateStage.WithContext(ctx).Create(tts...); err != nil {
				return err
			}
			if _, err := tx.TaskStage.WithContext(ctx).Where(tx.TaskStage.ID.In(
				sDelete...)).Delete(); err != nil {
				return err
			}
		}
		newTaskTemplate, err := o.findDB(ctx, tx, ids)
		if err != nil {
			return err
		}
		for _, t := range newTaskTemplate {
			cacheMap[int(t.ID)] = *t
		}
		o.setCacheMap(cacheMap)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (o *Operate) Delete(ids []int32) error {
	cacheMap := o.getCacheMap()
	sIds := make([]int32, 0, 20)
	for _, i := range ids {
		tt, ok := cacheMap[int(i)]
		if !ok {
			return errors.New(fmt.Sprintf("id: %d not found", i))
		}
		for _, s := range tt.Stages {
			sIds = append(sIds, s.ID)
		}
	}
	q := query.Use(o.db)
	ctx := context.Background()
	err := q.Transaction(func(tx *query.Query) error {
		if _, err := tx.TaskTemplate.WithContext(ctx).Where(
			tx.TaskTemplate.ID.In(ids...)).Delete(); err != nil {
			return err
		}
		if _, err := tx.TaskStage.WithContext(ctx).Where(
			tx.TaskStage.ID.In(sIds...)).Delete(); err != nil {
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

func (o *Operate) Execute(ctx context.Context, st e_task_template.SendTaskTemplate) (taskId string, err error) {
	task := o.generateTask(st)
	taskId, err = o.taskS.ExecuteReturnId(ctx, task)
	return
}

func (o *Operate) generateTask(st e_task_template.SendTaskTemplate) (task e_task.Task) {
	task = e_task.Task{
		TemplateID:     st.TemplateId,
		TriggerFrom:    st.TriggerFrom,
		TriggerAccount: st.TriggerAccount,
		Token:          st.Token,
	}
	ttList, err := o.findCache([]int32{int32(st.TemplateId)})
	if err != nil {
		task.Status = e_task.Status{TStatus: e_task.Failure}
		task.Message = &CannotFindTemplate
		return
	}
	n := time.Now()
	tt := e_task_template.Format(ttList)[0]
	task.TaskId = fmt.Sprintf("%v_%v_%v", st.TemplateId, tt.Name, n.UnixMicro())
	task.From = n
	task.Template = tt
	return
}
