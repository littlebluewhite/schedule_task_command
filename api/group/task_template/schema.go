package task_template

type SendTask struct {
	TriggerFrom    []map[string]string       `json:"trigger_from"`
	TriggerAccount string                    `json:"trigger_account" example:"Wilson"`
	Token          string                    `json:"token"`
	Variables      map[int]map[string]string `json:"variables"`
}
