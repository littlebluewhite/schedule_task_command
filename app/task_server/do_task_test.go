package task_server

import (
	"fmt"
	"github.com/littlebluewhite/schedule_task_command/entry/e_task_template"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetStages(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		ts := []e_task_template.StageItem{
			{StageNumber: 2, Name: "a", Mode: e_task_template.Execute},
			{StageNumber: 5, Name: "b", Mode: e_task_template.Execute},
			{StageNumber: 1, Name: "c", Mode: e_task_template.Monitor},
			{StageNumber: 2, Name: "d", Mode: e_task_template.Monitor},
			{StageNumber: 1, Name: "e", Mode: e_task_template.Monitor},
			{StageNumber: 1, Name: "f", Mode: e_task_template.Execute},
			{StageNumber: 1, Name: "g", Mode: e_task_template.Monitor},
			{StageNumber: 3, Name: "h", Mode: e_task_template.Monitor},
			{StageNumber: 1, Name: "i", Mode: e_task_template.Execute},
			{StageNumber: 6, Name: "j", Mode: e_task_template.Monitor},
		}
		gsr := getTaskStage(ts)
		fmt.Println(gsr.stageMap)
		//fmt.Printf("gsr: %+v", gsr)
		require.Contains(t, gsr.sns, int32(1))
		require.Contains(t, gsr.sns, int32(2))
		require.Contains(t, gsr.sns, int32(3))
		require.Contains(t, gsr.sns, int32(6))
		require.Contains(t, gsr.sns, int32(5))
	})
}
