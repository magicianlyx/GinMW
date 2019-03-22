package auth

import (
	"github.com/gin-gonic/gin"
	"git.corp.chaolian360.com/lrf123456/GinMW/hook"
)

// 外伸接口的userinfo参数是副本 跟中间件内部的不是同一个对象 在外面接口修改userinfo内的数据不会影响中间件内部的处理

type Validator func(c *hook.HttpRequest, u *UserInfo) error  // 返回err!=nil时会无权访问
type UnAuthResponse func(u *UserInfo, err error) interface{} // 无权限访问时返回给客户端的json response

type AccessLog func(u *UserInfo)              // 接入回调 用于日志打印
type UnAccessLog func(u *UserInfo, err error) // 无权限接入回调 用于日志打印

type MWAccessControl struct {
	ginHook *hook.GinHook
	rds     IUserInfoRead
}

type UserInfo struct {
	SessionId string
	User
}

func NewMWAccessControl(client IUserInfoRead, validator Validator, al AccessLog, ual UnAccessLog, uar UnAuthResponse) (*MWAccessControl) {
	if validator == nil {
		validator = func(c *hook.HttpRequest, u *UserInfo) error {
			return nil
		}
	}
	if al == nil {
		al = func(u *UserInfo) {
		}
	}
	if ual == nil {
		ual = func(u *UserInfo, err error) {
		}
	}
	if uar == nil {
		uar = func(u *UserInfo, err error) interface{} {
			return nil
		}
	}
	
	bh := func(c hook.IHttpContext) (error, error) {
		
		ui := &UserInfo{"", User{}}
		
		// 存储http请求用户信息到hook上下文
		defer c.SetHook("user", ui)
		
		sessid, err := c.GetGinContext().Cookie("PHPSESSID")
		if err != nil {
			return nil, ErrNoSessionId
		} else {
			ui.SessionId = sessid
		}
		
		// 从redis中获取请求的身份信息
		u, err := client.GetUserInfo(sessid)
		if err != nil || u == nil {
			// 无权限访问
			ual(ui, err)
			return nil, ErrRedisData
		} else {
			ui.User = u.Clone()
		}
		
		// 从gin context中获取http request信息
		hri, err := c.GetRequestInfo()
		if err != nil {
			// 无权限访问
			ual(ui, err)
			return nil, err
		}
		
		// 外部接口处理
		err = validator(hri, ui)
		if err != nil {
			// 无权限访问
			ual(ui, err)
			return nil, err
		}
		
		// 超管身份 完全开放访问
		if u.Role == "admin" {
			al(ui)
			return nil, nil
		}
		
		// 超管身份 完全开放访问
		if u.UserId == 1 {
			al(ui)
			return nil, nil
		}
		
		// 普通商家具有访问该页面的权限
		if ok := authCheck(c.GetGinContext().Request.URL.Path, u.AllPermission); ok {
			al(ui)
			return nil, nil
		} else {
			// 无权限访问
			ual(ui, ErrNoAuth)
			return nil, ErrNoAuth
		}
		
	}
	
	// 处理被拦截的http请求
	fh := func(c hook.IHttpContext, err error) error {
		u := getUserInfo(c)
		unAuthResponse := uar(u, err)
		c.GetGinContext().JSON(200, unAuthResponse)
		ual(u, err)
		return err
	}
	
	mwac := &MWAccessControl{hook.NewGinHook(bh, nil, fh, nil), client}
	
	return mwac
}

func (mwac *MWAccessControl) HandlerFunc() gin.HandlerFunc {
	return mwac.ginHook.HandlerFunc()
}

func getUserInfo(c hook.IHookContextRead) *UserInfo {
	u, ok := c.GetHook("user")
	if !ok {
		return nil
	}
	return u.(*UserInfo)
}
