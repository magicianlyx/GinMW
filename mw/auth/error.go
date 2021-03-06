package auth

import (
	"errors"
	"git.corp.chaolian360.com/lrf123456/GinMW/hook"
)

var (
	ErrGinRequestData   = hook.ErrGinRequestData
	ErrGinWriterInvalid = hook.ErrGinWriterInvalid
	ErrRedisData        = errors.New("redis internal data error")
	ErrNoSessionId      = errors.New("can not get session id from request cookie")
	ErrNoAuth           = errors.New("no authorization to access")
	ErrNoUser           = errors.New("can not get user from http context")
)
