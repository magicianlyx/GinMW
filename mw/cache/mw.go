package cache

import (
	"github.com/gin-gonic/gin"
	"git.corp.chaolian360.com/lrf123456/GinMW/hook"
)

type MWCache struct {
	ginHook *hook.GinHook
}

func NewMWCache(cache ICache, serializer ISerializer) (*MWCache, error) {
	if cache == nil || serializer == nil {
		return nil, ErrParameter
	}
	
	mwc := &MWCache{hook.NewGinHook()}
	
	mwc.ginHook.AddBeforeHandle(func(c *hook.HttpContext) (error, error) {
		req, err := hook.GetRequestInfo(c)
		if err != nil {
			return err, nil
		}
		
		if requestUID, err := serializer.RequestUUID(req); err != nil {
			return err, nil
		} else {
			
			// 没有获取到缓存
			if data, err := cache.Get(requestUID); err != nil {
				c.Set("requestUID", requestUID)
				return err, nil
			} else {
				if hri, err := serializer.DeserializeResponse(data); err != nil {
					// 解码http响应失败 从缓存中删除
					cache.Del(requestUID)
					c.Set("requestUID", requestUID)
					return err, nil
				} else {
					// 解码成功 还原到context上
					if err = hri.Restore(c); err != nil {
						// 还原失败 从缓存中删除
						cache.Del(requestUID)
						c.Set("requestUID", requestUID)
						return err, nil
					} else {
						// 还原成功
						c.GinContext.Abort()
						return nil, nil
					}
				}
			}
			
		}
	})
	
	mwc.ginHook.AddAfterHandle(func(c *hook.HttpContext) (error, error) {
		requestUID := ""
		if v, ok := c.Get("requestUID"); ok {
			requestUID = v.(string)
		} else {
			// 获取缓存必要参数失败
			return ErrGetGinRequestUID, nil
		}
		
		hri, err := hook.GetResponseInfo(c)
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
		
	})
	return mwc, nil
}

func (hc *MWCache) HandlerFunc() gin.HandlerFunc {
	return hc.ginHook.HandlerFunc()
}
