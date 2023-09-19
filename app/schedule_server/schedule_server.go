package schedule_server

import (
	"context"
	"fmt"
	"schedule_task_command/app/dbs"
	"schedule_task_command/dal/model"
	"schedule_task_command/entry/e_schedule"
	"schedule_task_command/util/logFile"
	"sync"
	"time"
)

type ScheduleServer[T, U any] struct {
	dbs   dbs.Dbs
	l     logFile.LogFile
	taskS taskServer
	timeS timeServer
}

func NewScheduleServer[T, U any](dbs dbs.Dbs, taskS taskServer, timeS timeServer) *ScheduleServer[T, U] {
	l := logFile.NewLogFile("app", "schedule_server")
	return &ScheduleServer[T, U]{
		dbs:   dbs,
		l:     l,
		taskS: taskS,
		timeS: timeS,
	}
}

func (s *ScheduleServer[T, U]) Start(ctx context.Context, interval, removeTime time.Duration) {
	s.l.Info().Println("Schedule server started")
	defer s.l.Error().Println("Schedule server stopped")
	wg := &sync.WaitGroup{}
	wg.Add(3)
	go func(wg *sync.WaitGroup) {
		s.listen(ctx, interval)
		wg.Done()
	}(wg)
	go func(wg *sync.WaitGroup) {
		s.taskS.Start(ctx, removeTime)
		wg.Done()
	}(wg)
	go func(wg *sync.WaitGroup) {
		s.timeS.Start(ctx)
		wg.Done()
	}(wg)
	wg.Wait()
}

func (s *ScheduleServer[T, U]) listen(ctx context.Context, duration time.Duration) {
	ticker := time.NewTicker(duration)
	for {
		select {
		case <-ctx.Done():
			s.l.Info().Println("Schedule server stopped")
			return
		case t := <-ticker.C:
			go s.checkSchedule(ctx, t)
			fmt.Println("Invoked at ", t)
		default:
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
				isTime := schedule.CheckTimeData(t)
				isActive := isTime && schedule.Enabled
				if isActive {
					// Task execute
					scheduleId := fmt.Sprintf("%d", schedule.ID)
					triggerFrom := []string{"schedule", scheduleId}
					now := time.Now()
					token := fmt.Sprintf("schedule-%s-%s-%s", scheduleId, schedule.Tags, now)
					s.taskS.Execute(ctx, int(schedule.TaskID), triggerFrom, "", token)
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

func (s *ScheduleServer[T, U]) GetTimeServer() U {
	return s.timeS.(U)
}

func (s *ScheduleServer[T, U]) GetTaskServer() T {
	return s.taskS.(T)
}
