package partioned_map

import "sync"

type partition[T any] struct {
	storage map[string]T
	sync.RWMutex
}

func (p *partition[T]) set(key string, value T) {
	p.Lock()
	p.storage[key] = value
	p.Unlock()
}

func (p *partition[T]) get(key string) (T, bool) {
	p.RLock()
	v, ok := p.storage[key]
	if !ok {
		p.RUnlock()
		return v, false
	}
	p.RUnlock()
	return v, true
}
