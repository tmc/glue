// cribbed from martini

package glue

import (
	"net/http"
	"reflect"
)

type ResponseWriter struct {
	http.ResponseWriter
	Size   int
	Status int
}

func NewResponseWriter(rw http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: rw,
	}
}

func (rw *ResponseWriter) WroteHeader() bool {
	return rw.Status != 0
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	if !rw.WroteHeader() {
		rw.WriteHeader(http.StatusOK)
	}
	size, err := rw.ResponseWriter.Write(b)
	rw.Size += size
	return size, err
}

func (rw *ResponseWriter) WriteHeader(status int) {
	rw.ResponseWriter.WriteHeader(status)
	rw.Status = status
}

type ResponseHandler func(http.ResponseWriter, []reflect.Value)

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
