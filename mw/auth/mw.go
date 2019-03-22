package auth

import (
	"github.com/gin-gonic/gin"
	"git.corp.chaolian360.com/lrf123456/GinMW/hook"
	"github.com/go-redis/redis"
	"fmt"
)

type UserInfo struct {
	SessionId string
	User
}

type IMWAccessControl interface {
	IUserInfoRead
	Validator(c hook.IHttpContext, u *UserInfo) error
	AccessLog(u *UserInfo)
	UnAccessLog(u *UserInfo, err error)
	UnAuthResponse(u *UserInfo, err error) interface{}
}

type MWAccessControl struct {
}

func (mwac *MWAccessControl) GetUserInfo(sessid string) (*User, error) {
	cli := redis.NewClient(&redis.Options{
		Addr:     "192.168.2.241:6379",
		Password: "",
		DB:       2,
		PoolSize: 10,
	})
	rds := InitRedis(cli)
	return rds.GetUserInfo(sessid)
}

func (mwac *MWAccessControl) Validator(c hook.IHttpContext, u *UserInfo) error {
	return nil
}

func (mwac *MWAccessControl) AccessLog(u *UserInfo) {
	fmt.Printf("用户%d接入成功\r\n", u.SelfId)
}

func (mwac *MWAccessControl) UnAccessLog(u *UserInfo, err error) {
	fmt.Printf("用户%s接入失败\r\n", u.SessionId)
}

func (mwac *MWAccessControl) UnAuthResponse(u *UserInfo, err error) interface{} {
	return struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}{
		400,
		"can not access",
	}
}

type MWAccessControlCenter struct {
	ginHook *hook.GinHook
}

func NewMWAccessControlCenter(mwac IMWAccessControl) (*MWAccessControlCenter) {
	
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
		u, err := mwac.GetUserInfo(sessid)
		if err != nil || u == nil {
			// 无权限访问
			return nil, ErrRedisData
		} else {
			ui.User = u.Clone()
		}
		
		// 外部接口处理
		err = mwac.Validator(c, ui)
		if err != nil {
			// 无权限访问
			return nil, err
		}
		
		// 超管身份 完全开放访问
		if u.Role == "admin" {
			mwac.AccessLog(ui)
			return nil, nil
		}
		
		// 超管身份 完全开放访问
		if u.UserId == 1 {
			mwac.AccessLog(ui)
			return nil, nil
		}
		
		// 普通商家具有访问该页面的权限
		if ok := authCheck(c.GetGinContext().Request.URL.Path, u.AllPermission); ok {
			mwac.AccessLog(ui)
			return nil, nil
		} else {
			// 无权限访问
			return nil, ErrNoAuth
		}
		
	}
	
	// 处理被拦截的http请求
	fh := func(c hook.IHttpContext, err error) error {
		u := getUserInfo(c)
		unAuthResponse := mwac.UnAuthResponse(u, err)
		c.GetGinContext().JSON(200, unAuthResponse)
		mwac.UnAccessLog(u, err)
		return err
	}
	
	mwacc := &MWAccessControlCenter{}
	mwacc.ginHook = hook.NewGinHook(bh, nil, fh, nil)
	return mwacc
}

func (mwacc *MWAccessControlCenter) HandlerFunc() gin.HandlerFunc {
	return mwacc.ginHook.HandlerFunc()
}

func getUserInfo(c hook.IHookContextRead) *UserInfo {
	u, ok := c.GetHook("user")
	if !ok {
		return nil
	}
	return u.(*UserInfo)
}
