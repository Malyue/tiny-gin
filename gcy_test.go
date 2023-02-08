package gcy

import (
	"fmt"
	"testing"
)

func newGCY() {
	c := New()
	// addRouter
	c.GET("/get", func(context *Context) {
		fmt.Fprint(context.Writer, "get successfully", context.Req.URL)
	})

	c.POST("/post", func(ctx *Context) {
		fmt.Fprint(ctx.Writer, "post successfully", ctx.Req.URL)
	})

	c.PUT("/put", func(ctx *Context) {
		fmt.Fprint(ctx.Writer, "put successfully", ctx.Req.URL)
	})

	c.DELETE("/delete", func(ctx *Context) {
		fmt.Fprint(ctx.Writer, "delete successfully", ctx.Req.URL)
	})

	c.OPTIONS("/options", func(ctx *Context) {
		fmt.Fprint(ctx.Writer, "options successfully", ctx.Req.URL)
	})

	c.Any("/any", func(ctx *Context) {
		fmt.Fprint(ctx.Writer, ctx.Method, "any successfully", ctx.Req.URL)
	})
	c.Run(":8080")
}

func TestRun(t *testing.T) {
	newGCY()
}
