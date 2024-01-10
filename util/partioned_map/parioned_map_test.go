package partioned_map

import (
	"fmt"
	"sync"
	"testing"
)

func BenchmarkStd(b *testing.B) {
	m := make(map[string]int)
	b.Run("set std concurrently", func(b *testing.B) {
		var wg sync.WaitGroup
		var mu sync.RWMutex
		for i := 0; i < b.N; i++ {
			wg.Add(1)
			i := i
			go func() {
				mu.Lock()
				m[fmt.Sprint(i)] = i
				mu.Unlock()
				wg.Done()
			}()
		}
		wg.Wait()
	})
}

func BenchmarkSyncStd(b *testing.B) {
	b.Run("set sync map std concurrently", func(b *testing.B) {
		var m sync.Map
		var wg sync.WaitGroup
		for i := 0; i < b.N; i++ {
			wg.Add(1)
			i := i
			go func() {
				m.Store(fmt.Sprint(i), i)
				wg.Done()
			}()
		}
		wg.Wait()
	})
}

func BenchmarkPartitioned(b *testing.B) {
	m := NewPartitionedMap[int](&hashSumPartitioner{1000}, 1000)
	b.Run("set partitioned concurrently", func(b *testing.B) {
		var wg sync.WaitGroup
		for i := 0; i < b.N; i++ {
			wg.Add(1)
			i := i
			go func() {
				m.Set(fmt.Sprint(i), i)
				wg.Done()
			}()
		}
		wg.Wait()
	})
}
