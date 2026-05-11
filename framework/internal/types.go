package internal

import (
	"fmt"

	"github.com/victormf2/framework/internal/sharehack"
	"github.com/victormf2/framework/types"
)

type IApplicationBuilder = types.IApplicationBuilder
type IEndpoint = types.IEndpoint
type IApplication = types.IApplication

func WithGoxError(err error) error {
	return fmt.Errorf("%w: %w", sharehack.GoxError, err)
}
