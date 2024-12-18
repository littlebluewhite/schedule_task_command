package command_server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/littlebluewhite/schedule_task_command/entry/e_command"
	"github.com/littlebluewhite/schedule_task_command/entry/e_command_template"
	"github.com/littlebluewhite/schedule_task_command/util"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func (c *CommandServer) requestProtocol(ctx context.Context, com e_command.Command) e_command.Command {
	for {
		select {
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.Canceled) {
				com.Status = e_command.Cancel
				if com.Message == nil {
					com.Message = &CommandCanceled
				}
			} else if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				com.Status = e_command.Failure
				com.Message = &CommandTimeout
			}
			return com
		default:
			// "condition not match" error cancel
			com.Message = nil

			switch com.CommandData.Protocol {
			case e_command_template.Http:
				com = c.doHttp(ctx, com)
			case e_command_template.Websocket:
			case e_command_template.Mqtt:
			case e_command_template.RedisTopic:
			default:
			}
			// variables error only
			if com.Message != nil {
				// variables failed or http request failed
				return com
			} else {
				// send command successfully
				if com.CommandData.Monitor == nil {
					// mode execute
					com.Status = e_command.Success
				} else {
					// mode monitor
					myErr := monitorData(com, *com.CommandData.Monitor)
					com.Message = myErr
					if myErr == nil {
						com.Status = e_command.Success
					}
				}
				c.l.Infof("id: %d \ncommand status: %v\nrequest result: %s\n", com.ID, com.Status, com.RespData)
				if com.Status == e_command.Success {
					// get return parser
					com.Return = parserData(com)
					return com
				} else {
					_ = util.SleepWithContext(ctx, time.Duration(com.CommandData.Timeout)*time.Millisecond)
				}
			}
		}
	}
}

func (c *CommandServer) doHttp(ctx context.Context, com e_command.Command) e_command.Command {
	var body io.Reader
	h := com.CommandData.Http
	var contentType string
	if h.Body != nil {
		changeBody, err := util.ChangeByteVariables(h.Body, com.Variables)
		if err != nil {
			com.Status = e_command.Failure
			com.Message = &URLVariables
			return com
		}
		switch h.BodyType {
		case e_command_template.Json:
			body = bytes.NewBuffer(changeBody)
			contentType = "application/json"
		case e_command_template.FormData:
			//TODO form data body
			contentType = "multipart/form-data"
		case e_command_template.XWWWFormUrlencoded:
			//TODO x_www_form_urlencoded body
			contentType = "application/x-www-form-urlencoded"
		default:
		}
	}
	header := make([]httpHeader, 0, 20)
	fullUrl, e := util.ChangeStringVariables(h.URL, com.Variables)
	if e != nil {
		com.Status = e_command.Failure
		com.Message = &URLVariables
		return com
	}
	// fullUrl has params
	if index := strings.Index(fullUrl, "?"); index != -1 {
		params := url.Values{}
		rUrl := fullUrl[:index]
		rParams := fullUrl[index+1:]
		pSlice := strings.Split(rParams, "&")
		for _, p := range pSlice {
			keyValue := strings.Split(p, "=")
			params.Add(keyValue[0], keyValue[1])
		}
		fullUrl = rUrl + "?" + params.Encode()
	}

	req, e := http.NewRequestWithContext(ctx, h.Method.String(), fullUrl, body)
	if e != nil {
		com.Status = e_command.Failure
		com.Message = &HttpTimeout
		return com
	}
	if h.Header != nil {
		hh, err := util.ChangeByteVariables(h.Header, com.Variables)
		if err != nil {
			com.Status = e_command.Failure
			com.Message = &HeaderVariables
			return com
		}
		if e := json.Unmarshal(hh, &header); e != nil {
			c.l.Errorf("id: %d header unmarshal failed", com.ID)
		}
	}
	for _, item := range header {
		if item.IsActive {
			req.Header.Set(item.Key, item.Value)
		}
	}
	req.Header.Set("Content-Type", contentType)
	var resp *http.Response
	resp1, e := c.httpClient.Do(req)
	if e != nil {
		c.l.Errorf("id: %d request failed, template id: %d", com.ID, com.TemplateId)
		if resp1 == nil {
			c.l.Errorf("request: %+v, and response is nil", req)
			com.Status = e_command.Failure
			com.Message = &RequestErr
			return com
		}
	}
	resp = resp1
	com.StatusCode = resp.StatusCode
	if respBody1, e := io.ReadAll(resp.Body); e != nil {
		com.RespData = []byte{}
		c.l.Errorf("id: %d request body failed", com.ID)
		return com
	} else {
		com.RespData = respBody1
		// check respBody1
		b, err := json.Marshal(com.RespData)
		if err != nil {
			// to json string
			b = append([]byte("\""), respBody1...)
			b = append(b, []byte("\"")...)
			com.RespData = b
		}

	}
	defer func() {
		if e := resp.Body.Close(); e != nil {
			c.l.Errorln("Response body closed failed")
		}
	}()
	c.l.Infof("id: %d \nrequest result: %s, status code: %d\n", com.ID, com.RespData, com.StatusCode)
	return com
}

