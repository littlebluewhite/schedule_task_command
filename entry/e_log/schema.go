package e_log

import "time"

type Log struct {
	Timestamp     float64   `json:"timestamp"`
	Account       string    `json:"account"`
	ContentLength int       `json:"content_length"`
	Datetime      time.Time `json:"datetime"`
	IP            string    `json:"ip"`
	Referer       string    `json:"referer"`
	ApiUrl        string    `json:"api_url"`
	Method        string    `json:"method"`
	Module        string    `json:"module"`
	StatusCode    int       `json:"status_code"`
	Token         string    `json:"token"`
	UserAgent     string    `json:"user_agent"`
	WebPath       string    `json:"web_path"`
}
