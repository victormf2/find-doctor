package generators

import "reflect"

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

	if endpointFuncType.NumOut() != 0 {
		panic("endpointFunc must return no value")
	}

	endpointDescriptorType := endpointFuncType.In(0)

	endpointGenerators[endpointDescriptorType] = reflectEndpointFunc
}
