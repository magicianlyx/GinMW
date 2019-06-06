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
	ErrNoUrlFormat      = errors.New("can not get url format from http context")
	ErrNoAuth           = errors.New("no authorization to access")
	ErrNoUser           = errors.New("can not get user from http context")
	ErrUnknownUrl       = errors.New("unknown tag of url")
	ErrRedisDisConnect  = errors.New("redis disconnect")
	// ErrRedisReadFail    = errors.New("can not read necessary  data from redis") // 用ErrRedisData代替
)
