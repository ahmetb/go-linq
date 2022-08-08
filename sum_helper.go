package linq

import "golang.org/x/exp/constraints"

func getAdder(data any) func(any) any {
	switch data.(type) {
	case int:
		return adder[int]()
	case int8:
		return adder[int8]()
	case int16:
		return adder[int16]()
	case int32:
		return adder[int32]()
	case int64:
		return adder[int64]()
	case uint:
		return adder[uint]()
	case uint8:
		return adder[uint8]()
	case uint16:
		return adder[uint16]()
	case uint32:
		return adder[uint32]()
	case uint64:
		return adder[uint64]()
	case float32:
		return adder[float32]()
	case float64:
		return adder[float64]()
	default:
		return nil
	}
}

type Number interface {
	constraints.Integer | constraints.Float
}

func adder[T Number]() func(any) any {
	var sum T = 0
	return func(i any) any {
		sum += i.(T)
		return sum
	}
}
