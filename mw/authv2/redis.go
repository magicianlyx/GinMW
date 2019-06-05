package auth

import (
	"github.com/go-redis/redis"
)

type RedisUserInfoRead struct {
	client *redis.Client
}

func NewRedisUserInfoRead(cli *redis.Client) (*RedisUserInfoRead) {
	return &RedisUserInfoRead{cli}
}

// 获取需要的部分user数据
func (c *RedisUserInfoRead) GetUserInfo(phpsessid string) (*User, error) {
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
