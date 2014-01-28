// Package glue provides a simple interface to writing HTTP services in Go
//
//
// Example:
//  g := glue.New()
//  g.Register(log.New(os.Stdout, "[glue example]", log.LstdFlags))
//  g.Add(loggers.NewApacheLogger())
//  g.Get("/{type}_teapot", func(r *http.Request) (int, string) {
//      return http.StatusTeapot, "that is " + r.URL.Query().Get(":type") + "!"
//  })
//  g.Get("/", http.FileServer(http.Dir("./public/")))
//  go g.Listen()
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
