package cache

import (
	"testing"
)

func TestGinMW(t *testing.T) {
	
	redis, err := InitRedisCache(3, "127.0.0.1", 6379, 0, "", 10)
	if err != nil {
		panic(err)
	}
	
	memory := InitMemCache(3)
	if err != nil {
		panic(err)
	}
	
	// 使用
	mw := &Serializer{}
	
	
	hcmw, _ := NewMWCache(redis, mw)
}
