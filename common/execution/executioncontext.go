package execution

type ExecutionContext struct {
	ConfigurationPhase []*Job
	ExecutionPhase     []*Job
}

func (ec *ExecutionContext) AddConfiguration(job *Job) {
	ec.ConfigurationPhase = append(ec.ConfigurationPhase, job)
}

func (ec *ExecutionContext) AddExecution(job *Job) {
	ec.ExecutionPhase = append(ec.ExecutionPhase, job)
}
