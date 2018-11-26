package waffle

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// Endpoint wraps a HTTP handler function with nicer arguments and return value
type Endpoint func(r *http.Request, query url.Values) Status

// ServeHTTP implements the necessary chaining functionality for HTTP middleware
func (f Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status := f(r, r.URL.Query())
	if status.Err != nil {
		status.ErrString = status.Err.Error()
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status.code)
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
