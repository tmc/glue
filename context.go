package glue

import (
	"log"
	"net/http"
	"reflect"

	"github.com/tmc/inj"
)

// Context represents the execution context for a request in Glue
// It is a DI (Dependency Injection) container and contains an augmented
// ResponseWriter
type Context struct {
	inj.Injector
	g  *Glue
	rw *ResponseWriter
}

// newContext creates a new Context and registers a few basic instances
func (g *Glue) newContext(w http.ResponseWriter, r *http.Request) *Context {
	ctx := &Context{inj.New(), g, newResponseWriter(w)}

	ctx.Register(r)
	ctx.Register(ctx.rw)
    // register our ResponseWriter as an http.ResponseWriter as well for
    // net/http HandlerFunc compatibility
	ctx.RegisterAs(ctx.rw, (*http.ResponseWriter)(nil))
	// register this instance with itself
	ctx.Register(*ctx)
	return ctx
}

// handle executes all of the Handler instances. Returned values from Handlers are
// ignored but errors panic. If a handler begins writing a response further handlers
// are not executed.
func (ctx *Context) handle() {
	for _, h := range append(ctx.g.handlers, ctx.g.defaultHandler) {
		vals, err := ctx.Call(h, ctx.g.Injector)

		// If a Handler returns values, and if the first value is a glue.AfterHandler
		// defer it to allow post-request logic
		if len(vals) > 0 {
			if vals[0].Type() == reflect.TypeOf(AfterHandler(nil)) {
				defer vals[0].Call([]reflect.Value{reflect.ValueOf(*ctx)})
			} else if len(vals) == 1 {
				log.Printf("glue: middleware didn't return a %T. It is instead of type: %+v\n", AfterHandler(nil), vals[0].Type())
			} else {
				log.Printf("glue: middleware didn't return a %T. It instead returned %d values: %+v", AfterHandler(nil), len(vals), vals)
			}
		}
		if err != nil {
			panic(err)
		}
		if ctx.rw.WroteHeader() {
			break
		}
	}
}
