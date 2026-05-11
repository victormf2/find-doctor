package generators

import (
	"context"
	"path"
	"reflect"
	"runtime"
	"strings"
)

type GeneratorTemplate struct {
	Tmpl                 string
	EndpointTemplateData any
}

var (
	endpointGenerators = map[reflect.Type]reflect.Value{}
)

func RegisterEndpointGenerator(endpointFunc any) {
	reflectEndpointFunc := reflect.ValueOf(endpointFunc)
	if reflectEndpointFunc.Kind() != reflect.Func {
		panic("endpointFunc must be a function")
	}

	endpointFuncType := reflectEndpointFunc.Type()
	if endpointFuncType.NumIn() != 1 {
		panic("endpointFunc must receive a single argument")
	}

	if endpointFuncType.NumOut() != 1 {
		panic("endpointFunc must return GeneratorTemplate")
	}

	endpointFuncOutType := endpointFuncType.Out(0)
	if endpointFuncOutType != reflect.TypeFor[GeneratorTemplate]() {
		panic("endpointFunc must return GeneratorTemplate")
	}

	endpointDescriptorType := endpointFuncType.In(0)

	endpointGenerators[endpointDescriptorType] = reflectEndpointFunc
}

type OperationDescriptor struct {
	Constructor   any
	OperationType reflect.Type
	Input         *OperationInput
	Output        *OperationOutput
}

type OperationTemplateData struct {
	// Package is a string representation of the operation package import
	Package string
	// Type is a string representation of the operation type reference
	Type string
	// Constructor is a string representation of the operation constructor reference
	Constructor string
	// Input is a string representation of the operation input type reference
	Input string
	// Output is a string representation of the operation output type reference
	Output string
}

func (d *OperationDescriptor) TemplateData() OperationTemplateData {
	// Extract package path and name from the operation type
	// Dereference pointer types to get the underlying type
	opType := d.OperationType
	if opType.Kind() == reflect.Pointer {
		opType = opType.Elem()
	}

	pkgPath := opType.PkgPath()
	pkgName := path.Base(pkgPath)

	data := OperationTemplateData{
		Package:     pkgPath,
		Type:        pkgName + "." + opType.Name(),
		Constructor: constructorName(d.Constructor, pkgName),
	}

	if d.Input != nil && d.Input.InputType != nil {
		inputType := *d.Input.InputType
		if inputType.Kind() == reflect.Pointer {
			inputType = inputType.Elem()
		}
		data.Input = pkgName + "." + inputType.Name()
	}

	if d.Output != nil && d.Output.OutputType != nil {
		outputType := *d.Output.OutputType
		if outputType.Kind() == reflect.Pointer {
			outputType = outputType.Elem()
		}
		data.Output = pkgName + "." + outputType.Name()
	}

	return data
}

func constructorName(constructor any, pkgName string) string {
	// runtime.FuncForPC gives us the full symbol name, e.g.:
	// "github.com/victormf2/myapi/send_user_message.NewOperation"
	// We replace the full package path prefix with the short package name.
	fullName := runtime.FuncForPC(reflect.ValueOf(constructor).Pointer()).Name()
	// fullName looks like: "github.com/.../send_user_message.NewOperation"
	// Find the last slash to isolate "send_user_message.NewOperation"
	if idx := strings.LastIndex(fullName, "/"); idx >= 0 {
		return fullName[idx+1:] // "send_user_message.NewOperation"
	}
	return fullName
}

type OperationInput struct {
	HasContext bool
	InputType  *reflect.Type
}

type OperationOutput struct {
	HasError   bool
	OutputType *reflect.Type
}

func NewOperationDescriptor(constructor any) *OperationDescriptor {

	operationDescriptor := &OperationDescriptor{
		Constructor: constructor,
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
	operationDescriptor.OperationType = operationType

	doMethod, found := operationType.MethodByName("Do")
	if !found {
		panic("operationConstructor must return a type with a method Do")
	}

	if doMethod.Func.IsNil() {
		// interface method
		operationDescriptor.Input = NewInterfaceMethodInput(doMethod)
	} else {
		// concrete type method
		operationDescriptor.Input = NewConcreteTypeMethodInput(doMethod)
	}
	operationDescriptor.Output = NewMethodOutput(doMethod)

	return operationDescriptor
}

func NewInterfaceMethodInput(method reflect.Method) *OperationInput {
	if !method.Func.IsNil() {
		panic(method.Name + " is not an interface method")
	}

	// interface methods have no receiver
	return NewMethodInput(method, 0)
}

func NewConcreteTypeMethodInput(method reflect.Method) *OperationInput {
	if method.Func.IsNil() {
		panic(method.Name + " is not a concrete type method")
	}

	// concrete type methods have receiver
	return NewMethodInput(method, 1)
}

func NewMethodInput(method reflect.Method, offset int) *OperationInput {
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
			operationInput.HasContext = true
		}

		operationInput.InputType = new(method.Type.In(inputIndex))
	}

	return operationInput

}

func NewMethodOutput(method reflect.Method) *OperationOutput {
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
			operationOutput.HasError = true
		}

		operationOutput.OutputType = new(method.Type.Out(0))
	}

	return operationOutput
}
