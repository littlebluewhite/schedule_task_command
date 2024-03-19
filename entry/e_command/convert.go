package e_command

func ToPub(c Command) (cp CommandPub) {
	m := ""
	if c.Message != nil {
		m = c.Message.Error()
	}
	cp.ID = c.ID
	cp.Token = c.Token
	cp.From = c.From
	cp.To = c.To
	cp.Variables = c.Variables
	cp.Source = c.Source
	cp.TriggerFrom = c.TriggerFrom
	cp.TriggerAccount = c.TriggerAccount
	cp.StatusCode = c.StatusCode
	cp.RespData = c.RespData
	cp.Status = c.Status
	cp.ClientMessage = c.ClientMessage
	cp.Message = m
	cp.TemplateID = c.TemplateId
	cp.CommandData = c.CommandData
	return
}

func ToPubSlice(cs []Command) []CommandPub {
	cps := make([]CommandPub, 0, len(cs))
	for _, c := range cs {
		cps = append(cps, ToPub(c))
	}
	return cps
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
