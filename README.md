# Waffle

Quick and dirty micro-framework for API servers.

Yes, there are thousands of these, but this one is mine.

## Features

- Simple, opinionated endpoint declaration
- Designed for multi-level routes with versions/categories/etc
- Generates documentation with examples!

Example:

```go
func main() {
    server := &http.Server{
        Addr: "0.0.0.0:80",
        Handler: handlers.CORS(
            handlers.AllowedHeaders([]string{"X-Requested-With"}),
            handlers.AllowedOrigins([]string{"*"}),
            handlers.AllowedMethods([]string{"HEAD", "GET", "POST", "PUT", "OPTIONS"}),
        )(waffle.Handlers{
            "v1": V1.Init(db),
        }.Router()),
    }
    server.ListenAndServe()
}

// v1 package

// V1 satisfies waffle.RouteHandler interface
type V1 struct {
    DB sql.DB
}
func Init(db sql.DB) *V1 { return &V1{DB: db} }

func (p *V1) Name() string { return "API Version 1.x" }
func (p *V1) Routes() []waffle.Route {
    return []waffle.Route{
        {
            Name:        "things",
            Method:      http.MethodGet,
            Path:        "/{thingid}",
            Description: "Returns things",
            Handler:     p.getThings,
        },
    }
}

type Thing struct {
    Name string `json:"name"`
}

// Satisfies the waffle.Exampler interface for generating docs
func (t Thing) Example(args ...interface{}) (result interface{}) {
    return Thing{Name: "Example name"}
}
```
