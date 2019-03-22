package GinMW

import (
	"github.com/gin-gonic/gin"
	"git.corp.chaolian360.com/lrf123456/GinMW/mw/cache"
	"git.corp.chaolian360.com/lrf123456/GinMW/hook"
	"git.corp.chaolian360.com/lrf123456/GinMW/mw/auth"
	"time"
	"testing"
)

func TestA(t *testing.T) {
	r := gin.Default()
	
	c, err := cache.InitRedisCache(10, "127.0.0.1", 6379, 0, "", 10)
	if err != nil {
		panic(err)
	}
	
	mw := &cache.Serializer{}
	hcmw, _ := cache.NewMWCache(c, mw, func(c hook.IHookContextRead, err error, isDeadly bool) {
	
	})
	
	acc := auth.NewMWAccessControl()
	
	r.GET("/wxsdk",
		hcmw.HandlerFunc(),
		acc.HandlerFunc(),
		wxsdk,
	)
	
	r.Run(":8081")
	// _ = hcmw
}

func wxsdk(c *gin.Context) {
	c.JSON(200, gin.H{
		"now": time.Now().Format("2006-01-02 15:04:05"),
	})
}
