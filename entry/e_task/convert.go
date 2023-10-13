package e_task

func ToPub(t Task) (tp TaskPub) {
	m := ""
	if t.Message != nil {
		m = t.Message.Error()
	}
	tp.TaskId = t.TaskId
	tp.Token = t.Token
	tp.From = t.From
	tp.To = t.To
	tp.Variables = t.Variables
	tp.TriggerFrom = t.TriggerFrom
	tp.TriggerAccount = t.TriggerAccount
	tp.Status = t.Status
	tp.Stages = t.Stages
	tp.ClientMessage = t.ClientMessage
	tp.Message = m
	tp.TemplateID = t.TemplateId
	tp.TaskData = t.TaskData
	return
}

func ToPubSlice(ts []Task) []TaskPub {
	tps := make([]TaskPub, 0, len(ts))
	for _, t := range ts {
		tps = append(tps, ToPub(t))
	}
	return tps
}

func S2Status(s *string) TStatus {
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
