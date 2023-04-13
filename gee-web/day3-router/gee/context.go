package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

//处理请求的req rsp 信息
type Context struct {
	Writer http.ResponseWriter
	Req    *http.Request

	Path   string
	Method string
	Params map[string]string

	StatusCode int
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Path:   r.URL.Path,
		Method: r.Method,
	}
}

func (r *Context) PostForm(k string) string {
	return r.Req.FormValue(k)
}

func (r *Context) Query(k string) string {
	return r.Req.URL.Query().Get(k)
}

func (r *Context) Param(k string) string {
	v, _ := r.Params[k]
	return v
}

func (r *Context) Status(code int) {
	r.StatusCode = code
	r.Writer.WriteHeader(code)
}

func (r *Context) SetHeader(k, v string) {
	r.Writer.Header().Set(k, v)
}

func (r *Context) String(code int, format string, values ...interface{}) {
	r.SetHeader("Content-Type", "text/plain")
	r.Status(code)
	r.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (r *Context) JSON(code int, obj interface{}) {
	r.SetHeader("Content-Type", "application/json")
	r.Status(code)
	encoder := json.NewEncoder(r.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(r.Writer, err.Error(), 500)
	}
}

func (r *Context) Data(code int, data []byte) {
	r.Status(code)
	r.Writer.Write(data)
}

func (r *Context) HTML(code int, html string) {
	r.SetHeader("Content-Type", "text/html")
	r.Status(code)
	r.Writer.Write([]byte(html))
}
