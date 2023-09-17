package time_template

type CheckTime struct {
	TriggerFrom    []string `json:"trigger_from"`
	TriggerAccount string   `json:"trigger_account"`
	Token          string   `json:"token"`
}
