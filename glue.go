package glue

import (
	"net/http"
	"os"

	"github.com/tmc/inj"
)

// Handler is a generic type that must be a callable function.
//
// It is invoked with the Call method of inj.Injector
// (http://godoc.org/github.com/tmc/inj#Injector.Call) which provides DI
// (Dependency Injection) based on the types of arguments it accepts.
//
// Accepting a glue.Context allows you to inspect the DI container and examine
// the currently registered types.
//
// The default registered ResponseHandler expects Handlers to return either one or two values.
//
// If one value, it should return a string or a byte slice.
// If two values, the first should be an int which will be used as the return code.
type Handler interface{}

// AfterHandler is a type that a glue Handler can return and have it invoked
// after the default handler. This allows middleware to execute logic after a
// response has started. See github.com/tmc/glue/loggers for an example.
type AfterHandler func(Context)

// Glue is the primary struct that exposes routing and Handler registration
type Glue struct {
	inj.Injector
	*router
	handlers       []Handler // the set of handlers invoked for every request
	defaultHandler Handler   // the handler that is handled last
}

// New prepares a new Glue instance and registers the default ResponseHandler
func New() *Glue {
	r := &router{}
	g := &Glue{inj.New(), r, []Handler{}, r.Handle}
	// register the default ResponseHandler
	g.Register(defaultResponseHandler())
	return g
}

// ServeHTTP satisfies the http.Handler interface
func (g *Glue) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.newContext(w, r).handle()
}

// Add adds a handler to the default set of handlers for a Glue instance
func (g *Glue) AddHandler(handler Handler) {
	//@todo verify is func
	g.handlers = append(g.handlers, handler)
}

// Listen attempts to ListenAndServe based on the environment variables HOST and PORT
func (g *Glue) Listen() error {
	port, host := os.Getenv("PORT"), os.Getenv("HOST")
	if port == "" {
		port = "5000"
	}
	return http.ListenAndServe(host+":"+port, g)
}
