package apihttp

import (
	"net/http"

	"github.com/victormf2/gosyringe"
	"github.com/victormf2/gox/problem"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request, err error, c *gosyringe.Container) {
	// Edit this function to add your own error handling logic.
	//
	// You can use the container to resolve any dependency from the request scope, for example,
	// to log context information, or audit related stuff.
	//
	// You can use a custom type in your application that implements the problem package interfaces
	// problem.IProblemType, problem.IProblemTitle, etc to output the result according to the
	// Problem Details RFC 9457 https://www.rfc-editor.org/rfc/rfc9457.html.
	//
	// Example: logging the error
	// slog.With("route", r.Pattern).Error(err.Error())

	problem.WriteErrorJSON(w, err)
}
