package linq

import (
	"fmt"
	"reflect"
	"strings"
)

// genericType represents a any reflect.Type.
type genericType int

var genericTp = reflect.TypeOf(new(genericType)).Elem()

// functionCache keeps genericFunc reflection objects in cache.
type functionCache struct {
	MethodName string
	ParamName  string
	FnValue    reflect.Value
	FnType     reflect.Type
	TypesIn    []reflect.Type
	TypesOut   []reflect.Type
}

// genericFunc is a type used to validate and call dynamic functions.
type genericFunc struct {
	Cache *functionCache
}

// Call calls a dynamic function.
func (g *genericFunc) Call(params ...interface{}) interface{} {
	paramsIn := make([]reflect.Value, len(params))
	for i, param := range params {
		paramsIn[i] = reflect.ValueOf(param)
	}
	paramsOut := g.Cache.FnValue.Call(paramsIn)
	if len(paramsOut) >= 1 {
		return paramsOut[0].Interface()
	}
	return nil
}

// newGenericFunc instantiates a new genericFunc pointer
func newGenericFunc(methodName, paramName string, fn interface{}, validateFunc func(*functionCache) error) (*genericFunc, error) {
	cache := &functionCache{}
	cache.FnValue = reflect.ValueOf(fn)

	if cache.FnValue.Kind() != reflect.Func {
		return nil, fmt.Errorf("%s: parameter [%s] is not a function type. It is a '%s'", methodName, paramName, cache.FnValue.Type())
	}
	cache.MethodName = methodName
	cache.ParamName = paramName
	cache.FnType = cache.FnValue.Type()
	numTypesIn := cache.FnType.NumIn()
	cache.TypesIn = make([]reflect.Type, numTypesIn)
	for i := 0; i < numTypesIn; i++ {
		cache.TypesIn[i] = cache.FnType.In(i)
	}

	numTypesOut := cache.FnType.NumOut()
	cache.TypesOut = make([]reflect.Type, numTypesOut)
	for i := 0; i < numTypesOut; i++ {
		cache.TypesOut[i] = cache.FnType.Out(i)
	}
	if err := validateFunc(cache); err != nil {
		return nil, err
	}

	return &genericFunc{Cache: cache}, nil
}

// simpleParamValidator creates a function to validate genericFunc based in the
// In and Out function parameters.
func simpleParamValidator(In []reflect.Type, Out []reflect.Type) func(cache *functionCache) error {
	return func(cache *functionCache) error {
		var isValid = func() bool {
			if In != nil {
				if len(In) != len(cache.TypesIn) {
					return false
				}
				for i, paramIn := range In {
					if paramIn != genericTp && paramIn != cache.TypesIn[i] {
						return false
					}
				}
			}
			if Out != nil {
				if len(Out) != len(cache.TypesOut) {
					return false
				}
				for i, paramOut := range Out {
					if paramOut != genericTp && paramOut != cache.TypesOut[i] {
						return false
					}
				}
			}
			return true
		}

		if !isValid() {
			return fmt.Errorf("%s: parameter [%s] has a invalid function signature. Expected: '%s', actual: '%s'", cache.MethodName, cache.ParamName, formatFnSignature(In, Out), formatFnSignature(cache.TypesIn, cache.TypesOut))
		}
		return nil
	}
}

// newElemTypeSlice creates a slice of items elem types.
func newElemTypeSlice(items ...interface{}) []reflect.Type {
	typeList := make([]reflect.Type, len(items))
	for i, item := range items {
		typeItem := reflect.TypeOf(item)
		if typeItem.Kind() == reflect.Ptr {
			typeList[i] = typeItem.Elem()
		}
	}
	return typeList
}

// formatFnSignature formats the func signature based in the parameters types.
func formatFnSignature(In []reflect.Type, Out []reflect.Type) string {
	paramInNames := make([]string, len(In))
	for i, typeIn := range In {
		if typeIn == genericTp {
			paramInNames[i] = "T"
		} else {
			paramInNames[i] = typeIn.String()
		}

	}
	paramOutNames := make([]string, len(Out))
	for i, typeOut := range Out {
		if typeOut == genericTp {
			paramOutNames[i] = "T"
		} else {
			paramOutNames[i] = typeOut.String()
		}
	}
	return fmt.Sprintf("func(%s)%s", strings.Join(paramInNames, ","), strings.Join(paramOutNames, ","))
}
