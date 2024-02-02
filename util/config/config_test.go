package config

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
)

func TestConfig(t *testing.T) {

	_, b, _, _ := runtime.Caller(0)
	rootPath := filepath.Dir(filepath.Dir(filepath.Dir(b)))
	t.Run("test1", func(t *testing.T) {
		a := NewConfig[Config](rootPath, "config", "config", Ini)
		fmt.Println(a)
	})
	t.Run("test1", func(t *testing.T) {
		a := NewConfig[Config](rootPath, "config", "config", Yaml)
		fmt.Println(a)
	})
}
