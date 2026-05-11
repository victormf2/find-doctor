package httpx

import (
	"encoding/json"
	"net/http"

	"github.com/victormf2/gox/internal/generators"
	"github.com/victormf2/gox/problem"
	"github.com/victormf2/gox/types"
)

func Http(pattern string, operationConstructor any) types.IEndpoint {

	// validates the provided pattern
	http.NewServeMux().HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {})

	httpEndpoint := &HttpEndpoint{
		pattern:             pattern,
		operationDescriptor: generators.NewOperationDescriptor(operationConstructor),
	}

	return httpEndpoint
}

type HttpEndpoint struct {
	operationDescriptor *generators.OperationDescriptor
	pattern             string
}

// Descriptor implements [IEndpoint].
func (h *HttpEndpoint) Descriptor() any {
	return h
}

var _ types.IEndpoint = &HttpEndpoint{}

func BindJSON(r *http.Request, value any) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(value)
	if err != nil {
		return problem.BadRequest("invalid payload",
			problem.WithDetail("expected payload in JSON format"),
		)
	}

	return nil
}

func WriteJSON(w http.ResponseWriter, value any) error {
	encoder := json.NewEncoder(w)
	err := encoder.Encode(value)
	if err != nil {
		return problem.InternalServerError("output error",
			problem.WithDetail("the request was successfully processed, but we failed to deliver the response"),
		)
	}

	return nil
}
