package util

import (
	"github.com/stretchr/testify/require"
	"testing"
)

var weekDay = []int{0, 1, 2, 3, 4, 5, 6}

func TestContain(t *testing.T) {
	t.Run("int forTest", func(t *testing.T) {
		value := Contains([]int{1, 2, 3, 4}, weekDay)
		require.Equal(t, value, true)
	})
	t.Run("int test2", func(t *testing.T) {
		value := Contains([]int{1, 4, 5, 7}, weekDay)
		require.Equal(t, value, false)
	})
	t.Run("int test3", func(t *testing.T) {
		value := Contains([]int{}, weekDay)
		require.Equal(t, value, true)
	})
	t.Run("int test4", func(t *testing.T) {
		value := Contains([]int{0, 8}, weekDay)
		require.Equal(t, value, false)
	})
	t.Run("string forTest", func(t *testing.T) {
		value := Contains([]string{"a", "c"}, []string{"a", "b", "c"})
		require.Equal(t, value, true)
	})
	t.Run("string test2", func(t *testing.T) {
		value := Contains([]string{"a", "c", "b"}, []string{"a", "b", "c"})
		require.Equal(t, value, true)
	})
}
