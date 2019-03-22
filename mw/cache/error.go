package cache

import (
	"errors"
	"git.corp.chaolian360.com/lrf123456/GinMW/hook"
)

var (
	ErrGinRequestData   = hook.ErrGinRequestData
	ErrGinWriterInvalid = hook.ErrGinWriterInvalid
	ErrParameter        = errors.New("parameter invalid")
	ErrJsonMarshal      = errors.New("marshal json fail")
	ErrJsonUnmarshal    = errors.New("unmarshal json fail")
	ErrGetGinRequestUID = errors.New("can not get parameter `request_uid` from gin's context structure")
	ErrCacheNoRecord    = errors.New("cache key not found")
)
