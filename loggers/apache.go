package loggers

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/tmc/glue"
)

// NewApacheLogger returns a glue.Handler that prints an Apache-style log
// message to the registered log.Logger within the glue.Context
// panics if no log.Logger is registered
func NewApacheLogger() glue.Handler {
	return func(w *glue.ResponseWriter, r *http.Request, logger *log.Logger) glue.AfterHandler {

		remoteAddr := r.RemoteAddr
		if index := strings.LastIndex(remoteAddr, ":"); index != -1 {
			remoteAddr = remoteAddr[:index]
		}
		referer := r.Referer()
		if "" == referer {
			referer = "-"
		}
		userAgent := r.UserAgent()
		if "" == userAgent {
			userAgent = "-"
		}
		start := time.Now()
		return func(c glue.Context) {

			logger.Printf(
				"%s %s %s [%v] \"%s %s %s\" %d %d \"%s\" \"%s\" %d\n",
				remoteAddr,
				"-", // remote logname, not supported
				"-", // @todo get username from request
				start.Format("02/Jan/2006:15:04:05 -0700"),
				r.Method,
				r.RequestURI,
				r.Proto,
				w.Status,
				w.Size,
				referer,
				userAgent,
				time.Now().Sub(start).Nanoseconds()/10e3,
			)
		}
	}
}
