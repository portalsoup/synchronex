package provisioner

import (
	"synchronex/src/hcl/context"
)

type Executable interface {
	Executor(context context.NexContext) Executor
}
