package cache

import (
	"github.com/json-iterator/go"
	"ginHook/hook"
)

type IRequestSerializer interface {
	RequestUUID(hri *hook.HttpRequest) (s string, err error)
}

type IResponseSerializer interface {
	SerializeResponse(hri *hook.HttpResponse) (data []byte, err error)
	DeserializeResponse(data []byte) (hri *hook.HttpResponse, err error)
}

type ISerializer interface {
	IRequestSerializer
	IResponseSerializer
}

type Serializer struct {
}


// http request 唯一性
func (*Serializer) RequestUUID(hri *hook.HttpRequest) (s string, err error) {
	return hri.Path, nil
}

// 序列化response结构体 可以根据需要只存储一部分参数
func (*Serializer) SerializeResponse(v *hook.HttpResponse) ([]byte, error) {
	bs, err := jsoniter.Marshal(v)
	if err != nil {
		return nil, ErrJsonMarshal
	} else {
		return bs, nil
	}
}

// 反序列化response结构体 注意要跟序列化函数配合
func (*Serializer) DeserializeResponse(data []byte) (hri *hook.HttpResponse, err error) {
	hri = &hook.HttpResponse{}
	err = jsoniter.Unmarshal(data, hri)
	if err != nil {
		return nil, ErrJsonUnmarshal
	} else {
		return hri, nil
	}
}
