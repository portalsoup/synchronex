package hcl

import (
	"path/filepath"
)

type NexModule struct {
	// "put" copy if not present
	// "sync" unconditionally replace with current version
	// "remove" delete the file at the dest
	Path string `hcl:"type,label"`
}

func (m NexModule) Executor(context NexContext) NexModuleExecutor {
	path, _ := filepath.Abs(m.Path)
	context.Path = path
	return NexModuleExecutor{
		Context:   context,
		NexModule: m,
	}
}

func (m NexModuleExecutor) Validate() {
	found := FindNexes(m.Context.Path)
	parsed := ParseNexFiles(&m.Context, found)

	for _, n := range parsed {
		n.Executor(*n.Context).Validate()
	}
}

type NexModuleExecutor struct {
	NexModule NexModule
	Context   NexContext
}

func (m NexModuleExecutor) Run() {
	found := FindNexes(m.Context.Path)
	parsed := ParseNexFiles(&m.Context, found)

	for _, n := range parsed {
		n.Executor(m.Context).Run()
	}
}
