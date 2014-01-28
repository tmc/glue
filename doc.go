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
//  g.Register(log.New(os.Stdout, "[glue example] ", log.LstdFlags))
//  g.Add(loggers.NewApacheLogger())
//  g.Get("/{type}_teapot", func(r *http.Request) (int, string) {
//      return http.StatusTeapot, "that is " + r.URL.Query().Get(":type") + "!"
//  })
//  g.Get("/", http.FileServer(http.Dir("./public/")))
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
