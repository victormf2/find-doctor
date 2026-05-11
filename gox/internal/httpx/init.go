package httpx

import (
	_ "embed"

	"github.com/victormf2/gox/generators"
)

//go:embed templates/routes.tmpl
var routesTmpl string

type HttpTemplateData struct {
	Pattern   string
	Operation generators.OperationTemplateData
}

func init() {
	generators.RegisterEndpointGenerator(func(endpoint *HttpEndpoint) generators.GeneratorTemplate {
		templateData := generators.GeneratorTemplate{
			Tmpl: routesTmpl,
			EndpointTemplateData: HttpTemplateData{
				Pattern:   endpoint.pattern,
				Operation: endpoint.operationDescriptor.TemplateData(),
			},
		}
		return templateData
	})
}
