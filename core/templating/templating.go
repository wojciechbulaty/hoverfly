package templating

import (
	"github.com/aymerick/raymond"
)

type TemplatingData struct {
	Request Request
}

type Request struct {
	QueryParam map[string][]string
	PathParam  []string
	Scheme     string
}

func ApplyTemplate() (string, error) {
	return raymond.Render(`
	Scheme: {{ Request.Scheme }}

	Query param value: {{ Request.QueryParam.Singular }}

	Query param value by index: {{ Request.QueryParam.Multiple.[0] }}
	Query param value by index: {{ Request.QueryParam.Multiple.[1] }}

	List of query param values: {{ Request.QueryParam.Multiple}}
	Looping through query params: {{#each Request.QueryParam.Multiple}}{{ this }} {{/each}}

	Path param value: {{ Request.PathParam.[0] }}
	Looping through path params: {{#each Request.QueryParam.Multiple}}{{ this }} {{/each}}

	`, TemplatingData{

		Request: Request{
			QueryParam: map[string][]string{
				"Singular": {"one"},
				"Multiple": {"one", "two"},
			},
			PathParam: []string{"foo", "bar"},
			Scheme:    "https",
		},
	})
}
