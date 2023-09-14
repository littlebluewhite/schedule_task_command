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

type ScheduleServer struct {
	dbs   dbs.Dbs
	l     logFile.LogFile
	taskS taskServer
	chs   chs
}

func NewScheduleServer(dbs dbs.Dbs, taskS taskServer) *ScheduleServer {
	l := logFile.NewLogFile("app", "schedule_server")
	mu := new(sync.RWMutex)
	return &ScheduleServer{
		dbs:   dbs,
		l:     l,
		taskS: taskS,
		chs: chs{
			mu: mu,
		},
	}
}

func (s *ScheduleServer) Start(ctx context.Context, duration time.Duration) {
	s.l.Info().Println("Schedule server started")
	ticker := time.NewTicker(duration)
	ctxChild, cancel := context.WithCancel(ctx)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			s.l.Info().Println("Schedule server")
			return
		case t := <-ticker.C:
			go s.checkSchedule(ctxChild, t)
			fmt.Println("Invoked at ", t)
		default:
		}
	}
}

func (s *ScheduleServer) checkSchedule(ctx context.Context, t time.Time) {
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
				td := schedule.GetTimeData()
				isTime := td.CheckTimeData(t)
				isActive := isTime && schedule.Enabled
				if isActive {
					// TODO Task
				}
			}(wg, sItem, t)
		}
	}
}

func (s *ScheduleServer) getSchedule() map[int]e_schedule.Schedule {
	cacheMap := make(map[int]e_schedule.Schedule)
	if x, found := s.dbs.GetCache().Get("schedules"); found {
		c := x.(map[int]model.Schedule)
		for key, value := range c {
			cacheMap[key] = e_schedule.Model2Entry(value)
		}
	}
	return cacheMap
}
