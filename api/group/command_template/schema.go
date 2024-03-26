package command_template

type SendCommand struct {
	TriggerFrom    []map[string]string `json:"trigger_from"`
	TriggerAccount string              `json:"trigger_account" example:"Wilson"`
	Token          string              `json:"token"`
	Variables      map[string]string   `json:"variables"`
}
