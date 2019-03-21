package cache


import (
	"time"
	"github.com/patrickmn/go-cache"
	"github.com/go-redis/redis"
	"strconv"
	"errors"
)

type ICache interface {
	Set(key string, v []byte) (error)
	Get(key string) ([]byte, error)
	Del(key string)
}

type RedisCache struct {
	validTime int
	cli       *redis.Client
}

func InitRedisCache(second int, host string, port int, db int, password string, poolsize int) (*RedisCache, error) {
	rc := &RedisCache{}
	rc.validTime = second
	opt := &redis.Options{
		Addr:     host + ":" + strconv.Itoa(port),
		DB:       db,
		Password: password,
		PoolSize: poolsize,
	}
	rc.cli = redis.NewClient(opt)
	if _, err := rc.cli.Ping().Result(); err != nil {
		return nil, err
	}
	return rc, nil
}

func (rc *RedisCache) Set(key string, v []byte) (error) {
	_, err := rc.cli.Set(key, string(v), time.Duration(rc.validTime)*time.Second).Result()
	return err
}

func (rc *RedisCache) Get(key string) ([]byte, error) {
	res, err := rc.cli.Get(key).Result()
	if err != nil {
		return nil, err
	} else {
		return []byte(res), nil
	}
}

func (rc *RedisCache) Del(key string) {
	rc.cli.Del(key)
}

type MemCache struct {
	cache *cache.Cache
}

func InitMemCache(second int) (*MemCache) {
	validTime := time.Second * time.Duration(second)
	mc := &MemCache{}
	mc.cache = cache.New(validTime, validTime)
	return mc
}

func (mc *MemCache) Set(key string, v []byte) (error) {
	mc.cache.SetDefault(key, v)
	return nil
}

func (mc *MemCache) Get(key string) ([]byte, error) {
	v, ok := mc.cache.Get(key)
	if ok {
		return v.([]byte), nil
	} else {
		return nil, errors.New("record not found")
	}
}

func (mc *MemCache) Del(key string) {
	mc.cache.Delete(key)
}
