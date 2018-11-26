package waffle

import (
	"encoding/json"
	"net/http"
	"net/url"
	"path"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Handlers represents a group of route handlers
type Handlers map[string]RouteHandler

// Router constructs a Gorilla mux router from the handlers
func (h Handlers) Router() (router *mux.Router) {
	router = mux.NewRouter().StrictSlash(true)
	router.Handle("/metrics", promhttp.Handler())
	for name, handler := range h {
		routes := handler.Routes()

		for _, route := range routes {
			router.Methods(route.Method).
				Path(path.Join("/", name, route.Path)).
				Name(route.Name).
				Handler(route.Handler)
		}

		router.Methods("GET").
			Path(path.Join("/", name, "docs")).
			Name("docs").
			Handler(docsWrapper(handler))
	}
	return
}

// Route represents an API route and its associated handler function
type Route struct {
	Name        string   `json:"name"`
	Method      string   `json:"method"`
	Path        string   `json:"path"`
	Description string   `json:"description"`
	Params      Exampler `json:"params"`
	Accepts     Exampler `json:"accepts"`
	Returns     Exampler `json:"returns"`
	Handler     Endpoint `json:"-"`
}

// RouteHandler represents an version group of API endpoints
type RouteHandler interface {
	Name() string
	Routes() []Route
}

// Exampler represents a resource that can give examples of itself
type Exampler interface {
	Example(...interface{}) interface{}
}

// Endpoint wraps a HTTP handler function with nicer arguments and return value
type Endpoint func(r *http.Request, query url.Values) Status

// ServeHTTP implements the necessary chaining functionality for HTTP middleware
func (f Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status := f(r, r.URL.Query())
	if status.Err != nil {
		status.ErrString = status.Err.Error()
	}
	w.WriteHeader(status.code)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(status)
	if err != nil {
		http.Error(w, "failed to encode response", 500)
	}
}

// Status is a custom status object returned by all endpoints.
type Status struct {
	code      int
	Err       error       `json:"-"`
	Result    interface{} `json:"result,omitempty"`
	ErrString string      `json:"error,omitempty"`
}

// Ok indicates a request was all good
func Ok(code int, result interface{}) Status {
	return Status{code: code, Result: result}
}

// Err is a helper for when things go wrong
func Err(code int, err error) Status {
	return Status{code: code, Err: err}
}
