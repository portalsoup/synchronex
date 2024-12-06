package execution

type Job[J any] interface {
	Validate() bool
	Execute()
	DifferencesFromState(state J) *J
	HashCode() (uint32, error)
}
