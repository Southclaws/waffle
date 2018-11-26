package waffle

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/alecthomas/template"
	"github.com/dyninc/qstring"
)

type routeForTemplate struct {
	Route
	ParamsSerialised  string
	AcceptsSerialised string
	ReturnsSerialised string
}

const documentationHeader = `# Server API: %s

This is an automatically generated documentation page for the %s API endpoints.

`

const documentationRouteTemplate = `## {{ .Name }}

` + "`" + `{{ .Method }}` + "`" + `: ` + "`" + `/{{ .Version }}{{ .Path }}` + "`" + `

{{ .Description }}
{{ if .Params }}
### Query parameters

Example: ` + "`" + `{{ .ParamsSerialised }}` + "`" + `
{{ end }}{{ if .AcceptsSerialised}}
### Accepts

` + "```json" + `
{{ .AcceptsSerialised }}
` + "```" + `
{{ else }}{{ end }}{{ if .ReturnsSerialised}}
### Returns

` + "```json" + `
{{ .ReturnsSerialised }}
` + "```" + `
{{ else }}{{ end }}
`

func docsWrapper(handler RouteHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf(documentationHeader, handler.Name(), handler.Name())))
		for _, route := range handler.Routes() {
			docsForRoute(route, w)
		}
		w.Header().Set("Content-Type", "text/markdown")
		return
	}
}

func docsForRoute(route Route, w io.Writer) {
	var err error

	obj := routeForTemplate{
		Route: route,
	}

	if route.Params != nil {
		obj.ParamsSerialised, err = qstring.MarshalString(route.Params.Example())
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if route.Accepts != nil {
		acceptsSerialised, err2 := json.MarshalIndent(route.Accepts.Example(), "", "    ")
		if err2 != nil {
			fmt.Println(err2)
			return
		}
		obj.AcceptsSerialised = string(acceptsSerialised)
	}

	if route.Returns != nil {
		returnsSerialised, err2 := json.MarshalIndent(route.Returns.Example(), "", "    ")
		if err2 != nil {
			fmt.Println(err2)
			return
		}
		obj.ReturnsSerialised = string(returnsSerialised)
	}

	tpl, err := template.New("doc").Parse(documentationRouteTemplate)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = tpl.Execute(w, obj)
	if err != nil {
		fmt.Println(err)
		return
	}
}
