package hook

import (
	"github.com/gin-gonic/gin"
)

type responseWriter struct {
	body []byte
	gin.ResponseWriter
}

func initResponseWrite(rw gin.ResponseWriter) *responseWriter {
	return &responseWriter{[]byte{}, rw}
}

func (w *responseWriter) Write(data []byte) (int, error) {
	w.body = make([]byte, len(data))
	copy(w.body, data)
	return w.ResponseWriter.Write(data)
}

func (w *responseWriter) WriteString(s string) (int, error) {
	w.body = make([]byte, len([]byte(s)))
	copy(w.body, []byte(s))
	return w.ResponseWriter.WriteString(s)
}

func (w *responseWriter) Status() int {
	return w.ResponseWriter.Status()
}

func (w *responseWriter) Written() bool {
	return w.ResponseWriter.Written()
}

func (w *responseWriter) GetBody() []byte {
	return w.body
}
