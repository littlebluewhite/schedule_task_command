package main

import (
	"fmt"
	"github.com/goccy/go-json"
)

func main() {
	v := make(map[string]string)
	//var v map[string]string
	b := []byte(`{"id": "1", "name": "aaa"}`)
	e := json.Unmarshal(b, &v)
	fmt.Println(e)
	a := map[string]string{
		"a": "fff",
		"b": "888",
	}
	fmt.Println(a == nil)
	if a != nil {
		for key, value := range a {
			v[key] = value
		}
	}
	fmt.Println(v)
}
