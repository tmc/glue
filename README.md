# glue
    import "github.com/tmc/glue"

Package glue provides a simple interface to writing HTTP services in Go

It aims to be small and as simple as possible while exposing a pleasant API.

Glue uses reflection and dependency injection to provide a flexible API for your
HTTP endpoints. There is an obvious tradeoff here. The cost of this flexibility
is some static safety and some performance overhead (though this appears
negligible in benchmarking).

Contributions welcome!

godoc: http://godoc.org/github.com/tmc/glue

Features

	* small (~250LOC)
	* compatible with the net/http Handler and HandleFunc interfaces.
	* provides mechanism for before and after request middleware


Basic Example:


```go
	package main
	import "github.com/tmc/glue"
	
	func main() {
	    g := glue.New()
	    g.Get("/", func() string {
	        return "hello world"
	    })
	    g.Listen() // listens on :5000 by default (uses PORT environtment variable)
	}
```

Example showing middleware, logging, routing, and static file serving:

```go
	g := glue.New()
	g.Register(log.New(os.Stdout, "[glue example] ", log.LstdFlags))
	g.Add(loggers.NewApacheLogger())
	g.Get("/{type}_teapot", func(r *http.Request) (int, string) {
	    return http.StatusTeapot, "that is " + r.URL.Query().Get(":type") + "!"
	})
	g.Get("/", http.FileServer(http.Dir("./public/")))
	go g.Listen() // listens on 5000 by default (uses PORT environtment variable)
	
	resp, err := http.Get("http://127.0.0.1:5000/purple_teapot")
	if err != nil {
	    panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(resp.Status, string(body))
	// Output:
	// 418 I'm a teapot that is purple!
```


## type AfterHandler
```go
type AfterHandler func(Context)
```
AfterHandler is a type that a glue Handler can return and have it invoked
after the default handler. This allows middleware to execute logic after a
response has started. See github.com/tmc/glue/loggers for an example.



## type Context
``` go
type Context struct {
    inj.Injector
    // contains filtered or unexported fields
}
```
Context represents the execution context for a request in Glue
It is a DI (Dependency Injection) container and contains an augmented
ResponseWriter



## type Glue
``` go
type Glue struct {
    inj.Injector
    // contains filtered or unexported fields
}
```
Glue is the primary struct that exposes routing and Handler registration


### func New
``` go
func New() *Glue
```
New prepares a new Glue instance and registers the default ResponseHandler


### func (\*Glue) Add
``` go
func (g *Glue) Add(handler Handler)
```
Add adds a handler to the default set of handlers for a Glue instance


### func (Glue) Delete
``` go
func (r Glue) Delete(pat string, h Handler) *mux.Route
```
Delete registers a pattern with a handler for DELETE requests.


### func (Glue) Get
``` go
func (r Glue) Get(pat string, h Handler) *mux.Route
```
Get registers a pattern with a handler for GET requests.


### func (Glue) Handle
``` go
func (r Glue) Handle(w http.ResponseWriter, req *http.Request, c Context)
```
Handle is a glue.Handler that does route matching and invokes the registered
glue.Handler for a route.

If a route is not found the NotFoundHandler is invoked.


### func (\*Glue) Listen
``` go
func (g *Glue) Listen()
```
Listen attempts to ListenAndServe based on the environment variables HOST and PORT


### func (Glue) Post
``` go
func (r Glue) Post(pat string, h Handler) *mux.Route
```
Post registers a pattern with a handler for POST requests.


### func (Glue) Put
``` go
func (r Glue) Put(pat string, h Handler) *mux.Route
```
Put registers a pattern with a handler for PUT requests.


### func (\*Glue) ServeHTTP
``` go
func (g *Glue) ServeHTTP(w http.ResponseWriter, r *http.Request)
```
ServeHTTP satisfies the http.Handler interface


## type Handler
``` go
type Handler interface{}
```
Handler is a generic type that must be a callable function.

It is invoked with the Call method of inj.Injector (http://godoc.org/github.com/tmc/inj#Injector.Call) which provides DI
(Dependency Injection) based on the types of arguments it accepts.

Accepting a glue.Context allows you to inspect the DI container and examine
the currently registered types.

The default registered ResponseHandler expects Handlers to return either one or two values.

If one value, it should return a string or a byte slice.
If two values, the first should be an int which will be used as the return code.



## type ResponseHandler
``` go
type ResponseHandler func(http.ResponseWriter, []reflect.Value)
```
ResponseHandler is a type that writes an HTTP response given a slice of reflect.Value

## type ResponseWriter
``` go
type ResponseWriter struct {
    http.ResponseWriter
    Size   int // the number of bytes that have been written as a response body
    Status int // the status code that has been written to the response (or zero if unwritten)
}
```
ResponseWriter is an augmented http.ResponseWriter that exposes some additional fields


### func (\*ResponseWriter) Write
``` go
func (rw *ResponseWriter) Write(b []byte) (int, error)
```
Write writes the data to the connection as part of an HTTP reply.
If WriteHeader has not yet been called, Write calls WriteHeader(http.StatusOK)
before writing the data.  If the Header does not contain a
Content-Type line, Write adds a Content-Type set to the result of passing
the initial 512 bytes of written data to DetectContentType.



### func (\*ResponseWriter) WriteHeader
``` go
func (rw *ResponseWriter) WriteHeader(status int)
```
WriteHeader sends an HTTP response header with status code.
If WriteHeader is not called explicitly, the first call to Write
will trigger an implicit WriteHeader(http.StatusOK).
Thus explicit calls to WriteHeader are mainly used to
send error codes.



### func (\*ResponseWriter) WroteHeader
``` go
func (rw *ResponseWriter) WroteHeader() bool
```
WroteHeader indicates if a header has been written (and a response has been started)
