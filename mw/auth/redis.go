package auth

import (
	"github.com/go-redis/redis"
)

type redisClient struct {
	client *redis.Client
}

func InitRedis(cli *redis.Client) (*redisClient) {
	return &redisClient{cli}
}

type User struct {
	SelfId        int
	UserId        int
	Role          string
	AllPermission []string
}

// 获取需要的redis数据
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
