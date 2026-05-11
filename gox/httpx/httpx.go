package httpx

import (
	"net/http"

	"github.com/victormf2/gox"
	"github.com/victormf2/gox/internal/httpx"
)

func Endpoint(pattern string, operationConstructor any) gox.IEndpoint {
	return httpx.Http(pattern, operationConstructor)
}

func BindJSON(r *http.Request, value any) error {
	return httpx.BindJSON(r, value)
}

func WriteJSON(w http.ResponseWriter, value any) error {
	return httpx.WriteJSON(w, value)
}
