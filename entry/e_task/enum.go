package e_task

import "github.com/goccy/go-json"

type TStatus int

const (
	Prepared TStatus = iota
	Process
	Success
	Failure
	Cancel
)

func (s TStatus) String() string {
	return [...]string{"Prepared", "Process", "Success", "Failure", "Cancel"}[s]
}

func (s TStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *TStatus) UnmarshalJSON(data []byte) error {
	var tStatus string
	err := json.Unmarshal(data, &tStatus)
	if err != nil {
		return err
	}
	*s = S2Status(&tStatus)
	return nil
}
