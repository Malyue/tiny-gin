## Gcy 

一个模仿gin框架实现的建议go web框架,参考了7-days-go的gee项目

首先从http请求来入手，先将http包的常用接口进行封装

首先需要一个全局统一调度资源的部分，即`Engine`,也由它来调用所需要的接口，包括初始化(`New`)，运行服务(`Run`)等,以及定义一个统一的函数类型,先从最简单的开始，即至少要有response和request

```go
// HandlerFunc defines the handler used by gcy
type HandlerFunc func(w http.ResponseWriter,req *http.Request)

// Engine defines the struct to Scheduling resource
type Engine struct {
    // define the variable you want
    router map[string]HandlerFunc
}

// New export the engine init to user
func New() *Engine {
	return &Engine{
		router: make(map[string]HandlerFunc),
	}
}

// Run implements the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// ServeHTTP implements the http.ListenAndServe handler
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler()
	} else {

	}
}
```

这里要注意，由于`Run`是对`http.ListenAndServe(addr,handler)`的封装，所以需要看他的源代码观察其参数，发现handler是一个接口，需要实现一个ServeHTTP方法，即只要传入了某个实现了ServeHTTP方法的实例，所有的http请求都会交给该实例进行处理，否则为nil的话则使用标准库中的实例来进行处理，这时就需要通过http.HandleFunc()来添加路由

> ServeHTTP方法的作用是解析路径，查找路由映射表，并执行查到的对应注册好的处理方法

这里先将此时的路由规则定义为`Method + "-" + pattern`,作为路由映射表的key，那这时候我们就要思考如何添加路由到engine的路由表里了，我们定义一个`addRoute`方法

```go
// Set the key and value in Engine.router
// addRoute add the route in engine
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}
```

为了使用的方便，这时再对`addRoute`方法进行上层的再一次封装

```go
// GET defines the HTTP Get request
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// POST defines the HTTP Post request
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// PUT defines the HTTP put request
func (engine *Engine) PUT(pattern string, handler HandlerFunc) {
	engine.addRoute("PUT", pattern, handler)
}

// DELETE defines the HTTP delete request
func (engine *Engine) DELETE(pattern string, handler HandlerFunc) {
	engine.addRoute("DELETE", pattern, handler)
}

// OPTIONS defines the HTTP options request
func (engine *Engine) OPTIONS(pattern string, handler HandlerFunc) {
	engine.addRoute("OPTIONS", pattern, handler)
}

// Any defines the all HTTP request method
func (engine *Engine) Any(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
	engine.addRoute("POST", pattern, handler)
	engine.addRoute("PUT", pattern, handler)
	engine.addRoute("DELETE", pattern, handler)
	engine.addRoute("OPTIONS", pattern, handler)
}
```

至此，一个最简单的web框架对http的封装就完成了，我们来测试一下：

```go
// gcy_test.go

func newGCY() {
	c := New()
	// addRouter
	c.GET("/get", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "get successfully", r.URL)
	})

	c.POST("/post", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "post successfully", r.URL)
	})

	c.PUT("/put", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "put successfully", r.URL)
	})

	c.DELETE("/delete", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "delete successfully", r.URL)
	})

	c.OPTIONS("/options", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "options successfully", r.URL)
	})

	c.Any("/any", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, r.Method, "any successfully", r.URL)
	})
	c.Run(":8080")
}

func TestRun(t *testing.T) {
	newGCY()
}

```

用curl命令执行一下，`curl http://localhost:8080/get`和`curl http://localhost:8080/any`

![image-20230201012354828](gcy_images\image-20230201012354828.png)

![image-20230201012459455](gcy_images\image-20230201012459455.png)

经过测试，没出现问题，至此，一个go-web框架对http最简单的封装就已经完成了.

这时候我们就该往下考虑，这个封装太简单了，如果我们要实现分组，动态路由，或者中间件等，那该如何实现？

首先要考虑一个共同的问题，即产生的信息存放在哪里？

这时候就要想起`Context---上下文`,它随着每一个请求的出现而产生，结束而销毁，和当前请求强相关的信息都应由`Context`来承载，在这里将设计结构，扩展性和复杂性都留在了内部，对外简化接口，`将一次会话的所有信息全储存了起来`

首先先定义`Context`结构，最基本的`r 和 w`肯定要有，除此之外，一次会话还需要有一个返回的结构，不过由于用户需求不同，返回的结构不同，所以这部分需要由用户自己来处理，那为了方便，还可以将`Method`和`Path`放到上下文中.

```go
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
}
```

那这时`HandleFunc`的参数类型变为`c *Context`，即需要修改之前的代码,,例如原先自定义的`ServerHTTP`,这时在这里我们顺便把路由抽象出来，在后面也需要用上.

```go
-------------
   gcy.go	|
-------------

// HandlerFunc defines the handler used by gcy
type HandlerFunc func(c *Context)

// Engine defines the struct to Scheduling resource
type Engine struct {
	router *router
}

// ServeHTTP implements the http.ListenAndServe handler
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	engine.router.handle(newContext(w,req))
}


--------------------------------------------------------

-----------------
   router.go	|
-----------------

// defines the router
type router struct {
	handlers map[string]HandlerFunc
}

// init the router
func newRouter() *router {
	return &router{
		handlers: make(map[string]HandlerFunc),
	}
}

// Set the key and value in Engine.router
// addRoute add the route in engine
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	engine.router.handlers[key] = handler
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

----------------------------------------------------

-----------------
   context.go   |
-----------------

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
	ErrorMsg string
}


func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Path:   r.URL.Path,
		Method: r.Method,
	}
}
```

我们既然已经将`Request和Writer`放到上下文中，那我们也可以提供对应封装好的上层接口,可以直接通过类似`c.Query`来调用,这里先封装常用的几个函数

```go
-----------------
   context.go   |
-----------------

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

```







