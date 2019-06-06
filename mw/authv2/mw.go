package auth

import (
	"github.com/gin-gonic/gin"
	"git.corp.chaolian360.com/lrf123456/GinMW/hook"
	"github.com/go-redis/redis"
	"fmt"
	"strings"
	"github.com/json-iterator/go"
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

// 根据http请求获取session
type IGetSession interface {
	GetSession(c *gin.Context) (session string, err error)
}

type GetSession struct {
}

func (*GetSession) GetSession(c *gin.Context) (session string, err error) {
	sessid, err := c.Cookie("PHPSESSID")
	if err != nil {
		return "", ErrNoSessionId
	} else {
		return sessid, nil
	}
}

// 根据session获取用户权限
type IGetUserInfo interface {
	GetUserInfo(sessionID string) (*User, error)
}

type UserInfo struct {
	client *redis.Client
}

func NewUserInfo(client *redis.Client) (*UserInfo, error) {
	if err := client.Ping().Err(); err != nil {
		return nil, ErrRedisDisConnect
	}
	return &UserInfo{client: client}, nil
}

// 获取需要的部分user数据
func (c *UserInfo) GetUserInfo(phpsessid string) (*User, error) {
	mss, err := c.client.HGetAll(phpsessid).Result()
	if err != nil {
		return nil, ErrRedisData
	}
	ru, err := DecodeRedisUser(mss, []int{
		PickSelfId,
		PickUserId,
		PickRole,
		PickAllPermission,
	})
	if err != nil {
		return nil, ErrRedisData
	}
	return &User{
		ru.SelfId,
		ru.UserId,
		ru.Role,
		ru.AllPermission,
	}, nil
}

// 根据请求编码url标识
type IEncodeUrl interface {
	EncodeUrl(c *gin.Context) (url string)
}

type EncodeUrlWithUrl struct {
}

func (*EncodeUrlWithUrl) EncodeUrl(c *gin.Context) (url string) {
	return fmt.Sprintf("%s", c.Request.URL.Path)
}

type EncodeUrlWithUrlMethod struct {
}

func (*EncodeUrlWithUrlMethod) EncodeUrl(c *gin.Context) (url string) {
	return fmt.Sprintf("%s:%s", c.Request.URL.Path, strings.ToLower(c.Request.Method))
}

// 根据url获取标识符
type ITagUrl interface {
	GetTagByUrl(url string) (tag string, err error)
}

type TagUrlFromRedis struct {
	client *redis.Client
}

func NewTagUrlFromRedis(client *redis.Client) (*TagUrlFromRedis, error) {
	if err := client.Ping().Err(); err != nil {
		return nil, ErrRedisDisConnect
	}
	return &TagUrlFromRedis{client: client}, nil
}

func (t *TagUrlFromRedis) GetTagByUrl(url string) (tag string, err error) {
	val, err := t.client.Get("sqlNode").Result()
	if err != nil {
		return "", ErrRedisData
	} else {
		vMap := map[string]string{}
		err = jsoniter.UnmarshalFromString(val, &vMap)
		if err != nil {
			return "", ErrRedisData
		}
		for k, v := range vMap {
			if url == v {
				return k, nil
			}
		}
		return "", ErrRedisData
	}
}

// 日志打印接口
type ILog interface {
	AccessLog(user *User, session string)
	UnAccessLog(user *User, session string)
}

// 接入控制
type IAccessController interface {
	GetSession(c *gin.Context) (session string, err error) // 根据http请求获取session
	EncodeUrl(c *gin.Context) (url string)                 // 编码url
	GetTagByUrl(url string) (tag string)                   // 根据url获取表示值
	GetUserInfo(sessionID string) (*User, error)           // 根据session获取用户信息
	AccessLog(user *User, session string)                  // 打印接入日志
	UnAccessLog(user *User, session string)                // 打印禁止接入日志
	UnAuthResponse(u *User, err error) interface{}         // 控制禁止接入时返回的json数据
}

type MWAccessControlCenter struct {
	ginHook *hook.GinHook
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
		session, err := iac.GetSession(c.GetGinContext())
		if err != nil {
			// 获取session失败
			// 接入失败
			iac.UnAccessLog(nil, session)
			return nil, ErrNoSessionId
		}
		
		url := iac.EncodeUrl(c.GetGinContext())
		
		tag := iac.GetTagByUrl(url)
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
