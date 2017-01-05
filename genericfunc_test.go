package linq

import (
	"errors"
	"reflect"
	"testing"
)

func TestNewGenericFunc(t *testing.T) {
	tests := []struct {
		methodName     string
		paramName      string
		function       interface{}
		validationFunc func(*functionCache) error
		exception      error
	}{
		{ // A valid function
			"TestNewGenericFunc", "test1",
			func(item int) bool { return item > 10 },
			simpleParamValidator(newElemTypeSlice(new(int)), newElemTypeSlice(new(bool))),
			nil,
		},
		{ // A valid generic function
			"TestNewGenericFunc", "test1",
			func(item int) bool { return item > 10 },
			simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(bool))),
			nil,
		},
		{ //returns error when the function parameter has not the function kind
			"TestNewGenericFunc", "test2",
			"Not a function",
			simpleParamValidator(nil, []reflect.Type{}),
			errors.New("TestNewGenericFunc: parameter [test2] is not a function type. It is a 'string'"),
		},
		{ // Returns error when expected parameters number are not equal
			"TestNewGenericFunc", "test3",
			func(idx, item int) {},
			simpleParamValidator(newElemTypeSlice(new(int)), []reflect.Type{}),
			errors.New("TestNewGenericFunc: parameter [test3] has a invalid function signature. Expected: 'func(int)', actual: 'func(int,int)'"),
		},
		{ // Returns error when expected parameters types are not equal
			"TestNewGenericFunc", "test4",
			func(items ...int) bool { return false },
			simpleParamValidator(newElemTypeSlice(new([]bool)), newElemTypeSlice(new(bool))),
			errors.New("TestNewGenericFunc: parameter [test4] has a invalid function signature. Expected: 'func([]bool)bool', actual: 'func([]int)bool'"),
		},
		{ // Returns error when expected returns number are not equal
			"TestNewGenericFunc", "test5",
			func(item int) bool { return item > 10 },
			simpleParamValidator(newElemTypeSlice(new(int)), []reflect.Type{}),
			errors.New("TestNewGenericFunc: parameter [test5] has a invalid function signature. Expected: 'func(int)', actual: 'func(int)bool'"),
		},
		{ // Returns error when expected return types are not equal
			"TestNewGenericFunc", "test6",
			func(items ...int) bool { return false },
			simpleParamValidator(newElemTypeSlice(new([]int)), newElemTypeSlice(new(int64))),
			errors.New("TestNewGenericFunc: parameter [test6] has a invalid function signature. Expected: 'func([]int)int64', actual: 'func([]int)bool'"),
		},
		{ // Returns error when expected return types are not equal
			"TestNewGenericFunc", "test7",
			func(items ...int) bool { return false },
			simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(int64))),
			errors.New("TestNewGenericFunc: parameter [test7] has a invalid function signature. Expected: 'func(T)int64', actual: 'func([]int)bool'"),
		},
	}

	for _, test := range tests {
		_, err := newGenericFunc(test.methodName, test.paramName, test.function, test.validationFunc)
		if !(err == test.exception || err.Error() == test.exception.Error()) {
			t.Errorf("Validate expect error: %s, actual: %s", test.exception, err)
		}
	}
}

func TestCall(t *testing.T) {
	tests := []struct {
		methodName     string
		paramName      string
		function       interface{}
		validationFunc func(*functionCache) error
		fnParameter    interface{}
		result         interface{}
		exception      error
	}{
		{ // A valid function and parameters
			"TestCall", "test1",
			func(i int) int { return i * 3 },
			simpleParamValidator(newElemTypeSlice(new(genericType)), newElemTypeSlice(new(int))),
			3,
			9,
			nil,
		},
		{ // Returns error when the required type doesn't match with the specification
			"TestCall", "test2",
			func(i int) int { return i * 3 },
			simpleParamValidator(newElemTypeSlice(new(int)), newElemTypeSlice(new(int))),
			"not a int",
			9,
			errors.New("reflect: Call using string as type int"),
		},
		{ // A valid function and parameters
			"TestCall", "test3",
			func(i int) {},
			simpleParamValidator(newElemTypeSlice(new(genericType)), []reflect.Type{}),
			3,
			nil,
			nil,
		},
	}

	for _, test := range tests {
		func() {
			defer func() {
				r := recover()
				if !(r == test.exception || r == test.exception.Error()) {
					t.Errorf("expect error: nil, actual: %s", r)
				}
			}()
			dynaFunc, err := newGenericFunc(test.methodName, test.paramName, test.function, test.validationFunc)
			if err != nil {
				t.Errorf("expect error: nil, actual: %s", err)
			}
			result := dynaFunc.Call(test.fnParameter)

			if result != nil && result != test.result {
				t.Errorf("expect result: %d, actual: %d", test.result, result)
			}
		}()
	}
}
