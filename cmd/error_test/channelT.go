package main

import "fmt"

func main() {
	_ = AA()
}

func AA() error {
	ch := make(chan error, 10)
	defer close(ch)
	for i := 0; i < 10; i++ {
		go func(id int) {
			if id == 11 {
				ch <- fmt.Errorf("id: %d, error", id)
			} else {
				fmt.Printf("id: %d, no error\n", id)
				ch <- nil
			}
		}(i)
	}
	for i := 0; i < 10; i++ {
		select {
		case e := <-ch:
			if e != nil {
				fmt.Println(e)
				return e
			}
		}
	}
	return nil
}
