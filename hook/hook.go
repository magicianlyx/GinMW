package hook

import (
	"github.com/gin-gonic/gin"
)

type IMW interface {
	HandlerFunc() gin.HandlerFunc
}

type GinHook struct {
	bh BeforeHandle
	ah AfterHandle
	fh FailHandler  // 致命错误 一般会处理http的响应结果
	eh ErrorHandler // 所有错误 一般实现逻辑是打印日志
}

func NewGinHook(bh BeforeHandle, ah AfterHandle, fh FailHandler, eh ErrorHandler) *GinHook {
	
	if bh == nil {
		bh = func(c IHttpContext) (error, error) {
			return nil, nil
		}
	}
	if ah == nil {
		ah = func(c IHttpContext) (error, error) {
			return nil, nil
		}
	}
	
	if fh == nil {
		fh = func(c IHttpContext, err error) error {
			return err
		}
	}
	
	if eh == nil {
		eh = func(c IHookContextRead, err error, isDeadly bool) {
		}
	}
	
	return &GinHook{
		bh,
		ah,
		fh,
		eh,
	}
}

func (gh *GinHook) HandlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		hc := newHttpContext(c)
		
		e1, e2 := gh.bh(hc)
		if e1 != nil {
			// 非致命错误
			gh.eh(hc, e1, false)
		}
		if e2 != nil {
			// 致命错误
			gh.eh(hc, e2, true)
			
			if !gh.fh(hc, e2) {
				c.Abort()
			} else {
				// 致命错误恢复
				// 执行下一个gin节点
			}
		}
		
		c.Next()
		
		e1, e2 = gh.ah(hc)
		if e1 != nil {
			// 非致命错误
			gh.eh(hc, e1, false)
		}
		if e2 != nil {
			// 致命错误
			gh.eh(hc, e2, true)
			
			if !gh.fh(hc, e2) {
				c.Abort()
			} else {
				// 致命错误恢复
				// 执行下一个gin节点
			}
		}
		
	}
}
