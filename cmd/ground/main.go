package main

import "fmt"

func main() {
	var a []int
	fmt.Println(a == nil)
	fmt.Println(a)
	b := make([]int, 0, 0)
	fmt.Println(b == nil)
	fmt.Println(b)
}
