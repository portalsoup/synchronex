package shell

type Shell interface {
	Name() string
	ExecuteCommand(cmd string)
}
