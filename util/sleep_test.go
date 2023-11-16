package util

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestSleep(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		err := SleepWithContext(ctx, 5*time.Second)
		fmt.Println("err: ", err)
	})
	t.Run("test2", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
		defer cancel()
		err := SleepWithContext(ctx, 5*time.Second)
		fmt.Println("err: ", err)

	})
	t.Run("test2", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		go func() {
			time.Sleep(2 * time.Second)
			cancel()
		}()
		err := SleepWithContext(ctx, 10*time.Second)
		fmt.Println("err: ", err)
	})
}
