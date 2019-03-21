package auth

import (
	"github.com/go-redis/redis"
	"github.com/gin-gonic/gin"
	"git.corp.chaolian360.com/lrf123456/GinMW/hook"
)

type ErrorLogHandler func(err error, isDeadly bool)
type HijackHandler func(c *gin.Context)
type Validator func(c *gin.Context, u *User) error

type MWAccessControl struct {
	ginHook *hook.GinHook
	rds     *redisClient
	// validator Validator   // 提供给外部的验证器 通过该接口能够获取http请求的用户身份信息
	// elh     ErrorLogHandler // 错误日志打印方法
	// hjh     HijackHandler   // 拦截的非法访问http请求处理
}

func NewMWAccessControl(client *redis.Client, validator Validator, elh ErrorLogHandler, hjh HijackHandler) (*MWAccessControl) {
	rds := InitRedis(client)
	
	mwac := &MWAccessControl{hook.NewGinHook(), rds}
	
	mwac.ginHook.AddBeforeHandle(func(c *hook.HttpContext) (error, error) {
		sessid, err := c.GinContext.Cookie("PHPSESSID")
		if err != nil {
			return nil, ErrNoSessionId
		}
		
		// 从redis中获取请求的身份信息
		u, err := rds.GetUserInfo(sessid)
		if err != nil {
			return nil, ErrRedisData
		}
		
		err = validator(c.GinContext, u)
		if err != nil {
			return nil, err
		}
		c.Set("user", u)
		return nil, nil
	})
	
	mwac.ginHook.AddBeforeHandle(func(c *hook.HttpContext) (error, error) {
		
		v, ok := c.Get("user")
		if !ok {
			return nil, ErrNoUser
		}
		u, ok := v.(*User)
		if !ok {
			return nil, ErrNoUser
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
		hjh(c.GinContext)
	})
	
	// 处理错误节点
	mwac.ginHook.AddErrorHandlerFunc(func(c *hook.HttpContext, err error, isDeadly bool) {
		elh(err, isDeadly)
	})
	
	return mwac
}

// func GetUserFromContext(c *gin.Context) (*User, bool) {
// 	if v, ok := c.Get("user"); ok {
// 		if u, ok1 := v.(*User); ok1 {
// 			return u, true
// 		}
// 	}
// 	return nil, false
// }

func (mwac *MWAccessControl) HandlerFunc() gin.HandlerFunc {
	return mwac.ginHook.HandlerFunc()
}
