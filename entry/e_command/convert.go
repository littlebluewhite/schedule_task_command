package e_command

func ToPub(c Command) (cp CommandPub) {
	cp.CommandId = c.CommandId
	cp.Token = c.Token
	cp.From = c.From
	cp.To = c.To
	cp.TriggerFrom = c.TriggerFrom
	cp.TriggerAccount = c.TriggerAccount
	cp.StatusCode = c.StatusCode
	cp.RespData = c.RespData
	cp.Status = c.Status
	cp.Message = c.Message
	cp.TemplateID = c.TemplateID
	cp.Template = c.Template
	return
}
