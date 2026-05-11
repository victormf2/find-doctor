package gox_test

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/victormf2/gox"
	"github.com/victormf2/gox/problem"
)

func ExampleGoxError() {

	invalidProblem := `{"type":{"non":"sense"}}`
	err := json.Unmarshal([]byte(invalidProblem), &problem.Problem{})

	isGoxError := errors.Is(err, gox.GoxError)
	isInvalidUnmarshalError := errors.As(err, new(&json.UnmarshalTypeError{}))

	fmt.Println(isGoxError)
	fmt.Println(isInvalidUnmarshalError)
	// Output:
	// true
	// true

}
