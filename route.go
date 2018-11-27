package waffle

import (
	"path"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// RouteHandler represents an version group of API endpoints
type RouteHandler interface {
	Name() string
	Routes() []Route
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

// Handlers represents a group of route handlers
type Handlers map[string]RouteHandler

// Router constructs a Gorilla mux router from the handlers
func (h Handlers) Router(prefix string) (router *mux.Router) {
	router = mux.NewRouter().StrictSlash(true)

	for name, handler := range h {
		for _, route := range handler.Routes() {
			router.Methods(route.Method).
				PathPrefix(prefix).
				Path(path.Join("/", name, route.Path)).
				Name(route.Name).
				Handler(route.Handler)
		}

		router.Methods("GET").
			PathPrefix(prefix).
			Path(path.Join("/", name, "docs")).
			Name("docs").
			Handler(docsWrapper(handler))
	}
	router.Handle("/metrics", promhttp.Handler())

	return
}

// Exampler represents a resource that can give examples of itself
type Exampler interface {
	Example(...interface{}) interface{}
}
