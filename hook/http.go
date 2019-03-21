package hook

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"bytes"
	"time"
)

// type IHttpContextRead interface {
// 	Get(key string) (value interface{}, exists bool)
// 	GetString(key string) (s string)
// 	GetBool(key string) (b bool)
// 	GetInt(key string) (i int)
// 	GetInt64(key string) (i64 int64)
// 	GetFloat64(key string) (f64 float64)
// 	GetTime(key string) (t time.Time)
// 	GetDuration(key string) (d time.Duration)
// 	GetStringSlice(key string) (ss []string)
// 	GetStringMap(key string) (sm map[string]interface{})
// 	GetStringMapString(key string) (sms map[string]string)
// 	GetStringMapStringSlice(key string) (smss map[string][]string)
// }
//
// type IHttpContextWrite interface {
// 	Set(key string, value interface{})
// }
//
// type IHttpContext interface {
// 	IHttpContextRead
// 	IHttpContextWrite
// }

type HttpContext struct {
	context    map[string]interface{} // 创建一个gin以外的context 不污染gin context中的字典
	GinContext *gin.Context
}

func NewHttpContext(c *gin.Context) *HttpContext {
	return &HttpContext{nil, c}
}

func (c *HttpContext) Set(key string, value interface{}) {
	if c.context == nil {
		c.context = make(map[string]interface{})
	}
	c.context[key] = value
}

func (c *HttpContext) Get(key string) (value interface{}, exists bool) {
	value, exists = c.context[key]
	return
}

func (c *HttpContext) GetString(key string) (s string) {
	if val, ok := c.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}

func (c *HttpContext) GetBool(key string) (b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

func (c *HttpContext) GetInt(key string) (i int) {
	if val, ok := c.Get(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}

func (c *HttpContext) GetInt64(key string) (i64 int64) {
	if val, ok := c.Get(key); ok && val != nil {
		i64, _ = val.(int64)
	}
	return
}

func (c *HttpContext) GetFloat64(key string) (f64 float64) {
	if val, ok := c.Get(key); ok && val != nil {
		f64, _ = val.(float64)
	}
	return
}

func (c *HttpContext) GetTime(key string) (t time.Time) {
	if val, ok := c.Get(key); ok && val != nil {
		t, _ = val.(time.Time)
	}
	return
}

func (c *HttpContext) GetDuration(key string) (d time.Duration) {
	if val, ok := c.Get(key); ok && val != nil {
		d, _ = val.(time.Duration)
	}
	return
}

func (c *HttpContext) GetStringSlice(key string) (ss []string) {
	if val, ok := c.Get(key); ok && val != nil {
		ss, _ = val.([]string)
	}
	return
}

func (c *HttpContext) GetStringMap(key string) (sm map[string]interface{}) {
	if val, ok := c.Get(key); ok && val != nil {
		sm, _ = val.(map[string]interface{})
	}
	return
}

func (c *HttpContext) GetStringMapString(key string) (sms map[string]string) {
	if val, ok := c.Get(key); ok && val != nil {
		sms, _ = val.(map[string]string)
	}
	return
}

func (c *HttpContext) GetStringMapStringSlice(key string) (smss map[string][]string) {
	if val, ok := c.Get(key); ok && val != nil {
		smss, _ = val.(map[string][]string)
	}
	return
}

type HttpRequest struct {
	Proto  string
	Method string
	Host   string
	Path   string
	Header http.Header
	Body   []byte
}

func GetRequestInfo(hc *HttpContext) (*HttpRequest, error) {
	c := hc.GinContext
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

type HttpResponse struct {
	proto  string
	header http.Header
	body   []byte
	status int
	c      *gin.Context
}

func GetResponseInfo(hc *HttpContext) (*HttpResponse, error) {
	c := hc.GinContext
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

func (hr *HttpResponse) AddHeader(k string, v string) {
	hr.c.Writer.Header().Set(k, v)
}

func (hr *HttpResponse) Write(body []byte) (int, error) {
	return hr.c.Writer.Write(hr.body)
}

func (hr *HttpResponse) Status(code int) {
	hr.c.Status(code)
}

func (hr *HttpResponse) Restore(hc *HttpContext) error {
	c := hc.GinContext
	
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
