package types

type IApplicationBuilder interface {
	AddEndpoint(endpoint IEndpoint)
	Build() IApplication
}

type IEndpoint interface {
	Descriptor() any
}

type IApplication interface {
	Endpoints() []IEndpoint
}
