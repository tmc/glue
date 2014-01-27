// Package glue provides a simple interface to writing HTTP services in Go
//
package glue

import (
	"net/http"
	"os"

	"github.com/tmc/inj"
)

// Handler is a generic type that must be callable
type Handler interface{}

type Glue struct {
	inj.Injector
	*Router
	handlers       []Handler // the set of handlers invoked for every request
	defaultHandler Handler   // the handler that is handled last
}

func New() *Glue {
	r := &Router{}
	g := &Glue{inj.New(),
		r,
		[]Handler{},
		r.Handle,
	}
	// register some expected types
	g.Register(defaultResponseHandler())
	return g
}

func (g *Glue) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.newContext(w, r).handle()
}

// Add adds a handler to the default set of handlers for a Glue instance
func (g *Glue) Add(handler Handler) {
	//@todo verify is func
	g.handlers = append(g.handlers, handler)
}

// Listen attempts to ListenAndServe based on the environment variables HOST and PORT
func (g *Glue) Listen() {
	port, host := os.Getenv("PORT"), os.Getenv("HOST")
	if port == "" {
		port = "5000"
	}
	http.ListenAndServe(host+":"+port, g)
}
