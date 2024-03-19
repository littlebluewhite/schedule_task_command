package time_template

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
	"path/filepath"
	"runtime"
	"schedule_task_command/app/dbs"
	"schedule_task_command/app/time_server"
	"schedule_task_command/entry/e_time_data"
	"schedule_task_command/entry/e_time_template"
	"schedule_task_command/util/config"
	"schedule_task_command/util/logFile"
	"testing"
	"time"
)

func setUpOperate() (o hOperate, l logFile.LogFile) {
	l = logFile.NewLogFile("test", "Operate.log")
	_, b, _, _ := runtime.Caller(0)
	rootPath := filepath.Dir(filepath.Dir(filepath.Dir(b)))
	c := config.NewConfig(rootPath, "config", "config", config.Yaml)
	DBS := dbs.NewDbs(l, true, c)
	timeS := time_server.NewTimeServer(DBS)
	o = NewOperate(DBS, timeS)
	return
}

func TestQuery(t *testing.T) {
	o, l := setUpOperate()
	t.Run("test find", func(t *testing.T) {
		l.Info().Println("tset Operate time template list")
		var i int32 = 20
		st1 := datatypes.NewTime(5, 12, 12, 0)
		st2 := datatypes.NewTime(8, 12, 12, 0)
		st3 := datatypes.NewTime(10, 12, 12, 0)
		st4 := datatypes.NewTime(20, 12, 12, 0)
		testTimeTemplates := []*e_time_template.TimeTemplateCreate{
			{Name: "test1",
				TimeData: e_time_data.TimeDatumCreate{
					RepeatType:      e_time_data.S2RepeatType(nil),
					StartDate:       time.Date(2023, 6, 18, 0, 0, 0, 0, time.Local),
					StartTime:       &st1,
					EndTime:         datatypes.NewTime(16, 9, 16, 0),
					IntervalSeconds: &i,
					ConditionType:   e_time_data.S2ConditionType(nil),
					TCondition:      []byte("[5, 1, 7]"),
				},
			},
			{Name: "test2",
				TimeData: e_time_data.TimeDatumCreate{
					RepeatType:      e_time_data.S2RepeatType(nil),
					StartDate:       time.Date(2023, 6, 18, 0, 0, 0, 0, time.Local),
					StartTime:       &st2,
					EndTime:         datatypes.NewTime(16, 9, 16, 0),
					IntervalSeconds: &i,
					ConditionType:   e_time_data.S2ConditionType(nil),
					TCondition:      []byte("[5, 1, 7]"),
				},
			},
			{Name: "test3",
				TimeData: e_time_data.TimeDatumCreate{
					RepeatType:      e_time_data.S2RepeatType(nil),
					StartDate:       time.Date(2023, 6, 18, 0, 0, 0, 0, time.Local),
					StartTime:       &st3,
					EndTime:         datatypes.NewTime(16, 9, 16, 0),
					IntervalSeconds: &i,
					ConditionType:   e_time_data.S2ConditionType(nil),
					TCondition:      []byte("[5, 1, 7]"),
				},
			},
			{Name: "test4",
				TimeData: e_time_data.TimeDatumCreate{
					RepeatType:      e_time_data.S2RepeatType(nil),
					StartDate:       time.Date(2023, 6, 18, 0, 0, 0, 0, time.Local),
					StartTime:       &st4,
					EndTime:         datatypes.NewTime(16, 9, 16, 0),
					IntervalSeconds: &i,
					ConditionType:   e_time_data.S2ConditionType(nil),
					TCondition:      []byte("[5, 1, 7]"),
				},
			},
		}
		timeTemplates, err := o.Create(testTimeTemplates)
		require.Nil(t, err)
		tIds := make([]int32, 0, 4)
		for _, item := range timeTemplates {
			tIds = append(tIds, item.ID)
		}
		fmt.Println("point")
		timeTemplates2, err := o.Find(tIds)
		require.Nil(t, err)
		require.Equal(t, len(timeTemplates), 4)
		require.Equal(t, timeTemplates2[0].Name, "test1")
		require.Equal(t, timeTemplates2[1].Name, "test2")
		require.Equal(t, timeTemplates2[2].Name, "test3")
		require.Equal(t, timeTemplates2[3].Name, "test4")
	})
	t.Run("test List", func(t *testing.T) {
		l.Info().Println("test time templates list")
		_, err := o.List()
		require.Nil(t, err)
	})
}

