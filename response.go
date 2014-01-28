package glue

import (
	"net/http"
	"reflect"
)

// responseWriter is an augmented http.ResponseWriter that exposes some additional fields
type responseWriter struct {
	http.ResponseWriter
	Size   int // the number of bytes that have been written as a response body
	Status int // the status code that has been written to the response (or zero if unwritten)
}

// newResponseWriter creates a new responseWriter given an http.ResponseWriter
func newResponseWriter(rw http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: rw,
	}
}

// WroteHeader indicates if a header has been written (and a response has been started)
func (rw *responseWriter) WroteHeader() bool {
	return rw.Status != 0
}

// Write writes the data to the connection as part of an HTTP reply.
// If WriteHeader has not yet been called, Write calls WriteHeader(http.StatusOK)
// before writing the data.  If the Header does not contain a
// Content-Type line, Write adds a Content-Type set to the result of passing
// the initial 512 bytes of written data to DetectContentType.
func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.WroteHeader() {
		rw.WriteHeader(http.StatusOK)
	}
	size, err := rw.ResponseWriter.Write(b)
	rw.Size += size
	return size, err
}

// WriteHeader sends an HTTP response header with status code.
// If WriteHeader is not called explicitly, the first call to Write
// will trigger an implicit WriteHeader(http.StatusOK).
// Thus explicit calls to WriteHeader are mainly used to
// send error codes.
func (rw *responseWriter) WriteHeader(status int) {
	rw.ResponseWriter.WriteHeader(status)
	rw.Status = status
}

// ResponseHandler is a type that writes an HTTP response given a slice of reflect.Value
type ResponseHandler func(http.ResponseWriter, []reflect.Value)

// defaultResponseHandler returns the default ResponseHandler
//
// If given one value it attempts to write it either as a byte slice or a string.
//
// If given more than one value and the first is an int it uses the first value as the response code
// and the second value as the response body. Additional values are ignored.
func defaultResponseHandler() ResponseHandler {
	return func(res http.ResponseWriter, vals []reflect.Value) {
		var v reflect.Value
		if len(vals) > 1 && vals[0].Kind() == reflect.Int {
			res.WriteHeader(int(vals[0].Int()))
			v = vals[1]
		} else if len(vals) > 0 {
			v = vals[0]
		}
		if v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.Kind() == reflect.Slice && v.Type() == reflect.TypeOf([]byte(nil)) {
			res.Write(v.Bytes())
		} else {
			res.Write([]byte(v.String()))
		}
	}
}