func monitorData(com e_command.Command, m e_command_template.Monitor) (err *util.MyErr) {
	if com.StatusCode != int(m.StatusCode) {
		err = &HttpCodeErr
		return
	}
	asserts := make([]assertResult, 0, len(m.MConditions))
	for _, condition := range m.MConditions {
		searchRule, _ := util.ChangeStringVariables(condition.SearchRule, com.Variables)
		value, _ := util.ChangeStringVariables(condition.Value, com.Variables)
		condition.Value = value
		ar := stringAnalyze(com.RespData, searchRule)
		assert := assertValue(ar, condition)
		asserts = append(asserts, assert)
	}
	logicResult := assertLogic(asserts)
	if !logicResult {
		err = &ConditionFailed
	}
	return
}

func parserData(com e_command.Command) map[string]string {
	parserReturn := make(map[string]string)
	for _, pr := range com.CommandData.ParserReturn {
		key, _ := util.ChangeStringVariables(pr.Key, com.Variables)
		searchRule, _ := util.ChangeStringVariables(pr.SearchRule, com.Variables)
		parserReturn[key] = fmt.Sprintf("%v", stringAnalyze(com.RespData, searchRule).valueResult)
	}
	// add default parser data
	addDefaultParserReturn(parserReturn, com)
	return parserReturn
}

func stringAnalyze(data []byte, rule string) (result analyzeResult) {
	r := strings.Split(rule, ".")
	// "root.person.[all]array.name
	var f []any
	var arrayFlag bool
	var d any
	e := json.Unmarshal(data, &d)
	if e != nil {
		return
	}
	f = append(f, d)
	for _, word := range r[1:] {
		var handleFunc func(word string, find []any) ([]any, bool)
		if strings.Index(word, "array") == -1 {
			handleFunc = handleKey
		} else {
			handleFunc = handleArray
		}
		var flag bool
		f, flag = handleFunc(word, f)
		if flag {
			arrayFlag = true
		}
	}
	if len(f) > 0 {
		result.getSuccess = true
	} else {
		return
	}
	if arrayFlag {
		result.arrayResult = f
	} else {
		result.valueResult = f[0]
	}
	return
}

func assertValue(ar analyzeResult, condition e_command_template.MCondition) (a assertResult) {
	a.order = condition.Order
	a.preLogicType = condition.PreLogicType
	if ar.getSuccess == false {
		return
	}
	if ar.valueResult != nil && util.Contains([]string{condition.CalculateType}, valueCalculate) {
		a.assertSuccess = assertSingle(ar.valueResult, condition.Value, condition.CalculateType)
	} else if ar.arrayResult != nil && util.Contains([]string{condition.CalculateType}, sliceCalculate) {
		a.assertSuccess = assertArray(ar.arrayResult, condition.Value, condition.CalculateType)
	}
	return
}

func assertSingle(result any, cv, c string) (r bool) {
	switch result.(type) {
	case string:
		r = assertString(result.(string), cv, c)
	case int:
		r = assertInt(result.(int), cv, c)
	case float64:
		r = assertFloat(result.(float64), cv, c)
	default:
		fmt.Printf("%T, %v", result, result)
	}
	return
}

