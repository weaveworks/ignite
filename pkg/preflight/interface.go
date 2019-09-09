package preflight

type Checker interface {
	Check() error
	Name() string
	Type() string
}
