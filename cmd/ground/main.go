package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	ctx := context.Background()
	baseURL := "http://192.168.1.10:9327/PropertyManagement/send_repair_request"
	params := url.Values{}
	params.Add("object_uid", "PM@MCB")
	params.Add("alarm_id", "asdfasdf")
	params.Add("alarm_level", "6")
	params.Add("title", "智慧電表 PM@MCB狀態異常")

	encodedParams := params.Encode()
	fullURL := baseURL + "?" + encodedParams

	//fullURL = "http://192.168.1.10:9327/PropertyManagement/send_repair_request?object_uid=PM@MCB&alarm_id=asdfasdf&alarm_level=6&title=智慧電表 PM@MCB狀態異常"
	body := strings.NewReader("") // Replace with actual body if needed
	req, e := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, body)
	if e != nil {
		panic(e)
	}
	req.Header.Set("Content-Type", "application/json")
	fmt.Println(req.Header)
	client := &http.Client{}
	resp1, e := client.Do(req)
	if e != nil {
		panic(e)
	}
	fmt.Println(resp1)

	a := "aaaaa?bbbbb"
	v := strings.Index(a, "?")
	fmt.Println(v)
	fmt.Println(a[:v])
	fmt.Println(a[v+1:])
}
