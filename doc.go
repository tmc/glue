// Package glue provides a simple interface to writing HTTP services in Go
//
// It aims to be small and as simple as possible while exposing a pleasant API.
//
// Glue uses reflection and dependency injection to provide a flexible API for your
// HTTP endpoints. There is an obvious tradeoff here. The cost of this flexibility
// is some static safety and some performance overhead (though this appears
// negligible in benchmarking).
//
// godoc: http://godoc.org/github.com/tmc/glue
//
// Features:
//
//  * small (~250LOC)
//  * compatible with the net/http Handler and HandleFunc interfaces.
//  * provides mechanism for before and after request middleware
//
// Basic Example:
//  package main
//  import "github.com/tmc/glue"
//
//  func main() {
//      g := glue.New()
//      g.Get("/", func() string {
//          return "hello world"
//      })
//      g.Listen() // listens on :5000 by default (uses PORT environtment variable)
//  }
//
// Example showing middleware, logging, routing, and static file serving:
//  g := glue.New()
//  // Register a new type with the underlying DI container
//  g.Register(log.New(os.Stderr, "[glue example] ", log.LstdFlags))
//  // Add a new glue.Handler that will be invoked for each request
//  g.AddHandler(loggers.NewApacheLogger())
//  // Add a handler using routing and parameter capture
//  g.Get("/{type}_teapot", func(r *http.Request) (int, string) {
//      return http.StatusTeapot, "that is " + r.URL.Query().Get(":type") + "!"
//  })
//  // Serve static files
//  g.Get("/", http.FileServer(http.Dir("./static/")))
//  go g.Listen() // listens on 5000 by default (uses PORT environtment variable)
//
//  resp, err := http.Get("http://127.0.0.1:5000/purple_teapot")
//  if err != nil {
//      panic(err)
//  }
//  defer resp.Body.Close()
//  body, err := ioutil.ReadAll(resp.Body)
//  fmt.Println(resp.Status, string(body))
//  // Output:
//  // 418 I'm a teapot that is purple!
//
// glue is influenced by martini and basically co-opts gorilla's pat muxing for routing.
package glue
