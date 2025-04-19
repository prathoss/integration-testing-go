package domain

var _ error = (*ErrNotFound)(nil)

type ErrNotFound struct {
	Msg string
}

func (e ErrNotFound) Error() string {
	return e.Msg
}

var _ error = (*ErrInvalid)(nil)

type ErrInvalid struct {
	Msg string
}

func (e ErrInvalid) Error() string {
	return e.Msg
}
