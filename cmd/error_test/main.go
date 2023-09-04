package main

import (
	"errors"
	"fmt"
)

func main() {
	a := CCC
	b := CCC
	fmt.Printf("type: %T\n", a)
	fmt.Printf("type: %T\n", b)
	fmt.Println(errors.Is(a, b))

	var c error = AAA{aaa: "aa"}
	var d error = AAA{aaa: "bb"}
	fmt.Printf("type: %T, %e\n", c, c)
	fmt.Printf("type: %T, %e\n", d, d)
	fmt.Println(errors.Is(c, d))
}

var CCC = errors.New("test error")

type AAA struct {
	aaa string
}

func (a AAA) Error() string {
	return "AAA error"
}
