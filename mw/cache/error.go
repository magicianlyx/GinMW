package cache

import (
	"errors"
	"ginHook/hook"
)

var (
	ErrGinRequestData   = hook.ErrGinRequestData
	ErrGinWriterInvalid = hook.ErrGinWriterInvalid
	ErrJsonMarshal      = errors.New("marshal json fail")
	ErrJsonUnmarshal    = errors.New("unmarshal json fail")
)
