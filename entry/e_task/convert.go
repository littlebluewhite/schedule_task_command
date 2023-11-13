package e_task

import "slices"

func ToPub(t Task) (tp TaskPub) {
	m := ""
	if t.Message != nil {
		m = t.Message.Error()
	}
	tp.ID = t.ID
	tp.Token = t.Token
	tp.From = t.From
	tp.To = t.To
	tp.Variables = t.Variables
	tp.TriggerFrom = t.TriggerFrom
	tp.TriggerAccount = t.TriggerAccount
	tp.Status = t.Status
	tp.StageNumber = t.StageNumber
	tp.Stages = t.Stages
	tp.FailedCommands = t.FailedCommands
	tp.ClientMessage = t.ClientMessage
	tp.Message = m
	tp.TemplateID = t.TemplateId
	tp.TaskData = t.TaskData
	return
}

func ToSimpleTask(t Task) (ts SimpleTask) {
	var si []StageItem
	var sns []int32
	for sn := range t.Stages {
		sns = append(sns, sn)
	}
	slices.Sort(sns)
	for _, stageNumber := range sns {
		si = append(si, t.Stages[stageNumber].Monitor...)
		si = append(si, t.Stages[stageNumber].Execute...)
	}
	ts.ID = t.ID
	ts.Status = int(t.Status)
	ts.StageNumber = t.StageNumber
	ts.StageItems = si
	return
}

func ToStageItemStatus(t Task) (r []int) {
	var si []StageItem
	var sns []int32
	for sn := range t.Stages {
		sns = append(sns, sn)
	}
	slices.Sort(sns)
	for _, stageNumber := range sns {
		si = append(si, t.Stages[stageNumber].Monitor...)
		si = append(si, t.Stages[stageNumber].Execute...)
	}
	for _, s := range si {
		r = append(r, int(s.Status))
	}
	return
}

func ToPubSlice(ts []Task) []TaskPub {
	tps := make([]TaskPub, 0, len(ts))
	for _, t := range ts {
		tps = append(tps, ToPub(t))
	}
	return tps
}

func ToSimpleTaskSlice(ts []Task) []SimpleTask {
	tss := make([]SimpleTask, 0, len(ts))
	for _, t := range ts {
		tss = append(tss, ToSimpleTask(t))
	}
	return tss
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
