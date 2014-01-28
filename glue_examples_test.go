package glue_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/lib/pq"
	"github.com/tmc/glue"
	"github.com/tmc/glue/loggers"
)

// Exapmle showing the use of Listen, routing, logging and static file serving
func ExampleGlue_Listen() {
	g := glue.New()
	g.Register(log.New(os.Stderr, "[glue example] ", log.LstdFlags))
	g.Add(loggers.NewApacheLogger())
	g.Get("/{type}_teapot", func(r *http.Request) (int, string) {
		return http.StatusTeapot, "that is " + r.URL.Query().Get(":type") + "!"
	})
	g.Get("/", http.FileServer(http.Dir("./static/")))
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

// Example showing hello world in Glue
func ExampleNew() {
	g := glue.New()

	g.Get("/", func() string {
		return "hello world"
	})
	go http.ListenAndServe(":5001", g)
	resp, _ := http.Get("http://127.0.0.1:5001/")
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	// Output:
	// hello world
}

// Example showing returning a status code and a byte slice
func ExampleNew_multiReturnAndByteSlice() {
	g := glue.New()

	g.Get("/teapot", func(r *http.Request) (int, []byte) {
		teapot, _ := json.Marshal(struct {
			Teapot struct {
				IsReady bool
			}
		}{})
		return http.StatusTeapot, teapot
	})
	go http.ListenAndServe(":5002", g)
	resp, _ := http.Get("http://127.0.0.1:5002/teapot")
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	// Output:
	// {"Teapot":{"IsReady":false}}
}

// Example showing the creation of a Glue instance and the registration of a sql.DB
func ExampleNew_registerExampleWithDB() {
	g := glue.New()
	url := os.Getenv("DATABASE_URL")
	connection, _ := pq.ParseURL(url)
	db, _ := sql.Open("postgres", connection)

	g.Register(db)
	g.Get("/", func(db *sql.DB) string {
		return fmt.Sprintf("db: %T\n", db)
	})
	go http.ListenAndServe(":5003", g)
	resp, _ := http.Get("http://127.0.0.1:5003/")
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(resp.Status, string(body))
	// Output:
	// 200 OK db: *sql.DB
}
