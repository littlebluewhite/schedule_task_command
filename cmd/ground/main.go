package main

import "fmt"

func main() {
	b := []byte("asf")
	for _, value := range b {
		fmt.Printf("%T, %v", value, value)
	}
	fmt.Println([]byte("\""))
}
