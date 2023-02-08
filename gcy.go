package gcy

import (
	"net/http"
)

// HandlerFunc defines the handler used by gcy
type HandlerFunc func(c *Context)

// Engine defines the struct to Scheduling resource
type Engine struct {
	router *Router
}

// New export the engine init to user
func New() *Engine {
	return &Engine{
		router: newRouter(),
	}
}

// GET defines the HTTP Get request
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.router.addRoute("GET", pattern, handler)
}

// POST defines the HTTP Post request
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.router.addRoute("POST", pattern, handler)
}

// PUT defines the HTTP put request
func (engine *Engine) PUT(pattern string, handler HandlerFunc) {
	engine.router.addRoute("PUT", pattern, handler)
}

// DELETE defines the HTTP delete request
func (engine *Engine) DELETE(pattern string, handler HandlerFunc) {
	engine.router.addRoute("DELETE", pattern, handler)
}

// OPTIONS defines the HTTP options request
func (engine *Engine) OPTIONS(pattern string, handler HandlerFunc) {
	engine.router.addRoute("OPTIONS", pattern, handler)
}

// Any defines the all HTTP request method
func (engine *Engine) Any(pattern string, handler HandlerFunc) {
	engine.router.addRoute("GET", pattern, handler)
	engine.router.addRoute("POST", pattern, handler)
	engine.router.addRoute("PUT", pattern, handler)
	engine.router.addRoute("DELETE", pattern, handler)
	engine.router.addRoute("OPTIONS", pattern, handler)
}

// Run implements the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// ServeHTTP implements the http.ListenAndServe handler
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	engine.router.handle(newContext(w, req))
}
