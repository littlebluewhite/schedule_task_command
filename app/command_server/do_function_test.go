package command_server

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"schedule_task_command/entry/e_command_template"
	"testing"
)

func TestMonitorData(t *testing.T) {
	t.Run("test1", func(t *testing.T) {

	})
}

func TestStringAnalyze(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		jsonFile, err := os.Open("test/sa_test1.json")
		require.NoError(t, err)
		defer jsonFile.Close()
		data, err := io.ReadAll(jsonFile)
		require.NoError(t, err)
		result := stringAnalyze(data, "root.[]array.user.book.[]array.author")
		require.Equal(t, true, result.getSuccess)
		require.Contains(t, result.arrayResult, "cc")
		require.Contains(t, result.arrayResult, "dd")
	})
	t.Run("test2", func(t *testing.T) {
		jsonFile, err := os.Open("test/sa_test1.json")
		require.NoError(t, err)
		defer jsonFile.Close()
		data, err := io.ReadAll(jsonFile)
		require.NoError(t, err)
		result := stringAnalyze(data, "root.[1]array.user.book.[]array.author")
		require.Equal(t, true, result.getSuccess)
		require.Contains(t, result.arrayResult, "tt")
		require.Contains(t, result.arrayResult, "hh")
	})
	t.Run("test3", func(t *testing.T) {
		jsonFile, err := os.Open("test/sa_test1.json")
		require.NoError(t, err)
		defer jsonFile.Close()
		data, err := io.ReadAll(jsonFile)
		require.NoError(t, err)
		result := stringAnalyze(data, "root.[]array.user.secret")
		require.Equal(t, true, result.getSuccess)
		require.Contains(t, result.arrayResult, "123456")
	})
	t.Run("test4", func(t *testing.T) {
		jsonFile, err := os.Open("test/sa_test1.json")
		require.NoError(t, err)
		defer jsonFile.Close()
		data, err := io.ReadAll(jsonFile)
		require.NoError(t, err)
		result := stringAnalyze(data, "root.[1]array.user.secret")
		require.Equal(t, true, result.getSuccess)
		require.Equal(t, result.valueResult, "123456")
		require.Nil(t, result.arrayResult)
	})
	t.Run("test5", func(t *testing.T) {
		jsonFile, err := os.Open("test/sa_test1.json")
		require.NoError(t, err)
		defer jsonFile.Close()
		data, err := io.ReadAll(jsonFile)
		require.NoError(t, err)
		result := stringAnalyze(data, "root.[2]array.user.secret")
		require.Equal(t, false, result.getSuccess)
	})
	t.Run("test6", func(t *testing.T) {
		jsonFile, err := os.Open("test/sa_test1.json")
		require.NoError(t, err)
		defer jsonFile.Close()
		data, err := io.ReadAll(jsonFile)
		require.NoError(t, err)
		result := stringAnalyze(data, "[2]array.user.secret")
		require.Equal(t, false, result.getSuccess)
	})
	t.Run("test7", func(t *testing.T) {
		jsonFile, err := os.Open("test/sa_test1.json")
		require.NoError(t, err)
		defer jsonFile.Close()
		data, err := io.ReadAll(jsonFile)
		require.NoError(t, err)
		result := stringAnalyze(data, "root.[]array.user")
		require.Equal(t, true, result.getSuccess)
		fmt.Println(result.arrayResult)
		fmt.Println(result.valueResult)
	})
}

