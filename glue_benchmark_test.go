package glue_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/tmc/glue"
)

func init() {
	g := glue.New()
	g.Get("/a", func() string {
		return "a"
	})

	http.HandleFunc("/b", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprint(rw, "b")
	})

	http.Handle("/a", g)
	go http.ListenAndServe(":6001", nil)
}

func BenchmarkGlueIO(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Sleep(200 * time.Millisecond)
		resp, _ := http.Get("http://127.0.0.1:6001/a")
		ioutil.ReadAll(resp.Body)
	}
}

func BenchmarkHttpIO(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Sleep(200 * time.Millisecond)
		resp, _ := http.Get("http://127.0.0.1:6001/b")
		ioutil.ReadAll(resp.Body)
	}
}

func BenchmarkHttp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Sleep(200 * time.Millisecond)
		http.Get("http://127.0.0.1:6001/b")
	}
}

func BenchmarkGlue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Sleep(200 * time.Millisecond)
		http.Get("http://127.0.0.1:6001/a")
	}
}
