package gox

import (
	"github.com/victormf2/gox/internal"
	"github.com/victormf2/gox/internal/errorsx"
	"github.com/victormf2/gox/types"
)

type IApplication = types.IApplication
type IApplicationBuilder = types.IApplicationBuilder
type IEndpoint = types.IEndpoint

func New() IApplicationBuilder {
	return internal.NewApplication()
}

var GoxError = errorsx.GoxError
