package httpx

import (
	"context"
	"net/http"
	"reflect"

	"github.com/victormf2/framework/types"
)

func Http(pattern string, operationConstructor any) types.IEndpoint {

	// validates the provided pattern
	http.NewServeMux().HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {})

	httpEndpoint := &HttpEndpoint{
		pattern:             pattern,
		operationDescriptor: GetOperationDescriptor(operationConstructor),
	}

	return httpEndpoint
}

func GetOperationDescriptor(constructor any) *OperationDescriptor {

	operationDescriptor := &OperationDescriptor{
		constructor: constructor,
	}

	reflectOperationConstructor := reflect.ValueOf(constructor)
	if reflectOperationConstructor.Kind() != reflect.Func {
		panic("operationConstructor must be a function")
	}

	constructorType := reflectOperationConstructor.Type()
	if constructorType.NumOut() <= 0 || constructorType.NumOut() > 2 {
		panic("operationConstructor must either return a single value or a value and an error")
	}

	operationType := constructorType.Out(0)
	operationDescriptor.operationType = operationType

	doMethod, found := operationType.MethodByName("Do")
	if !found {
		panic("operationConstructor must return a type with a method Do")
	}

	if doMethod.Func.IsNil() {
		// interface method
		operationDescriptor.input = GetInterfaceMethodInput(doMethod)
	} else {
		// concrete type method
		operationDescriptor.input = GetConcreteTypeMethodInput(doMethod)
	}
	operationDescriptor.output = GetMethodOutput(doMethod)

	return operationDescriptor
}

func GetInterfaceMethodInput(method reflect.Method) *OperationInput {
	if !method.Func.IsNil() {
		panic(method.Name + " is not an interface method")
	}

	// interface methods have no receiver
	return GetMethodInput(method, 0)
}

func GetConcreteTypeMethodInput(method reflect.Method) *OperationInput {
	if method.Func.IsNil() {
		panic(method.Name + " is not a concrete type method")
	}

	// concrete type methods have receiver
	return GetMethodInput(method, 1)
}

func GetMethodInput(method reflect.Method, offset int) *OperationInput {
	operationInput := &OperationInput{}

	if method.Type.NumIn() > (2 + offset) {
		panic(method.Name + " method must either receive 0 arguments, a single argument, or a context.Context and other argument")
	}

	if method.Type.NumIn() > (0 + offset) {
		inputIndex := (0 + offset)
		if method.Type.NumIn() == (2 + offset) {
			// has context
			if method.Type.In(0+offset) != reflect.TypeFor[context.Context]() {
				panic(method.Name + " method must either receive 0 arguments, a single argument, or a context.Context and other argument")
			}
			inputIndex = (1 + offset)
			operationInput.hasContext = true
		}

		operationInput.inputType = new(method.Type.In(inputIndex))
	}

	return operationInput

}

func GetMethodOutput(method reflect.Method) *OperationOutput {
	operationOutput := &OperationOutput{}

	if method.Type.NumOut() > 2 {
		panic(method.Name + " method must either return no value, a single value, or a value and an error")
	}

	if method.Type.NumOut() > 0 {
		if method.Type.NumOut() == 2 {
			// has error
			if method.Type.Out(1) != reflect.TypeFor[error]() {
				panic(method.Name + " method must either return no value, a single value, or a value and an error")
			}
			operationOutput.hasError = true
		}

		operationOutput.outputType = new(method.Type.Out(0))
	}

	return operationOutput
}

type HttpEndpoint struct {
	operationDescriptor *OperationDescriptor
	pattern             string
}

type OperationDescriptor struct {
	constructor   any
	operationType reflect.Type
	input         *OperationInput
	output        *OperationOutput
}

type OperationInput struct {
	hasContext bool
	inputType  *reflect.Type
}

type OperationOutput struct {
	hasError   bool
	outputType *reflect.Type
}

// Descriptor implements [IEndpoint].
func (h *HttpEndpoint) Descriptor() any {
	return h
}

var _ types.IEndpoint = &HttpEndpoint{}
