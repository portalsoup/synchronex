package hcl

import (
	"synchronex/src/hcl/context"
)

type NexModule struct {
	// "put" copy if not present
	// "sync" unconditionally replace with current version
	// "remove" delete the file at the dest
	Path string `hcl:"type,label"`
}

func (m NexModule) Executor(context context.NexContext) NexModuleExecutor {

	return NexModuleExecutor{
		Context: context,
		Path:    m.Path,
	}
}

func (m NexModuleExecutor) Validate() {
	found := FindNexes(m.Path)
	parsed := ParseNexFiles(found)

	for _, n := range parsed {
		n.Validate()
	}
}

type NexModuleExecutor struct {
	NexModule NexModule
	Context   context.NexContext
	Path      string
}

func (m NexModuleExecutor) Run() {
	found := FindNexes(m.Path)
	parsed := ParseNexFiles(found)

	for _, n := range parsed {
		n.Executor(m.Context).Run()
	}
}
