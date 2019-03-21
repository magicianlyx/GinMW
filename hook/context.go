package hook

import (
	"time"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"bytes"
	"net/http"
)

// 钩子字典只读接口 返回的结果必须是副本 防止中间件内部数据被污染
type IHookContextRead interface {
	GetHook(key string) (value interface{}, exists bool)
	GetHookString(key string) (s string)
	GetHookBool(key string) (b bool)
	GetHookInt(key string) (i int)
	GetHookInt64(key string) (i64 int64)
	GetHookFloat64(key string) (f64 float64)
	GetHookTime(key string) (t time.Time)
	GetHookDuration(key string) (d time.Duration)
	GetHookStringSlice(key string) (ss []string)
	GetHookStringMap(key string) (sm map[string]interface{})
	GetHookStringMapString(key string) (sms map[string]string)
	GetHookStringMapStringSlice(key string) (smss map[string][]string)
}

// 钩子字典只写接口
type IHookContextWrite interface {
	SetHook(key string, value interface{})
}

// 钩子字典访问接口
type IHookContext interface {
	IHookContextRead
	IHookContextWrite
}

// gin.context 接口
type IGinContext interface {
	GetGinContext() *gin.Context
	GetResponseInfo() (*HttpResponse, error)
	Restore(hr *HttpResponse) error
	GetRequestInfo() (*HttpRequest, error)
}

type IHttpContext interface {
	IHookContextRead
	IHookContextWrite
	IGinContext
}

type HttpContext struct {
	hookContext map[string]interface{} // 创建一个gin以外的context 不污染gin context中的字典
	ginc        *gin.Context
}

func newHttpContext(c *gin.Context) *HttpContext {
	return &HttpContext{nil, c}
}

func (c *HttpContext) SetHook(key string, value interface{}) {
	if c.hookContext == nil {
		c.hookContext = make(map[string]interface{})
	}
	c.hookContext[key] = value
}

func (c *HttpContext) GetHook(key string) (value interface{}, exists bool) {
	value, exists = c.hookContext[key]
	return
}

func (c *HttpContext) GetHookString(key string) (s string) {
	if val, ok := c.GetHook(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}

func (c *HttpContext) GetHookBool(key string) (b bool) {
	if val, ok := c.GetHook(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

func (c *HttpContext) GetHookInt(key string) (i int) {
	if val, ok := c.GetHook(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}

func (c *HttpContext) GetHookInt64(key string) (i64 int64) {
	if val, ok := c.GetHook(key); ok && val != nil {
		i64, _ = val.(int64)
	}
	return
}

func (c *HttpContext) GetHookFloat64(key string) (f64 float64) {
	if val, ok := c.GetHook(key); ok && val != nil {
		f64, _ = val.(float64)
	}
	return
}

func (c *HttpContext) GetHookTime(key string) (t time.Time) {
	if val, ok := c.GetHook(key); ok && val != nil {
		t, _ = val.(time.Time)
	}
	return
}

func (c *HttpContext) GetHookDuration(key string) (d time.Duration) {
	if val, ok := c.GetHook(key); ok && val != nil {
		d, _ = val.(time.Duration)
	}
	return
}

func (c *HttpContext) GetHookStringSlice(key string) (ss []string) {
	if val, ok := c.GetHook(key); ok && val != nil {
		ss, _ = val.([]string)
	}
	return
}

func (c *HttpContext) GetHookStringMap(key string) (sm map[string]interface{}) {
	if val, ok := c.GetHook(key); ok && val != nil {
		sm, _ = val.(map[string]interface{})
	}
	return
}

func (c *HttpContext) GetHookStringMapString(key string) (sms map[string]string) {
	if val, ok := c.GetHook(key); ok && val != nil {
		sms, _ = val.(map[string]string)
	}
	return
}

func (c *HttpContext) GetHookStringMapStringSlice(key string) (smss map[string][]string) {
	if val, ok := c.GetHook(key); ok && val != nil {
		smss, _ = val.(map[string][]string)
	}
	return
}

func (c *HttpContext) GetGinContext() *gin.Context {
	return c.ginc
}

func (hc *HttpContext) GetResponseInfo() (*HttpResponse, error) {
	c := hc.GetGinContext()
	rw := c.Writer.(*responseWriter)
	proto := c.Request.Proto
	header := rw.Header()
	body := rw.GetBody()
	status := rw.Status()
	return &HttpResponse{
		proto,
		header,
		body,
		status,
		c,
	}, nil
}

func (hc *HttpContext) Restore(hr *HttpResponse) error {
	c := hc.GetGinContext()
	
	// body
	_, err := c.Writer.Write(hr.body)
	if err != nil {
		return ErrGinWriterInvalid
	}
	
	// header
	for k, vals := range hr.header {
		for _, v := range vals {
			c.Writer.Header().Set(k, v)
		}
	}
	
	// status
	c.Status(hr.status)
	return nil
}

func (hc *HttpContext) GetRequestInfo() (*HttpRequest, error) {
	c := hc.GetGinContext()
	rw := initResponseWrite(c.Writer)
	c.Writer = rw
	
	proto := c.Request.Proto
	method := c.Request.Method
	host := c.Request.Host
	path := c.Request.URL.Path
	
	body, err := c.GetRawData()
	defer func() {
		c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
	}()
	if err != nil {
		return nil, ErrGinRequestData
	}
	
	header := http.Header{}
	for k, vals := range c.Request.Header {
		for _, v := range vals {
			header.Set(k, v)
		}
	}
	
	return &HttpRequest{
		proto,
		method,
		host,
		path,
		header,
		body,
	}, nil
}
