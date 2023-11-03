package redis_stream

func CreateRedisStreamCom() map[string]interface{} {
	return map[string]interface{}{
		"command":                  "",
		"timestamp":                "",
		"data":                     "",
		"callback_command":         "",
		"callback_channel":         "",
		"is_wait_call_back":        "",
		"callback_token":           "",
		"callback_timeout":         "",
		"callback_until_feed_back": "",
		"command_sk":               "",
		"send_pattern":             "",
		"status_code":              "",
	}
}
