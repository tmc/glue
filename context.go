package glue

import (
	"net/http"

	"github.com/tmc/inj"
)

type Context struct {
	inj.Injector
	*Glue
	rw *ResponseWriter
}

func (g *Glue) newContext(w http.ResponseWriter, r *http.Request) *Context {
	ctx := &Context{inj.New(), g, NewResponseWriter(w)}

	ctx.Register(r)
	// register our ResponseWriter as an http.ResponseWriter or net/http HandlerFunc compatibility
	ctx.RegisterAs(ctx.rw, (*http.ResponseWriter)(nil))
	// register this instance with itself
	ctx.Register(*ctx)
	return ctx
}

func (ctx *Context) handle() {
	for _, h := range append(ctx.handlers, ctx.defaultHandler) {
		_, err := ctx.Call(h)
		if err != nil {
			panic(err)
		}
		if ctx.rw.WroteHeader() {
			break
		}
	}
}
