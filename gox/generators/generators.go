package generators

import (
	"github.com/victormf2/gox/internal/generators"
)

type GeneratorTemplate = generators.GeneratorTemplate
type OperationDescriptor = generators.OperationDescriptor
type OperationTemplateData = generators.OperationTemplateData

func NewOperationDescriptor(constructor any) *OperationDescriptor {
	return generators.NewOperationDescriptor(constructor)
}

func RegisterEndpointGenerator(endpointFunc any) {
	generators.RegisterEndpointGenerator(endpointFunc)
}
