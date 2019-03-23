package hook

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type HttpRequest struct {
	Proto  string
	Method string
	Host   string
	Path   string
	Header http.Header
	Body   []byte
	c      *gin.Context
}

func (hr *HttpRequest) Query(key string) string {
	return hr.c.Query(key)
}

func (hr *HttpRequest) Cookie(name string) (string, error) {
	return hr.c.Cookie(name)
}

type HttpResponse struct {
	Proto      string
	Header     http.Header
	Body       []byte
	StatusCode int
	c          *gin.Context
}

func (hr *HttpResponse) AddHeader(k string, v string) {
	hr.c.Writer.Header().Set(k, v)
}

func (hr *HttpResponse) Write(body []byte) (int, error) {
	return hr.c.Writer.Write(hr.Body)
}

func (hr *HttpResponse) Status(code int) {
	hr.c.Status(code)
}
