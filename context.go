package gcy

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"
)

// abortIndex represents a typical value used in abort functions
const abortIndex int8 = math.MaxInt8 >> 1

// defines the context to store the session info
type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	// response info
	StatusCode int
	ErrorMsg   string
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Path:   r.URL.Path,
		Method: r.Method,
	}
}

// Query returns the keyed url query value if it exists otherwise it returns an empty string `("")`
func (c *Context) Query(key string) (value string) {
	return c.Req.URL.Query().Get(key)
}

// DefaultQuery returns the url query value if it exist
// otherwise returns the specified defalut value
func (c *Context) DefaultQuery(key string, defaultValue string) (value string) {
	value = c.Req.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	} else {
		return value
	}
}

// PostForm returns the specified key from a post urlencoded form or multipart form when it exist
// otherwise it returns an empty string ("")
func (c *Context) PostForm(key string) (value string) {
	return c.Req.PostForm.Get(key)
}

// DefaultPostForm
func (c *Context) DefaultPostForm(key string, defaultValue string) (value string) {
	value = c.Req.PostForm.Get(key)
	if value == "" {
		return defaultValue
	} else {
		return value
	}
}

// IsWebsocket returns whether the request headers indicate a websocket handshake is being initiated by the client
func (c *Context) IsWebsocket() bool {
	if strings.Contains(strings.ToLower(c.Req.Header.Get("Connection")), "upgrade") &&
		strings.EqualFold(c.Req.Header.Get("Upgrade"), "websocket") {
		return true
	}
	return false
}

// Status set the HTTP response code
func (c *Context) Status(code int) {
	c.Writer.WriteHeader(code)
	c.StatusCode = code
}

// Header shortcut for c.Writer.Header().Set(key,value)
// if value == "",removes the header
func (c *Context) Header(key, value string) {
	if value == "" {
		c.Writer.Header().Del(key)
		return
	}
	c.Writer.Header().Set(key, value)
}

// String write the string into the response body
func (c *Context) String(code int, format string, values ...interface{}) {
	c.Header("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// JSON write the json into the response body
func (c Context) JSON(code int, obj interface{}) {
	c.Header("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// Data write the data into the response body
func (c Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}
