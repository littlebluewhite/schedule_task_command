package util

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestChangeVariables(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		m := map[string]string{"name": "wilson", "id": "123456wwww"}
		data, err := ChangeByteVariables([]byte("{name:{{name}}, id:{{id}}}, last_name:{{name}}}"), m)
		fmt.Println(string(data))
		require.NoError(t, err)
	})
	t.Run("test2", func(t *testing.T) {
		m := map[string]string{"id": "123456wwww"}
		data, err := ChangeByteVariables([]byte("name: {{name}}, id:{{id}}"), m)
		require.Nil(t, data)
		require.Error(t, err)
	})
}
func TestChangeStringVariables(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		m := map[string]string{"name": "wilson", "id": "123456wwww"}
		data := "{name:{{name}}, id:{{id}}, last_name:{{name}}}"
		data2, err := ChangeStringVariables(data, m)
		fmt.Println("data: ", data)
		fmt.Println(data2)
		require.NoError(t, err)
	})
	t.Run("test2", func(t *testing.T) {
		m := map[string]string{"id": "123456wwww"}
		data, err := ChangeStringVariables("name: {{name}}, id:{{id}}", m)
		require.Equal(t, data, "")
		require.Error(t, err)
	})
}
func TestSlice(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		SliceT([]int{1, 2, 3, 4, 5, 6, 7, 8})
	})
	t.Run("test2", func(t *testing.T) {
		MapT()
	})
}
