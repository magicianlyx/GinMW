package hook

import (
	"github.com/gin-gonic/gin"
)

type IMW interface {
	HandlerFunc() gin.HandlerFunc
}

type GinHook struct {
	bhm *BeforeHandleMap
	ahm *AfterHandleMap
	fhm *FailHandlerMap  // 致命错误 一般会处理http的响应结果
	ehm *ErrorHandlerMap // 所有错误 一般实现逻辑是打印日志
}

func NewGinHook() *GinHook {
	return &GinHook{
		NewBeforeHandleMap(),
		NewAfterHandleMap(),
		NewFailHandlerMap(),
		NewErrorHandlerMap(),
	}
}

func (gh *GinHook) AddFailHandlerFunc(fh FailHandler) {
	gh.fhm.Add(fh)
}
func (gh *GinHook) DelFailHandlerFunc(fh FailHandler) {
	gh.fhm.Del(fh)
}

func (gh *GinHook) AddErrorHandlerFunc(fh ErrorHandler) {
	gh.ehm.Add(fh)
}
func (gh *GinHook) DelErrorHandlerFunc(fh ErrorHandler) {
	gh.ehm.Del(fh)
}

func (gh *GinHook) AddBeforeHandle(bh BeforeHandle) {
	gh.bhm.Add(bh)
}

func (gh *GinHook) DelBeforeHandle(bh BeforeHandle) {
	gh.bhm.Del(bh)
}

func (gh *GinHook) AddAfterHandle(ah AfterHandle) {
	gh.ahm.Add(ah)
}

func (gh *GinHook) DelAfterHandle(ah AfterHandle) {
	gh.ahm.Del(ah)
}

func (gh *GinHook) HandlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		hc := NewHttpContext(c)
		
		gh.bhm.m.Range(func(_, value interface{}) bool {
			bh := value.(BeforeHandle)
			e1, e2 := bh(hc)
			if e1 != nil {
				// 非致命错误
				gh.ehm.InvokeAll(hc, e1, false)
				return true
			}
			if e2 != nil {
				// 致命错误
				gh.fhm.InvokeAll(hc, e2)
				gh.ehm.InvokeAll(hc, e2, true)
				c.Abort()
				return false
			}
			return true
		})
		
		c.Next()
		
		gh.ahm.m.Range(func(_, value interface{}) bool {
			ah := value.(AfterHandle)
			e1, e2 := ah(hc)
			if e1 != nil {
				// 非致命错误
				gh.ehm.InvokeAll(hc, e1, false)
			}
			if e2 != nil {
				// 致命错误
				gh.fhm.InvokeAll(hc, e2)
				gh.ehm.InvokeAll(hc, e2, true)
				c.Abort()
				return false
			}
			return true
		})
		
	}
}
