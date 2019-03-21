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
	// unauthResponse interface{} // 没有权限访问时会返回这个结构体的json到客户端
}

func NewMWAccessControl(client *redis.Client, validator Validator, elh ErrorLogHandler, unauthResponse interface{}) (*MWAccessControl) {
	rds := InitRedis(client)
	
	bh := func(c hook.IHttpContext) (error, error) {
		sessid, err := c.GetGinContext().Cookie("PHPSESSID")
		if err != nil {
			return nil, ErrNoSessionId
		}
		
		// 从redis中获取请求的身份信息
		u, err := rds.GetUserInfo(sessid)
		if err != nil {
			return nil, ErrRedisData
		}
		
		// 外部接口处理
		err = validator(c.GetGinContext(), u)
		if err != nil {
			return nil, err
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
		if ok := authCheck(c.GetGinContext().Request.URL.Path, u.AllPermission); ok {
			return nil, nil
		}
		
		// 无权限访问
		return nil, ErrNoAuth
		
	}
	
	// 处理被拦截的http请求
	fh := func(c hook.IHttpContext, err error) error {
		// hjh(c.GetGinContext())
		c.GetGinContext().JSON(200, unauthResponse)
		return err
	}
	
	// 处理错误节点
	eh := func(c hook.IHookContextRead, err error, isDeadly bool) {
		elh(err, isDeadly)
	}
	
	mwac := &MWAccessControl{hook.NewGinHook(bh, nil, fh, eh), rds}
	
	return mwac
}

func (mwac *MWAccessControl) HandlerFunc() gin.HandlerFunc {
	return mwac.ginHook.HandlerFunc()
}
