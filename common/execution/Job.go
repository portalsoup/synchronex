package execution

type Job interface {
	validation() (bool, error)
	execution() error
	ToString() string
}
