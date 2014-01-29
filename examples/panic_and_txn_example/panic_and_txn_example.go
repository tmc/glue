// Small usage example with panic recovery and database transaction management
package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/lib/pq"
	"github.com/tmc/glue"
)

func main() {
	g := glue.New()
	url := os.Getenv("DATABASE_URL")
	connection, _ := pq.ParseURL(url)
	db, _ := sql.Open("postgres", connection)

	g.AddHandler(DBTxnManagement())
	g.AddHandler(PanicRecover)

	g.Register(db)
	g.Get("/div/{n}", func(r *http.Request) string {
		n, _ := strconv.Atoi(r.URL.Query().Get(":n"))
		return fmt.Sprint(1.0 / n)
	})

	go g.Listen()
	doreq := func(n int) {
		url := fmt.Sprintf("http://127.0.0.1:5000/div/%d", n)
		resp, _ := http.Get(url)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("RESPONSE:", resp.Status, string(body))
	}
	doreq(1)
	doreq(0)
}

func PanicRecover() glue.AfterHandler {
	return func(c glue.Context) {
		if err := recover(); err != nil {
			log.Println("encountered panic:", err)

			c.Call(func(rw http.ResponseWriter) {
				rw.WriteHeader(500)
				fmt.Fprintln(rw, err)
			})
		}
	}
}
func DBTxnManagement() glue.Handler {
	return func(rw *glue.ResponseWriter, r *http.Request, db *sql.DB) glue.AfterHandler {
		log.Println("(pretend) BEGIN TXN")

		return func(ctx glue.Context) {
			if rw.Status == 500 {
				log.Println("(pretend) ROLLBACK!")
			} else {
				log.Println("(pretend) COMMIT")
			}
		}
	}
}
