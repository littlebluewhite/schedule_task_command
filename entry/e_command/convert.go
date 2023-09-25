package e_command

func ToPub(c Command) (cp CommandPub) {
	m := ""
	if c.Message != nil {
		m = c.Message.Error()
	}
	cp.CommandId = c.CommandId
	cp.Token = c.Token
	cp.From = c.From
	cp.To = c.To
	cp.TriggerFrom = c.TriggerFrom
	cp.TriggerAccount = c.TriggerAccount
	cp.StatusCode = c.StatusCode
	cp.RespData = c.RespData
	cp.Status = c.Status
	cp.Message = m
	cp.TemplateID = c.TemplateId
	cp.Template = c.Template
	return
}

func S2Status(s *string) Status {
	if s == nil {
		return Prepared
	}
	switch *s {
	case "Process":
		return Process
	case "Success":
		return Success
	case "Failure":
		return Failure
	case "Cancel":
		return Cancel
	default:
		return Prepared
	}
}
