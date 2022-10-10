package quetaro

type ErrWithJob struct {
	cause error
	Id    string
	Name  string
}

func (e *ErrWithJob) Error() string {
	return e.cause.Error()
}

func (e *ErrWithJob) Unwrap() error {
	return e.cause
}
