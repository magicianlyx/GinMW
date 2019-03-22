package GinMW

import (
	"github.com/gin-gonic/gin"
	"git.corp.chaolian360.com/lrf123456/GinMW/mw/cache"
	"git.corp.chaolian360.com/lrf123456/GinMW/hook"
	"github.com/go-redis/redis"
	"git.corp.chaolian360.com/lrf123456/GinMW/mw/auth"
	"fmt"
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
	
	cli := redis.NewClient(&redis.Options{
		Addr:     "192.168.2.241:6379",
		Password: "",
		DB:       2,
		PoolSize: 10,
	})
	rds := auth.InitRedis(cli)
	
	acc := auth.NewMWAccessControl(
		rds,
		nil,
		func(u *auth.UserInfo) {
			fmt.Printf("用户%d接入成功\r\n", u.SelfId)
		},
		func(u *auth.UserInfo, err error) {
			fmt.Printf("用户%s接入失败\r\n", u.SessionId)
		},
		func(u *auth.UserInfo, err error) interface{} {
			return struct {
				Code int    `json:"code"`
				Msg  string `json:"msg"`
			}{
				400,
				"can not access",
			}
		},
	)
	
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
