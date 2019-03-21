package auth

import (
	"strconv"
	"github.com/json-iterator/go"
)

const (
	PickStrPermission = iota
	PickUserId
	PickSelfId
	PickButton
	PickSystem
	PickRole
	PickNav
	PickAllPermission
)

var (
	pickUserMap = map[int]func(map[string]string, *RedisUser) error{
		PickStrPermission: decodeStrPermission,
		PickUserId:        decodeUserId,
		PickSelfId:        decodeSelfId,
		PickButton:        decodeButton,
		PickSystem:        decodeSystem,
		PickRole:          decodeRole,
		PickNav:           decodeNav,
		PickAllPermission: decodeAllPermission,
	}
)

type MenuEntry struct {
	Sign     string      `json:"sign"`
	Name     string      `json:"name"`
	Children []MenuEntry `json:"children,omitempty"`
}

type RedisUser struct {
	StrPermission string              `mapstructure:"strPermission"`
	UserId        int                 `mapstructure:"user_id"`
	SelfId        int                 `mapstructure:"self_id"`
	Button        map[string][]string `mapstructure:"button"`
	System        []string            `mapstructure:"system"`
	Role          string              `mapstructure:"role"`
	Nav           []MenuEntry         `mapstructure:"nav"`
	AllPermission []string            `mapstructure:"allPermission"`
}

func decodeStrPermission(mss map[string]string, u *RedisUser) error {
	if u == nil {
		u = &RedisUser{}
	}
	var err error
	u.StrPermission = mss["strPermission"]
	return err
}

func decodeUserId(mss map[string]string, u *RedisUser) error {
	if u == nil {
		u = &RedisUser{}
	}
	var err error
	u.UserId, err = strconv.Atoi(mss["user_id"])
	return err
}
func decodeSelfId(mss map[string]string, u *RedisUser) error {
	if u == nil {
		u = &RedisUser{}
	}
	var err error
	u.SelfId, err = strconv.Atoi(mss["self_id"])
	return err
}

func decodeButton(mss map[string]string, u *RedisUser) error {
	if u == nil {
		u = &RedisUser{}
	}
	var err error
	u.Button = map[string][]string{}
	err = jsoniter.UnmarshalFromString(mss["button"], &u.Button)
	return err
}

func decodeSystem(mss map[string]string, u *RedisUser) error {
	if u == nil {
		u = &RedisUser{}
	}
	var err error
	u.System = []string{}
	err = jsoniter.UnmarshalFromString(mss["system"], &u.System)
	return err
}

func decodeRole(mss map[string]string, u *RedisUser) error {
	if u == nil {
		u = &RedisUser{}
	}
	var err error
	
	u.Role = mss["role"]
	return err
}

func decodeNav(mss map[string]string, u *RedisUser) error {
	if u == nil {
		u = &RedisUser{}
	}
	var err error
	
	u.Nav = []MenuEntry{}
	err = jsoniter.UnmarshalFromString(mss["nav"], &u.Nav)
	return err
}

func decodeAllPermission(mss map[string]string, u *RedisUser) error {
	if u == nil {
		u = &RedisUser{}
	}
	var err error
	
	u.AllPermission = []string{}
	err = jsoniter.UnmarshalFromString(mss["allPermission"], &u.AllPermission)
	return err
}

func DecodeRedisUser(mss map[string]string, needPick []int) (*RedisUser, error) {
	u := &RedisUser{}
	for _, mark := range needPick {
		if err := pickUserMap[mark](mss, u); err != nil {
			return nil, err
		}
	}
	return u, nil
}