func TestAssertValue(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		ar := analyzeResult{getSuccess: false, valueResult: "", arrayResult: []any{}}
		c := e_command_template.MCondition{Order: 1, CalculateType: "=", Value: "wilson"}
		result := assertValue(ar, c)
		require.Equal(t, result.assertSuccess, false)
	})
	t.Run("test2", func(t *testing.T) {
		ar := analyzeResult{getSuccess: true, valueResult: "",
			arrayResult: []any{"cc", "dd", "ee", "tt", "hh", "cc", "cc", "nn", "hh"}}
		c := e_command_template.MCondition{Order: 1, CalculateType: "=", Value: "cc"}
		result := assertValue(ar, c)
		require.Equal(t, result.assertSuccess, false)
	})
	t.Run("test3", func(t *testing.T) {
		ar := analyzeResult{getSuccess: true, valueResult: "",
			arrayResult: []any{"cc", "dd", "ee", "tt", "hh", "cc", "cc", "nn", "hh"}}
		c := e_command_template.MCondition{Order: 1, CalculateType: "include", Value: "tt"}
		result := assertValue(ar, c)
		require.Equal(t, result.assertSuccess, true)
	})
	t.Run("test4", func(t *testing.T) {
		ar := analyzeResult{getSuccess: true, valueResult: "",
			arrayResult: []any{"cc", "dd", "ee", "tt", "hh", "cc", "cc", "nn", "hh"}}
		c := e_command_template.MCondition{Order: 1, CalculateType: "exclude", Value: "c"}
		result := assertValue(ar, c)
		require.Equal(t, result.assertSuccess, true)
	})
	t.Run("test5", func(t *testing.T) {
		ar := analyzeResult{getSuccess: true, valueResult: "",
			arrayResult: []any{"cc", "dd", "ee", "tt", "hh", "cc", "cc", "nn", "hh", 12}}
		c := e_command_template.MCondition{Order: 1, CalculateType: "exclude", Value: "12"}
		result := assertValue(ar, c)
		require.Equal(t, result.assertSuccess, false)
	})
	t.Run("test6", func(t *testing.T) {
		ar := analyzeResult{getSuccess: true, valueResult: "",
			arrayResult: []any{"cc", "dd", "ee", "tt", "hh", "cc", "cc", "nn", "hh", 12}}
		c := e_command_template.MCondition{Order: 1, CalculateType: "include", Value: "12"}
		result := assertValue(ar, c)
		require.Equal(t, result.assertSuccess, true)
	})
	t.Run("test7", func(t *testing.T) {
		ar := analyzeResult{getSuccess: true, valueResult: "",
			arrayResult: []any{"cc", "dd", "ee", "tt", "hh", "cc", "cc", "nn", "hh", 12.5}}
		c := e_command_template.MCondition{Order: 1, CalculateType: "exclude", Value: "12.50"}
		result := assertValue(ar, c)
		require.Equal(t, result.assertSuccess, false)
	})
	t.Run("test8", func(t *testing.T) {
		ar := analyzeResult{getSuccess: true, valueResult: "",
			arrayResult: []any{"cc", "dd", "ee", "tt", "hh", "cc", "cc", "nn", "hh", 12.5}}
		c := e_command_template.MCondition{Order: 1, CalculateType: "include", Value: "12.50"}
		result := assertValue(ar, c)
		require.Equal(t, result.assertSuccess, true)
	})
	t.Run("test9", func(t *testing.T) {
		ar := analyzeResult{getSuccess: true,
			valueResult: "99",
			arrayResult: nil}
		c := e_command_template.MCondition{Order: 1, CalculateType: "=", Value: "99"}
		result := assertValue(ar, c)
		require.Equal(t, result.assertSuccess, true)
	})
	t.Run("test10", func(t *testing.T) {
		ar := analyzeResult{getSuccess: true,
			valueResult: 99,
			arrayResult: nil}
		c := e_command_template.MCondition{Order: 1, CalculateType: "=", Value: "99"}
		result := assertValue(ar, c)
		require.Equal(t, result.assertSuccess, true)
	})
	t.Run("test11", func(t *testing.T) {
		ar := analyzeResult{getSuccess: true,
			valueResult: 99,
			arrayResult: nil}
		c := e_command_template.MCondition{Order: 1, CalculateType: "<=", Value: "100"}
		result := assertValue(ar, c)
		require.Equal(t, result.assertSuccess, true)
	})
	t.Run("test12", func(t *testing.T) {
		ar := analyzeResult{getSuccess: true,
			valueResult: 99,
			arrayResult: nil}
		c := e_command_template.MCondition{Order: 1, CalculateType: ">=", Value: "97.5"}
		result := assertValue(ar, c)
		require.Equal(t, result.assertSuccess, true)
	})
	t.Run("test13", func(t *testing.T) {
		ar := analyzeResult{getSuccess: true,
			valueResult: "99",
			arrayResult: nil}
		c := e_command_template.MCondition{Order: 1, CalculateType: "<=", Value: "102.5"}
		result := assertValue(ar, c)
		require.Equal(t, result.assertSuccess, true)
	})
}

func TestAssertLogic(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		and := "and"
		or := "or"
		asserts := []assertResult{
			{5, true, &and},
			{3, true, &or},
			{4, true, &and},
			{1, true, nil},
			{2, false, &and},
		}
		result := assertLogic(asserts)
		require.Equal(t, true, result)
	})
	t.Run("test2", func(t *testing.T) {
		and := "and"
		//or := "or"
		asserts := []assertResult{
			{1, true, nil},
			{5, true, &and},
			{2, false, &and},
			{4, true, &and},
			{3, true, &and},
		}
		result := assertLogic(asserts)
		require.Equal(t, false, result)
	})
	t.Run("test2", func(t *testing.T) {
		and := "and"
		or := "or"
		asserts := []assertResult{
			{1, false, nil},
			{5, true, &or},
			{2, false, &and},
			{4, false, &or},
			{3, false, &and},
		}
		result := assertLogic(asserts)
		require.Equal(t, true, result)
	})
}
