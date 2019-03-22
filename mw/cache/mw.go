package cache

import (
	"github.com/gin-gonic/gin"
	"git.corp.chaolian360.com/lrf123456/GinMW/hook"
)

type MWCache struct {
	ginHook *hook.GinHook
}

func NewMWCache(cache IMWCache, serializer ISerializer, eh hook.ErrorHandler) (*MWCache, error) {
	if cache == nil || serializer == nil {
		return nil, ErrParameter
	}
	
	bh := func(c hook.IHttpContext) (error, error) {
		req, err := c.GetRequestInfo()
		if err != nil {
			return err, nil
		}
		
		if requestUID, err := serializer.RequestUUID(req); err != nil {
			return err, nil
		} else {
			
			// 没有获取到缓存
			if data, err := cache.Get(requestUID); err != nil {
				c.SetHook("requestUID", requestUID)
				return err, nil
			} else {
				if hri, err := serializer.DeserializeResponse(data); err != nil {
					// 解码http响应失败 从缓存中删除
					cache.Del(requestUID)
					c.SetHook("requestUID", requestUID)
					return err, nil
				} else {
					// 解码成功 还原到context上
					if err = c.Restore(hri); err != nil {
						// 还原失败 从缓存中删除
						cache.Del(requestUID)
						c.SetHook("requestUID", requestUID)
						return err, nil
					} else {
						// 还原成功
						c.GetGinContext().Abort()
						return nil, nil
					}
				}
			}
			
		}
	}
	
	ah := func(c hook.IHttpContext) (error, error) {
		requestUID := ""
		if v, ok := c.GetHook("requestUID"); ok {
			requestUID = v.(string)
		} else {
			// 获取缓存必要参数失败
			return ErrGetGinRequestUID, nil
		}
		
		hri, err := c.GetResponseInfo()
		if err != nil {
			// 获取response报文失败
			return err, nil
		}
		data, err := serializer.SerializeResponse(hri)
		if err != nil {
			// 序列化失败
			return err, nil
		}
		
		err = cache.Set(requestUID, data)
		if err != nil {
			// 设置缓存失败
			return err, nil
		}
		return nil, nil
	}
	
	// 任何致命错误皆恢复
	fh := func(c hook.IHttpContext, err error) error {
		return nil
	}
	
	mwc := &MWCache{hook.NewGinHook(bh, ah, fh, eh)}
	
	return mwc, nil
}

func (hc *MWCache) HandlerFunc() gin.HandlerFunc {
	return hc.ginHook.HandlerFunc()
}
