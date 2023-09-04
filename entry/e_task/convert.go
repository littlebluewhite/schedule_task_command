package e_task

func ToPub(t Task) (tp TaskPub) {
	tp.TaskId = t.TaskId
	tp.Token = t.Token
	tp.From = t.From
	tp.To = t.To
	tp.TriggerFrom = t.TriggerFrom
	tp.TriggerAccount = t.TriggerAccount
	tp.Status = t.Status
	tp.Message = t.Message
	tp.TemplateID = t.TemplateID
	tp.Template = t.Template
	return
}
