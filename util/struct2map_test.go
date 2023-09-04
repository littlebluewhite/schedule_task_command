package util

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"schedule_task_command/dal/model"
	"testing"
	"time"
)

func TestStructToMap1(t *testing.T) {
	var i int32 = 20
	sample := model.TimeTemplate{
		ID:   1,
		Name: "test", TimeDataID: 1,
		TimeData: model.TimeDatum{
			RepeatType:      nil,
			StartDate:       time.Date(2023, 6, 9, 0, 0, 0, 0, time.Local),
			StartTime:       []byte("13:12:12"),
			EndTime:         []byte("16:17:16"),
			IntervalSeconds: &i,
			ConditionType:   nil,
			TCondition:      []byte("[5, 6, 8]"),
		},
	}
	res := StructToMap(&sample)
	require.NotNil(t, res)
	fmt.Printf("%+v\n", res)
	jbyt, err := json.Marshal(sample)
	require.NoError(t, err)
	fmt.Println(string(jbyt))

}

func TestMapDeleteNil(t *testing.T) {
	m := map[string]interface{}{"created_at": nil, "id": 1, "updated_at": nil}
	MapDeleteNil(m)
	P(m)
	require.NotContains(t, m, "created_at")
}
