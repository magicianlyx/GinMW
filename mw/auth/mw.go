package auth

import (
	"github.com/go-redis/redis"
	"github.com/gin-gonic/gin"
	"ginHook/hook"
)

type ErrorLogHandler func(msg string)
type HijackHandler func(c *gin.Context)

type MWAccessControl struct {
	ginHook *hook.GinHook
	rds     *redisClient
	elh     ErrorLogHandler // 错误日志打印方法
	hjh     HijackHandler   // 拦截的非法访问http请求处理
}

func NewMWAccessControl(client *redis.Client, elh ErrorLogHandler, hjh HijackHandler) *MWAccessControl {
	rds := InitRedis(client)
	
	mwac := &MWAccessControl{hook.NewGinHook(), rds, elh, hjh}
	
	mwac.ginHook.AddBeforeHandle(func(c *hook.HttpContext) (error, error) {
		sessid, err := c.GinContext.Cookie("PHPSESSID")
		if err != nil {
			return nil, ErrSessionId
		}
		
		// 从redis中获取请求的身份信息
		u, err := rds.GetUserInfo(sessid)
		if err != nil {
			return nil, ErrRedisData
		}
		
		// 超管身份 完全开放访问
		if u.Role == "admin" {
			return nil, nil
		}
		
		// 超管身份 完全开放访问
		if u.UserId == 1 {
			return nil, nil
		}
		
		// 普通商家具有访问该页面的权限
		if ok := authCheck(c.GinContext.Request.URL.Path, u.AllPermission); ok {
			return nil, nil
		}
		
		// 无权限访问
		return nil, ErrNoAuth
	})
	
	// 处理被拦截的http请求
	mwac.ginHook.AddFailHandlerFunc(func(c *hook.HttpContext, err error) {
		mwac.hjh(c.GinContext)
	})
	
	mwac.ginHook.AddErrorHandlerFunc(func(c *hook.HttpContext, err error) {
		mwac.elh(err.Error())
	})
	
	return mwac
}

func (mwac *MWAccessControl) HandlerFunc() gin.HandlerFunc {
	return mwac.ginHook.HandlerFunc()
}
