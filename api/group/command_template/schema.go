package command_template

type SendCommand struct {
	TriggerFrom    []string `json:"trigger_from" example:"[task execute]"`
	TriggerAccount string   `json:"trigger_account" example:"Wilson"`
	Token          string   `json:"token"`
}
