package e_task

func ToPub(t Task) (tp TaskPub) {
	tp.TaskId = t.TaskId
	tp.Token = t.Token
	tp.From = t.From
	tp.To = t.To
	tp.TriggerFrom = t.TriggerFrom
	tp.TriggerAccount = t.TriggerAccount
	tp.Status = t.Status
	tp.Message = t.Message.Error()
	tp.TemplateID = t.TemplateId
	tp.Template = t.Template
	return
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
