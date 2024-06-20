package shell

type Shell interface {
	ExecuteCommand(cmd string)
}
