package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {
	i, err := strconv.ParseInt("1405544146", 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(i, 0)
	fmt.Println(tm)
}
