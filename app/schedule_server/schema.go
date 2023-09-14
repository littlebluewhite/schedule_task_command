package schedule_server

import (
	"sync"
)

type taskServer interface{}

type chs struct {
	mu *sync.RWMutex
}

type ScheduleSer interface{}
