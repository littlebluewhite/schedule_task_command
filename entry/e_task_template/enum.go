package e_task_template

import (
	"fmt"
	"github.com/goccy/go-json"
	"schedule_task_command/util"
)

type Mode int

const (
	NoneMode Mode = iota
	Monitor
	Execute
)

func (m Mode) String() string {
	return [...]string{"", "monitor", "execute"}[m]
}

func (m Mode) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.String())
}

func (m *Mode) UnmarshalJSON(data []byte) error {
	var tStatus string
	err := json.Unmarshal(data, &tStatus)
	if err != nil {
		return err
	}
	*m = S2Mode(&tStatus)
	return nil
}

func TaskTemplateNotFound(id int) util.MyErr {
	e := fmt.Sprintf("task template id: %d not found", id)
	return util.MyErr(e)
}
