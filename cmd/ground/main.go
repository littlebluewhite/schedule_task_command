package main

import "fmt"

func main() {
	ch := make(chan int)

	go func() {
		ch <- 42
		close(ch)
	}()

	// Assuming some work is being done here...

	value := <-ch
	fmt.Println(value)
}
