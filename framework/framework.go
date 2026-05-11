package framework

import (
	"github.com/victormf2/framework/internal"
	"github.com/victormf2/framework/internal/sharehack"
	"github.com/victormf2/framework/types"
)

type IApplication = types.IApplication
type IApplicationBuilder = types.IApplicationBuilder
type IEndpoint = types.IEndpoint

func New() IApplicationBuilder {
	return internal.NewApplication()
}

var GoxError = sharehack.GoxError
