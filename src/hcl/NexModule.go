package hcl

import (
	"log"
)

type NexModule struct {
	// "put" copy if not present
	// "sync" unconditionally replace with current version
	// "remove" delete the file at the dest
	Path string `hcl:"type,label"`
}

func (m NexModule) Executor(context NexContext) NexModuleExecutor {
	log.Printf("context inside module executor: %s", context)
	return NexModuleExecutor{
		Context: context,
		Path:    m.Path,
	}
}

func (m NexModuleExecutor) Validate() {
	found := FindNexes(m.Path)
	parsed := ParseNexFiles(&m.Context, found)

	for _, n := range parsed {
		n.Executor(*n.Context).Validate()
	}
}

type NexModuleExecutor struct {
	NexModule NexModule
	Context   NexContext
	Path      string
}

func (m NexModuleExecutor) Run() {
	found := FindNexes(m.Path)
	parsed := ParseNexFiles(&m.Context, found)

	for _, n := range parsed {
		log.Printf("About to dereference a context: %s", m)
		n.Executor(m.Context).Run()
	}
}
