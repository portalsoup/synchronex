package hcl

type Executable interface {
	Executor(context NexContext) Executor
}
