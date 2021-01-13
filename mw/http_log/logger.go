package http_log

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

var space = []byte(" ")
var enter = []byte("\r")
var feed = []byte("\n")

func Log(log *log.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Writer = initResponseWrite(c.Writer)
		RequestLog(log, c)
		c.Next()
		ResponseLog(log, c)
	}

}

func RequestLog(log *log.Logger, c *gin.Context) {
	method := c.Request.Method
	url := c.Request.URL.String()
	proto := c.Request.Proto
	header := c.Request.Header
	body, _ := c.GetRawData()

	data := mergeBytes(
		[]byte(method), space, []byte(url), space, []byte(proto), enter, feed,
		mergeHeader(header), feed,
		body,
	)

	log.Writer().Write(data)
}

func ResponseLog(log *log.Logger, c *gin.Context) {
	rw := c.Writer.(*responseWriter)

	proto := c.Request.Proto
	status := strconv.Itoa(rw.Status())
	statusText := http.StatusText(rw.Status())

	header := rw.Header()
	body := rw.GetBody()

	data := mergeBytes(
		[]byte(proto), space, []byte(status), space, []byte(statusText), enter, feed,
		mergeHeader(header), feed,
		body,
	)

	log.Writer().Write(data)
}

func mergeBytes(datas ...[]byte) []byte {
	r := make([]byte, 0)
	for i := range datas {
		r = append(r, datas[i]...)
	}
	return r
}

func mergeHeader(headers http.Header) []byte {
	r := make([]byte, 0)
	for k := range headers {
		t := []byte(fmt.Sprintf("%s:%s%s%s", k, headers[k], enter, feed))
		r = append(r, t...)
	}
	return r
}
