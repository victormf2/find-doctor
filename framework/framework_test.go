package framework_test

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/victormf2/framework"
	"github.com/victormf2/framework/problem"
)

func ExampleGoxError() {

	invalidProblem := `{"type":{"non":"sense"}}`
	err := json.Unmarshal([]byte(invalidProblem), &problem.Problem{})

	isGoxError := errors.Is(err, framework.GoxError)
	isInvalidUnmarshalError := errors.As(err, new(&json.UnmarshalTypeError{}))

	fmt.Println(isGoxError)
	fmt.Println(isInvalidUnmarshalError)
	// Output:
	// true
	// true

}
