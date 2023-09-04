package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()
	p := fmt.Println
	p(time.Unix(1654963822, 4567).UnixNano())
	a := time.Date(2023, time.June, 16, 12, 0, 0, 0, time.Local)
	b := time.Date(2023, time.June, 16, 12, 0, 0, 0, time.UTC)
	p(now.In(time.UTC))
	p(a)
	p(b)
	p(b.Sub(a))
	p(now.Location())
}