func TestCreate(t *testing.T) {
	o, l := setUpOperate()
	t.Run("create success", func(t *testing.T) {
		l.Info().Println("test Operate time template create")
		var i int32 = 300
		st := datatypes.NewTime(12, 15, 12, 0)
		testTimeTemplate := []*e_time_template.TimeTemplateCreate{
			{Name: "test6", TimeData: e_time_data.TimeDatumCreate{
				RepeatType:      e_time_data.S2RepeatType(nil),
				StartDate:       time.Date(2023, 6, 16, 0, 0, 0, 0, time.Local),
				StartTime:       &st,
				EndTime:         datatypes.NewTime(13, 21, 13, 0),
				IntervalSeconds: &i,
				ConditionType:   e_time_data.S2ConditionType(nil),
				TCondition:      []byte("[1, 7, 3, 4]"),
			}},
		}
		result, err := o.Create(testTimeTemplate)
		fmt.Println(result)
		require.Nil(t, err)
		require.Equal(t, result[0].Name, "test6")
	})
	t.Run("create fail", func(t *testing.T) {

		l.Info().Println("test Operate time template create")
		var i int32 = 300
		st := datatypes.NewTime(8, 12, 12, 0)
		testTimeTemplate := []*e_time_template.TimeTemplateCreate{
			{Name: "test6", TimeData: e_time_data.TimeDatumCreate{
				RepeatType:      e_time_data.S2RepeatType(nil),
				StartDate:       time.Date(2023, 6, 19, 0, 0, 0, 0, time.Local),
				StartTime:       &st,
				EndTime:         datatypes.NewTime(13, 9, 13, 0),
				IntervalSeconds: &i,
				ConditionType:   e_time_data.S2ConditionType(nil),
				TCondition:      []byte("[1, 8, 3, 4]"),
			}},
		}
		result, err := o.Create(testTimeTemplate)
		fmt.Println(result)
		require.Nil(t, result)
		require.Error(t, err)
	})
}

func TestUpdate(t *testing.T) {
	o, l := setUpOperate()
	t.Run("update", func(t *testing.T) {
		var s = "monthly_day"
		name := "test1"
		startTime := datatypes.NewTime(8, 12, 12, 0)
		testTimeTemplate := []*e_time_template.TimeTemplateUpdate{
			{Name: &name, ID: 1,
				TimeData: &e_time_data.TimeDatumUpdate{
					RepeatType:      e_time_data.S2RepeatType(nil),
					StartDate:       time.Date(2023, 6, 18, 0, 0, 0, 0, time.Local),
					StartTime:       &startTime,
					EndTime:         datatypes.NewTime(16, 55, 16, 0),
					IntervalSeconds: nil,
					ConditionType:   e_time_data.S2ConditionType(&s),
					TCondition:      []byte("[5, 1, 7]"),
				},
			},
		}
		l.Info().Println("test Operate time template update")
		err := o.Update(testTimeTemplate)
		require.Nil(t, err)
	})
}

func TestDelete(t *testing.T) {
	o, l := setUpOperate()
	t.Run("delete", func(t *testing.T) {
		var i int32 = 20
		st1 := datatypes.NewTime(5, 12, 12, 0)
		st2 := datatypes.NewTime(8, 12, 12, 0)
		st3 := datatypes.NewTime(10, 12, 12, 0)
		st4 := datatypes.NewTime(20, 12, 12, 0)
		testTimeTemplates := []*e_time_template.TimeTemplateCreate{
			{Name: "apple",
				TimeData: e_time_data.TimeDatumCreate{
					RepeatType:      e_time_data.S2RepeatType(nil),
					StartDate:       time.Date(2023, 6, 18, 0, 0, 0, 0, time.Local),
					StartTime:       &st1,
					EndTime:         datatypes.NewTime(16, 9, 16, 0),
					IntervalSeconds: &i,
					ConditionType:   e_time_data.S2ConditionType(nil),
					TCondition:      []byte("[5, 1, 7]"),
				},
			},
			{Name: "dog",
				TimeData: e_time_data.TimeDatumCreate{
					RepeatType:      e_time_data.S2RepeatType(nil),
					StartDate:       time.Date(2023, 6, 18, 0, 0, 0, 0, time.Local),
					StartTime:       &st2,
					EndTime:         datatypes.NewTime(16, 9, 16, 0),
					IntervalSeconds: &i,
					ConditionType:   e_time_data.S2ConditionType(nil),
					TCondition:      []byte("[5, 1, 7]"),
				},
			},
			{Name: "banana",
				TimeData: e_time_data.TimeDatumCreate{
					RepeatType:      e_time_data.S2RepeatType(nil),
					StartDate:       time.Date(2023, 6, 18, 0, 0, 0, 0, time.Local),
					StartTime:       &st3,
					EndTime:         datatypes.NewTime(16, 9, 16, 0),
					IntervalSeconds: &i,
					ConditionType:   e_time_data.S2ConditionType(nil),
					TCondition:      []byte("[5, 1, 7]"),
				},
			},
			{Name: "cherry",
				TimeData: e_time_data.TimeDatumCreate{
					RepeatType:      e_time_data.S2RepeatType(nil),
					StartDate:       time.Date(2023, 6, 18, 0, 0, 0, 0, time.Local),
					StartTime:       &st4,
					EndTime:         datatypes.NewTime(16, 9, 16, 0),
					IntervalSeconds: &i,
					ConditionType:   e_time_data.S2ConditionType(nil),
					TCondition:      []byte("[5, 1, 7]"),
				},
			},
		}
		l.Info().Println("test Operate time template delete")
		timeTemplates, err := o.Create(testTimeTemplates)
		require.Nil(t, err)
		ids := make([]int32, 0, len(timeTemplates))
		for _, tt := range timeTemplates {
			ids = append(ids, tt.ID)
		}
		err = o.Delete(ids)
		require.Nil(t, err)
	})
}
