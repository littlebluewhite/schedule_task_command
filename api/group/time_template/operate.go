package time_template

import (
	"context"
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"schedule_task_command/api"
	"schedule_task_command/app/dbs"
	"schedule_task_command/dal/model"
	"schedule_task_command/dal/query"
	"schedule_task_command/entry/e_time"
	"schedule_task_command/entry/e_time_template"
	"schedule_task_command/util"
	"strconv"
	"time"
)

type Operate struct {
	db    *gorm.DB
	cache *cache.Cache
	rdb   *redis.Client
	timeS api.TimeServer
}

func NewOperate(dbs dbs.Dbs, timeS api.TimeServer) *Operate {
	o := &Operate{
		db:    dbs.GetSql(),
		cache: dbs.GetCache(),
		rdb:   dbs.GetRdb(),
		timeS: timeS,
	}
	err := o.ReloadCache()
	if err != nil {
		panic("initial time template Operate error")
	}
	return o
}

func (o *Operate) getCacheMap() map[int]model.TimeTemplate {
	var cacheMap map[int]model.TimeTemplate
	if x, found := o.cache.Get("timeTemplates"); found {
		cacheMap = x.(map[int]model.TimeTemplate)
	} else {
		return make(map[int]model.TimeTemplate)
	}
	return cacheMap
}

func (o *Operate) setCacheMap(cacheMap map[int]model.TimeTemplate) {
	o.cache.Set("timeTemplates", cacheMap, cache.NoExpiration)
}

func (o *Operate) listDB() ([]*model.TimeTemplate, error) {
	t := query.Use(o.db).TimeTemplate
	ctx := context.Background()
	timeTemplates, err := t.WithContext(ctx).Preload(field.Associations).Find()
	if err != nil {
		return nil, err
	}
	return timeTemplates, nil
}

func (o *Operate) listCache() ([]model.TimeTemplate, error) {
	var tt []model.TimeTemplate
	cacheMap := o.getCacheMap()
	for _, value := range cacheMap {
		tt = append(tt, value)
	}
	return tt, nil
}

func (o *Operate) List() ([]model.TimeTemplate, error) {
	return o.listCache()
}

func (o *Operate) ReloadCache() (e error) {
	tt, err := o.listDB()
	if err != nil {
		e = err
		return
	}
	cacheMap := make(map[int]model.TimeTemplate)
	for i := 0; i < len(tt); i++ {
		entry := tt[i]
		cacheMap[int(entry.ID)] = *entry
	}
	o.setCacheMap(cacheMap)
	return
}

func (o *Operate) findDB(ctx context.Context, q *query.Query, ids []int32) ([]*model.TimeTemplate, error) {
	t := q.TimeTemplate
	timeTemplates, err := t.WithContext(ctx).Preload(field.Associations).Where(t.ID.In(ids...)).Find()
	if err != nil {
		return nil, err
	}
	return timeTemplates, nil
}

