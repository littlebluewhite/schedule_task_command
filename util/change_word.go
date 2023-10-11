package util

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

func ChangeByteVariables(data []byte, variables map[string]string) ([]byte, error) {
	dataC := slices.Clone(data)
	r, _ := regexp.Compile(`{{([^{}]*)}}`)
	c := r.FindAllIndex(dataC, -1)
	for i := len(c) - 1; i >= 0; i-- {
		start := c[i][0]
		end := c[i][1]
		v := string(dataC[start+2 : end-2])
		word, ok := variables[v]
		if !ok {
			return nil, fmt.Errorf("variable: %v not found", v)
		}
		wordBytes := []byte(word)
		dataC = slices.Replace(dataC, start, end, wordBytes...)
	}
	return dataC, nil
}

func ChangeStringVariables(data string, variables map[string]string) (string, error) {
	r, _ := regexp.Compile(`{{([^{}]*)}}`)
	c := r.FindAllStringIndex(data, -1)
	variableSlice := make([]string, 0, len(c))
	for i := 0; i < len(c); i++ {
		start := c[i][0]
		end := c[i][1]
		variableSlice = append(variableSlice, data[start+2:end-2])
	}
	for _, v := range variableSlice {
		word, ok := variables[v]
		if !ok {
			return "", fmt.Errorf("variable: %v not found", v)
		}
		oldWord := fmt.Sprintf("{{%s}}", v)
		data = strings.Replace(data, oldWord, word, 1)
	}
	return data, nil
}

func SliceT(data []int) {
	s := []int{9, 9, 9, 9, 9}
	d := slices.Replace(data, 2, 4, s...)
	fmt.Println(d)
	a := "asdfasdf"
	a = strings.Replace(a, "asd", "wwwwww", -1)
	fmt.Println(a)
}

func MapT() {
	a := make(map[string]string)
	fmt.Println(a == nil)
	a["aa"] = "aa"
}
