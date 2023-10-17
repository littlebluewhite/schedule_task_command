package main

import (
	"fmt"
	"sync/atomic"
)

func main() {
	a := atomic.Uint64{}
	a.Store(10)
	b := a.Add(1)
	fmt.Println(b)
	c := a.Load()
	fmt.Println(c)
}
