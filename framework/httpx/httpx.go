package httpx

import (
	"encoding/json"
	"net/http"

	"github.com/victormf2/framework"
	"github.com/victormf2/framework/internal/httpx"
	"github.com/victormf2/framework/problem"
)

func Endpoint(pattern string, operationConstructor any) framework.IEndpoint {
	return httpx.Http(pattern, operationConstructor)
}

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
