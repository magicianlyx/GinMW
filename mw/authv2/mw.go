package auth

import (
	"github.com/gin-gonic/gin"
	"git.corp.chaolian360.com/lrf123456/GinMW/hook"
)

type User struct {
	SelfId        int
	UserId        int
	Role          string
	AllPermission []string
}

// 生成副本
func (u *User) Clone() User {
	if u == nil {
		return User{}
	}
	allPermission := make([]string, len(u.AllPermission))
	copy(allPermission, u.AllPermission)
	return User{
		u.SelfId,
		u.UserId,
		u.Role,
		allPermission,
	}
}

// 根据session获取用户权限
type IGetUserInfo interface {
	GetUserInfo(sessionID string) (*User, error)
}

// 根据请求获取标识符
type IGetRequestTag interface {
	GetTagByUrl(c *gin.Context) (tag string)
}

// 日志打印接口
type ILog interface {
	AccessLog(user *User, session string)
	UnAccessLog(user *User, session string)
}

// 接入控制
type IAccessController interface {
	GetTagByUrl(c *gin.Context) (tag string)
	GetUserInfo(sessionID string) (*User, error)
	AccessLog(user *User, session string)
	UnAccessLog(user *User, session string)
	UnAuthResponse(u *User, err error) interface{}
}

type MWAccessControlCenter struct {
	ginHook *hook.GinHook
}

func getSession(c *gin.Context) (sessionID string, err error) {
	
	sessid, err := c.Cookie("PHPSESSID")
	if err != nil {
		return "", ErrNoSessionId
	} else {
		return sessid, nil
	}
}

func NewMWAccessControlCenter(iac IAccessController) (*MWAccessControlCenter) {
	if iac == nil {
		return &MWAccessControlCenter{hook.NewGinHook(nil, nil, nil, nil)}
	}
	
	getUserInfo := func(c hook.IHookContextRead) *User {
		u, ok := c.GetHook("user")
		if !ok {
			return nil
		}
		return u.(*User)
	}
	
	bh := func(c hook.IHttpContext) (error, error) {
		session, err := getSession(c.GetGinContext())
		if err != nil {
			// 获取session失败
			// 接入失败
			iac.UnAccessLog(nil, session)
			return nil, ErrNoSessionId
		}
		
		tag := iac.GetTagByUrl(c.GetGinContext())
		if tag == "" {
			// 无法获取url的tag
			// 接入失败
			iac.UnAccessLog(nil, session)
			return nil, ErrUnknownUrl
		}
		
		user, err := iac.GetUserInfo(session)
		if err != nil {
			// 获取session用户失败
			iac.UnAccessLog(user, session)
			return nil, ErrNoUser
		}
		
		// 存储http请求用户信息到hook上下文
		c.SetHook("user", user)
		
		// 用户没有权限接入
		if user.AllPermission == nil {
			iac.UnAccessLog(user, session)
			return nil, ErrNoAuth
		}
		
		access := false
		for i := range user.AllPermission {
			permission := user.AllPermission[i]
			if tag == permission {
				// 有权限接入
				access = true
				break
			}
		}
		
		if access {
			// 接入成功 打印日志
			iac.AccessLog(user, session)
			return nil, nil
		} else {
			// 接入失败
			iac.UnAccessLog(user, session)
			return nil, ErrNoAuth
		}
		
	}
	
	// 处理被拦截的http请求
	fh := func(c hook.IHttpContext, err error) error {
		u := getUserInfo(c)
		unAuthResponse := iac.UnAuthResponse(u, err)
		c.GetGinContext().JSON(200, unAuthResponse)
		return err
	}
	
	mwacc := &MWAccessControlCenter{}
	mwacc.ginHook = hook.NewGinHook(bh, nil, fh, nil)
	return mwacc
}

func (mwacc *MWAccessControlCenter) HandlerFunc() gin.HandlerFunc {
	return mwacc.ginHook.HandlerFunc()
}
