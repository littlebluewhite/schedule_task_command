package task_template

type SendTask struct {
	TriggerFrom    []string          `json:"trigger_from" example:"[task execute]"`
	TriggerAccount string            `json:"trigger_account" example:"Wilson"`
	Token          string            `json:"token"`
	Variables      map[string]string `json:"variables"`
}
