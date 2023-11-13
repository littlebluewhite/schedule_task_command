package main

import (
	"fmt"
	"github.com/goccy/go-json"
)

func main() {
	var a []int32
	b := []byte("[1, 2, 3, 4]")
	e := json.Unmarshal(b, &a)

	fmt.Println(e)
	fmt.Println(a)
}

type stageMapValue2 struct {
	a map[string]int
}