func (o *Operate) findCache(ids []int32) ([]model.TimeTemplate, error) {
	tt := make([]model.TimeTemplate, 0, len(ids))
	var cacheMap map[int]model.TimeTemplate
	if x, found := o.cache.Get("timeTemplates"); found {
		cacheMap = x.(map[int]model.TimeTemplate)
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

func (o *Operate) Find(ids []int32) ([]model.TimeTemplate, error) {
	return o.findCache(ids)
}

func (o *Operate) Create(c []*e_time_template.TimeTemplateCreate) ([]model.TimeTemplate, error) {
	q := query.Use(o.db)
	ctx := context.Background()
	cacheMap := o.getCacheMap()
	timeTemplates := e_time_template.CreateConvert(c)
	result := make([]model.TimeTemplate, 0, len(timeTemplates))
	err := q.Transaction(func(tx *query.Query) error {
		if err := tx.TimeTemplate.WithContext(ctx).CreateInBatches(timeTemplates, 100); err != nil {
			return err
		}
		for _, t := range timeTemplates {
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

func (o *Operate) Update(u []*e_time_template.TimeTemplateUpdate) error {
	cacheMap := o.getCacheMap()
	tt, e := e_time_template.UpdateConvert(cacheMap, u)
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
			td := t["time_data"].(map[string]interface{})
			delete(t, "time_data")
			delete(t, "updated_at")
			delete(t, "created_at")
			delete(td, "id")
			if _, err := tx.TimeTemplate.WithContext(ctx).Where(tx.TimeTemplate.ID.Eq(item.ID)).Updates(
				t); err != nil {
				return err
			}
			if _, err := tx.TimeDatum.WithContext(ctx).Where(tx.TimeDatum.ID.Eq(item.TimeDataID)).Updates(
				td); err != nil {
				return err
			}
		}
		// update cache
		newTimeTemplate, err := o.findDB(ctx, tx, ids)
		if err != nil {
			return err
		}
		for _, t := range newTimeTemplate {
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
	tdId := make([]int32, 0, len(ids))
	for _, i := range ids {
		t, ok := cacheMap[int(i)]
		if !ok {
			return errors.New(fmt.Sprintf("id: %d not found", i))
		}
		tdId = append(tdId, t.TimeDataID)
	}
	q := query.Use(o.db)
	ctx := context.Background()
	err := q.Transaction(func(tx *query.Query) error {
		if _, err := tx.TimeTemplate.WithContext(ctx).Where(
			tx.TimeTemplate.ID.In(ids...)).Delete(); err != nil {
			return err
		}
		if _, err := tx.TimeDatum.WithContext(ctx).Where(
			tx.TimeDatum.ID.In(tdId...)).Delete(); err != nil {
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

func (o *Operate) CheckTime(id int, c CheckTime) (isTime bool, err error) {
	st := SendTime{
		TemplateId:     id,
		TriggerFrom:    c.TriggerFrom,
		TriggerAccount: c.TriggerAccount,
		Token:          c.Token,
		Time:           c.Time,
	}
	pt := o.generatePublishTime(st)
	isTime, err = o.timeS.Execute(pt)
	return
}

func (o *Operate) generatePublishTime(st SendTime) (pt e_time.PublishTime) {
	var ti time.Time
	if st.Time == nil {
		ti = time.Now()
	} else {
		ti = *st.Time
	}
	pt = e_time.PublishTime{
		TemplateId:     st.TemplateId,
		TriggerFrom:    st.TriggerFrom,
		TriggerAccount: st.TriggerAccount,
		Token:          st.Token,
		Time:           ti,
	}
	ttList, err := o.findCache([]int32{int32(st.TemplateId)})
	if err != nil {
		pt.Status = e_time.Failure
		pt.Message = &CannotFindTemplate
		return
	}
	pt.TimeData = e_time_template.Format(ttList)[0].TimeData
	return
}

var StreamComMap = make(map[string]func(rsc map[string]interface{}) (string, error))

func (o *Operate) getStreamComMap() map[string]func(rsc map[string]interface{}) (string, error) {
	StreamComMap["check_time"] = o.streamCheckTime
	return StreamComMap
}

func (o *Operate) streamCheckTime(rsc map[string]interface{}) (result string, err error) {
	var entry SendTime
	err = json.Unmarshal([]byte(rsc["data"].(string)), &entry)
	if err != nil {
		return
	}
	timestamp, err := strconv.ParseInt(rsc["timestamp"].(string), 10, 64)
	if err != nil {
		return
	}
	t := time.Unix(timestamp, 0)
	ct := CheckTime{
		TriggerAccount: entry.TriggerAccount,
		Token:          rsc["callback_token"].(string),
		Time:           &t,
	}
	ct.TriggerFrom = append(entry.TriggerFrom, "stream execute timeTemplate")
	isTime, err := o.CheckTime(entry.TemplateId, ct)
	if err != nil {
		return
	}
	result = strconv.FormatBool(isTime)
	return
}
