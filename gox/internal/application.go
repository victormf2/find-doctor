package internal

import "slices"

type Application struct {
	endpoints []IEndpoint
}

// AddEndpoint implements [IApplicationBuilder].
func (a *Application) AddEndpoint(endpoint IEndpoint) {
	a.endpoints = append(a.endpoints, endpoint)
}

// Build implements [IApplicationBuilder].
func (a *Application) Build() IApplication {
	return a
}

// Endpoints implements [IApplication].
func (a *Application) Endpoints() []IEndpoint {
	return slices.Clone(a.endpoints)
}

var _ IApplication = &Application{}
var _ IApplicationBuilder = &Application{}

func NewApplication() *Application {
	return &Application{}
}
