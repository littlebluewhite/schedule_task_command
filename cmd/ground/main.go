package main

import (
	"fmt"
	"github.com/goccy/go-json"
)

func main() {
	v := make(map[string]string)
	if err := json.Unmarshal([]byte("[]"), &v); err != nil {
		fmt.Println(err)
	}
	for key, value := range map[string]string{} {
		v[key] = value
	}
	fmt.Println(v == nil)
	fmt.Println(v)
}
