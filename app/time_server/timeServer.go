package time_server

import (
	"context"
	"fmt"
	"schedule_task_command/app/dbs"
	"schedule_task_command/dal/model"
	"schedule_task_command/util/logFile"
	"sync"
	"time"
)

var (
	timeServerLog logFile.LogFile
)

func init() {
	timeServerLog = logFile.NewLogFile("app", "time_server")
}

type TimeServer interface {
	Start(ctx context.Context)
	ReloadSchedule(msm map[int]model.Schedule)
}

type timeServer[T any] struct {
	mu       *sync.RWMutex
	dbs      dbs.Dbs
	duration time.Duration
	schedule map[int]schedule
}

func NewTimeServer[T any](dbs dbs.Dbs, duration time.Duration) TimeServer {
	mu := new(sync.RWMutex)
	return &timeServer[T]{
		dbs:      dbs,
		duration: duration,
		mu:       mu,
	}
}

func (ts *timeServer[T]) Start(ctx context.Context) {
	timeServerLog.Info().Println("Time server start")
	ticker := time.NewTicker(ts.duration)
	ctxChild, cancel := context.WithCancel(ctx)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			timeServerLog.Info().Println("Time server stop gracefully")
			return
		case t := <-ticker.C:
			go ts.checkSchedule(ctxChild, t)
			fmt.Println("Invoked at ", t)
		}
	}
}

func (ts *timeServer[T]) checkSchedule(ctx context.Context, t time.Time) {
	select {
	case <-ctx.Done():
		return
	default:
		var wg sync.WaitGroup
		for _, s := range ts.schedule {
			wg.Add(1)
			go func(s schedule, t time.Time, wg *sync.WaitGroup, mu *sync.RWMutex) {
				mu.RLock()
				defer mu.RUnlock()
				defer wg.Done()
				isActive := checkScheduleActive(s, t)
				if isActive {
					// TODO execute task
				}
				timeServerLog.Info().Printf("id: %v, active: %v\n", s.ID, isActive)
			}(s, t, &wg, ts.mu)
		}
		wg.Wait()
	}
}

func (ts *timeServer[T]) ReloadSchedule(msm map[int]model.Schedule) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.schedule = modelMap2scheduleMap(msm)
}
