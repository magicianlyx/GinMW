package auth

import (
	"github.com/go-redis/redis"
)

type User struct {
	SelfId        int
	UserId        int
	Role          string
	AllPermission []string
}

// 生成副本
func (u *User) Clone() User {
	if u == nil{
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

type IUserInfoRead interface {
	GetUserInfo(string) (*User, error)
}

type redisClient struct {
	client *redis.Client
}

func InitRedis(cli *redis.Client) (*redisClient) {
	return &redisClient{cli}
}

// 获取需要的部分user数据
func (c *redisClient) GetUserInfo(phpsessid string) (*User, error) {
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
