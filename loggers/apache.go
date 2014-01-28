package loggers

import (
	"log"
	"net/http"
	"time"

	"github.com/tmc/glue"
)

// NewApacheLogger returns a glue.Handler that prints an Apache-style log
// message to the registered log.Logger within the glue.Context
// panics if no log.Logger is registered
func NewApacheLogger() glue.Handler {
	return func(w http.ResponseWriter, r *http.Request, logger *log.Logger) glue.AfterHandler {
		start := time.Now()
		return func(c glue.Context) {
			log.Println(r.Method, r.RequestURI, "took", time.Now().Sub(start))
		}
	}
}
