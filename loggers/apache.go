package loggers

import (
	"log"
	"net/http"
	"github.com/tmc/glue"
)

func NewApacheLogger() glue.Handler {
	return func(w http.ResponseWriter, r *http.Request, c glue.Context) {
		log.Println("TODO")
	}
}
