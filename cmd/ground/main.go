package main

import "fmt"

func main() {
	a := 1
	b := 2
	c := 3
	switch {
	case a == 1:
		fmt.Println("a ok")
	case b == 2:
		fmt.Println("b ok")
	default:
		fmt.Println("c ok", c)
	}
}

type stageMapValue2 struct {
	a map[string]int
}
