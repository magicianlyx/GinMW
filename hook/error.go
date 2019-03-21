package hook

import (
	"errors"
)

var (
	ErrGinRequestData   = errors.New("can not get necessary data from http request")
	ErrGinWriterInvalid = errors.New("gin response writer invalid")
)
