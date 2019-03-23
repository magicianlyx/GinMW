package GinMW

import (
	"github.com/gin-gonic/gin"
	"git.corp.chaolian360.com/lrf123456/GinMW/mw/cache"
	"git.corp.chaolian360.com/lrf123456/GinMW/hook"
	"time"
	"testing"
)

func TestA(t *testing.T) {
	r := gin.Default()

	c := cache.NewMemCache(10)
	mw := &cache.Serializer{}
	hcmw, _ := cache.NewMWCache(c, mw, func(c hook.IHookContextRead, err error, isDeadly bool) {
	
	})

	r.GET("/wxsdk",
		hcmw.HandlerFunc(),
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
