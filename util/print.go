package util

import "fmt"

func P[T any](v T) {
	fmt.Println(v)
}