func assertString(v string, cv, c string) (r bool) {
	switch c {
	case "=":
		if v == cv {
			r = true
		}
	case "!=":
		if v != cv {
			r = true
		}
	default:
		vNum, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return
		}
		_, err = strconv.ParseFloat(cv, 64)
		if err != nil {
			return
		}
		return assertFloat(vNum, cv, c)
	}
	return
}

func assertInt(v int, cv, c string) (r bool) {
	_, err := strconv.ParseFloat(cv, 64)
	if err != nil {
		return
	}
	return assertFloat(float64(v), cv, c)
}

func assertFloat(v float64, cv, c string) (r bool) {
	cNum, err := strconv.ParseFloat(cv, 64)
	if err != nil {
		return
	}
	switch c {
	case "=":
		if v == cNum {
			r = true
		}
	case "!=":
		if v != cNum {
			r = true
		}
	case "<":
		if v < cNum {
			r = true
		}
	case "<=":
		if v <= cNum {
			r = true
		}
	case ">":
		if v > cNum {
			r = true
		}
	case ">=":
		if v >= cNum {
			r = true
		}
	}
	return
}

func assertArray(result []any, cv, calculateType string) (r bool) {
	switch calculateType {
	case "include":
		r = handleInclude(result, cv)
	case "exclude":
		r = handleExclude(result, cv)
	}
	return
}

func handleInclude(data []any, cv string) (r bool) {
	for _, item := range data {
		switch item.(type) {
		case string:
			if item.(string) == cv {
				r = true
				return
			}
		case float64:
			cNum, err := strconv.ParseFloat(cv, 64)
			if err != nil {
				continue
			}
			if item.(float64) == cNum {
				r = true
				return
			}
		case int:
			cNum, err := strconv.ParseInt(cv, 10, 64)
			if err != nil {
				continue
			}
			if item.(int) == int(cNum) {
				r = true
				return
			}
		default:
			fmt.Printf("%T, %v", item, item)
			continue
		}
	}
	return
}

func handleExclude(data []any, cv string) (r bool) {
	for _, item := range data {
		switch item.(type) {
		case string:
			if item.(string) == cv {
				return
			}
		case float64:
			cNum, err := strconv.ParseFloat(cv, 64)
			if err != nil {
				continue
			}
			if item.(float64) == cNum {
				return
			}
		case int:
			cNum, err := strconv.ParseInt(cv, 10, 64)
			if err != nil {
				continue
			}
			if item.(int) == int(cNum) {
				return
			}
		default:
			continue
		}
	}
	r = true
	return
}

func handleArray(word string, find []any) (result []any, flag bool) {
	re, _ := regexp.Compile(`\[([0-9]*)]`)
	indexes := re.FindStringSubmatchIndex(word)
	index := word[indexes[2]:indexes[3]]
	if index == "" {
		result = handleArrayAll(find)
		flag = true
	} else {
		result = handleArrayIndex(index, find)
	}
	return
}

func handleArrayAll(find []any) (result []any) {
	for _, item := range find {
		s, ok := item.([]any)
		if !ok {
			continue
		}
		for _, v := range s {
			result = append(result, v)
		}
	}
	return
}

func handleArrayIndex(index string, find []any) (result []any) {
	for _, item := range find {
		num, err := strconv.ParseInt(index, 10, 64)
		if err != nil {
			continue
		}
		s, ok := item.([]any)
		if !ok {
			continue
		}
		if num < 0 || int(num) >= len(s) {
			continue
		}
		result = append(result, s[num])
	}
	return
}

func handleKey(word string, find []any) (result []any, flag bool) {
	for _, item := range find {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		item, ok = m[word]
		if !ok {
			continue
		}
		result = append(result, item)
	}
	return
}

func assertLogic(asserts []assertResult) (result bool) {
	sort.Slice(asserts, func(i, j int) bool {
		return asserts[i].order < asserts[j].order
	})
	orSlice := make([]bool, 0, len(asserts))
	pre := true
	for i, assert := range asserts {
		if assert.preLogicType == nil && i == 0 {
			pre = assert.assertSuccess
			continue
		}
		switch *assert.preLogicType {
		case "and":
			pre = pre && assert.assertSuccess
		case "or":
			orSlice = append(orSlice, pre)
			pre = assert.assertSuccess
		}
	}
	orSlice = append(orSlice, pre)
	result = util.Contains[bool]([]bool{true}, orSlice)
	return
}
