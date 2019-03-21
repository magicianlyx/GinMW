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
}

type HttpResponse struct {
	proto  string
	header http.Header
	body   []byte
	status int
	c      *gin.Context
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
