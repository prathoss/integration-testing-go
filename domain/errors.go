package domain

var _ error = (*ErrNotFound)(nil)

type ErrNotFound struct {
	Msg string
}

func NewErrNotFound(msg string) *ErrNotFound {
	return &ErrNotFound{Msg: msg}
}

func (e ErrNotFound) Error() string {
	return e.Msg
}

var _ error = (*ErrInvalid)(nil)

type ErrInvalid struct {
	Msg string
}

func NewErrInvalid(msg string) *ErrInvalid {
	return &ErrInvalid{Msg: msg}
}

func (e ErrInvalid) Error() string {
	return e.Msg
}
