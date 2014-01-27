// largely cribbed from gorilla/pat

package glue

import (
	"net/http"
	"net/url"
	"path"
	"strings"
	"github.com/gorilla/mux"
)

func NewRouter() *Router {
	return &Router{}
}

// Router is a request router that implements a pat-like API.
//
// pat docs: http://godoc.org/github.com/bmizerany/pat
type Router struct {
	mux.Router
	NotFoundHandler Handler
}

// Add registers a pattern with a handler for the given request method.
func (r *Router) Add(meth, pat string, h Handler) *mux.Route {
	return r.NewRoute().PathPrefix(pat).Handler(routeHandler{Handler: h}).Methods(meth)
}

// Delete registers a pattern with a handler for DELETE requests.
func (r *Router) Delete(pat string, h Handler) *mux.Route {
	return r.Add("DELETE", pat, h)
}

// Get registers a pattern with a handler for GET requests.
func (r *Router) Get(pat string, h Handler) *mux.Route {
	return r.Add("GET", pat, h)
}

// Post registers a pattern with a handler for POST requests.
func (r *Router) Post(pat string, h Handler) *mux.Route {
	return r.Add("POST", pat, h)
}

// Put registers a pattern with a handler for PUT requests.
func (r *Router) Put(pat string, h Handler) *mux.Route {
	return r.Add("PUT", pat, h)
}

func (r *Router) Handle(w http.ResponseWriter, req *http.Request, c Context) {
	// Clean path to canonical form and redirect.
	if p := cleanPath(req.URL.Path); p != req.URL.Path {
		w.Header().Set("Location", p)
		w.WriteHeader(http.StatusMovedPermanently)
		return
	}
	var match mux.RouteMatch
	var handler http.Handler
	if matched := r.Match(req, &match); matched {
		if rhandler, ok := match.Handler.(routeHandler); ok {
			rhandler.ctx = c
			handler = rhandler
		}
		// @todo coersion could fail
		registerVars(req, match.Vars)
	}

	if handler == nil {
		if r.NotFoundHandler == nil {
			handler = http.NotFoundHandler()
		} else {
			handler = routeHandler{r.NotFoundHandler, c}
		}
	}
	handler.ServeHTTP(w, req)
}

// registerVars adds the matched route variables to the URL query.
func registerVars(r *http.Request, vars map[string]string) {
	parts, i := make([]string, len(vars)), 0
	for key, value := range vars {
		parts[i] = url.QueryEscape(":"+key) + "=" + url.QueryEscape(value)
		i++
	}
	q := strings.Join(parts, "&")
	if r.URL.RawQuery == "" {
		r.URL.RawQuery = q
	} else {
		r.URL.RawQuery += "&" + q
	}
}

// cleanPath returns the canonical path for p, eliminating . and .. elements.
// Borrowed from the net/http package.
func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)
	// path.Clean removes trailing slash except for root;
	// put the trailing slash back if necessary.
	if p[len(p)-1] == '/' && np != "/" {
		np += "/"
	}
	return np
}

type routeHandler struct {
	Handler
	ctx Context
}

func (rh routeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vals, err := rh.ctx.Call(rh.Handler)
	if err != nil {
		panic(err)
	}
	// if the handler returned something, write it to the http response
	if len(vals) > 0 {
		_, err := rh.ctx.Call(func(rw http.ResponseWriter, rHandler ResponseHandler) {
			rHandler(rw, vals)
		}, rh.ctx.Glue.Injector)
		if err != nil {
			panic(err)
		}
	}
}
