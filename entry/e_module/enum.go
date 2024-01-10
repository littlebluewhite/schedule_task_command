package e_module

type Module int

const (
	None Module = iota
	Command
	Task
)

func (m Module) String() string {
	return [...]string{"", "command", "task"}[m]
}
