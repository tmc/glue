package glue_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/tmc/glue"
	"github.com/tmc/glue/loggers"
)

func ExampleGlue_Listen() {
	g := glue.New()
	g.Register(log.New(os.Stderr, "[glue example] ", log.LstdFlags))
	g.Add(loggers.NewApacheLogger())
	g.Get("/{type}_teapot", func(r *http.Request) (int, string) {
		return http.StatusTeapot, "that is " + r.URL.Query().Get(":type") + "!"
	})
	g.Get("/", http.FileServer(http.Dir("./public/")))
	go g.Listen()

	resp, err := http.Get("http://127.0.0.1:5000/purple_teapot")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(resp.Status, string(body))
	// Output:
	// 418 I'm a teapot that is purple!
}
