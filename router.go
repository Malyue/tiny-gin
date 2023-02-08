package gcy

import (
	"fmt"
	"log"
	"net/http"
)

// defines the router
type Router struct {
	handlers map[string]HandlerFunc
}

// init the router
func newRouter() *Router {
	return &Router{
		handlers: make(map[string]HandlerFunc),
	}
}

// Set the key and value in Engine.router
// addRoute add the route in engine
func (router *Router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	router.handlers[key] = handler
}

// implements the ServerHTTP handle
func (router *Router) handle(c *Context) {
	key := c.Req.Method + "-" + c.Req.URL.Path
	if handler, ok := router.handlers[key]; ok {
		handler(c)
	} else {
		c.Writer.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(c.Writer, "404 NOT FOUND:%s\n", c.Req.URL)
	}
}
