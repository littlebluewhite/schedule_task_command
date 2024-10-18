package schedule_server

import (
	"context"
	"fmt"
	"schedule_task_command/api"
	"schedule_task_command/dal/model"
	"schedule_task_command/entry/e_command_template"
	"schedule_task_command/entry/e_schedule"
	"schedule_task_command/entry/e_task"
	"schedule_task_command/entry/e_task_template"
	"schedule_task_command/entry/e_time"
	"schedule_task_command/util/my_log"
	"sync"
	"time"
)

type ScheduleServer[T, U any] struct {
	dbs   api.Dbs
	l     api.Logger
	taskS taskServer
	timeS timeServer
	wg    *sync.WaitGroup
}

func NewScheduleServer[T, U any](dbs api.Dbs, taskS taskServer, timeS timeServer) *ScheduleServer[T, U] {
	l := my_log.NewLog("app/schedule_server")
	return &ScheduleServer[T, U]{
		dbs:   dbs,
		l:     l,
		taskS: taskS,
		timeS: timeS,
		wg:    new(sync.WaitGroup),
	}
}

func (s *ScheduleServer[T, U]) Start(ctx context.Context, interval, removeTime time.Duration) {
	s.l.Infoln("Schedule server started")

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.listen(ctx, interval)
	}()
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.taskS.Start(ctx, removeTime)
	}()
	go func() {
		s.timeS.Start(ctx)
	}()
}

func (s *ScheduleServer[T, U]) Close() {
	s.wg.Wait()
	s.taskS.Close()
	s.l.Infoln("Schedule server stop gracefully")
}

func (s *ScheduleServer[T, U]) listen(ctx context.Context, duration time.Duration) {
	ticker := time.NewTicker(duration)
	for {
		select {
		case <-ctx.Done():
			s.l.Infoln("Schedule server stopped")
			return
		case t := <-ticker.C:
			go s.checkSchedule(ctx, t)
			//fmt.Println("Invoked at ", t)
		}
	}
}

func (s *ScheduleServer[T, U]) checkSchedule(ctx context.Context, t time.Time) {
	select {
	case <-ctx.Done():
		return
	default:
		wg := new(sync.WaitGroup)
		cacheMap := s.getSchedule()
		wg.Add(len(cacheMap))
		for _, sItem := range cacheMap {
			go func(wg *sync.WaitGroup, schedule e_schedule.Schedule, t time.Time) {
				defer wg.Done()
				if !schedule.Enabled {
					return
				}
				scheduleId := fmt.Sprintf("%d", schedule.ID)
				token := fmt.Sprintf("schedule-%s-%s", scheduleId, t)
				triggerFrom := []map[string]string{{"schedule": scheduleId}}
				// Check time
				pt := e_time.PublishTime{
					TriggerFrom: triggerFrom,
					Token:       token,
					Time:        t,
					TimeData:    schedule.TimeData,
				}
				isTime, _ := s.timeS.Execute(pt)
				if isTime {
					// Task execute
					s.l.Infof("id: %s execute", scheduleId)
					st := e_task_template.SendTaskTemplate{
						TemplateId:     int(schedule.TaskTemplateID),
						Source:         "Schedule",
						TriggerFrom:    triggerFrom,
						TriggerAccount: "system",
						Token:          token,
					}
					task := s.generateTask(st)
					_ = s.taskS.ExecuteWait(ctx, task)
				}
			}(wg, sItem, t)
		}
		wg.Wait()
	}
}

func (s *ScheduleServer[T, U]) getSchedule() map[int]e_schedule.Schedule {
	cacheMap := make(map[int]e_schedule.Schedule)
	if x, found := s.dbs.GetCache().Get("schedules"); found {
		c := x.(map[int]model.Schedule)
		for key, value := range c {
			cacheMap[key] = e_schedule.Model2Entry(value)
		}
	}
	return cacheMap
}

func (s *ScheduleServer[T, U]) generateTask(st e_task_template.SendTaskTemplate) (task e_task.Task) {
	task = e_task.Task{
		TemplateId:     st.TemplateId,
		Source:         st.Source,
		TriggerFrom:    st.TriggerFrom,
		TriggerAccount: st.TriggerAccount,
		Token:          st.Token,
	}
	var cacheMap map[int]model.TaskTemplate
	if x, found := s.dbs.GetCache().Get("taskTemplates"); found {
		cacheMap = x.(map[int]model.TaskTemplate)
	}
	mt, ok := cacheMap[st.TemplateId]
	if !ok {
		task.Status = e_task.Failure
		task.Message = &e_command_template.CannotFindTemplate
		return
	}
	task.TaskData = e_task_template.Format([]model.TaskTemplate{mt})[0]

	trigger := map[string]string{"task_template": fmt.Sprintf("%d", st.TemplateId)}
	task.TriggerFrom = append(task.TriggerFrom, trigger)
	return
}

func (s *ScheduleServer[T, U]) GetTimeServer() U {
	return s.timeS.(U)
}

func (s *ScheduleServer[T, U]) GetTaskServer() T {
	return s.taskS.(T)
}
