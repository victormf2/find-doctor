package main

import (
	apihttp "find-doctor/api/http"
	"net/http"

	"github.com/victormf2/gosyringe"
)

func main() {
	c := gosyringe.NewContainer()
	apihttp.AddRoutes(c, http.DefaultServeMux)
	http.ListenAndServe(":8080", nil)
}
