package linq

type intConverter func(any) int64

func getIntConverter(data any) intConverter {
	switch data.(type) {
	case int:
		return func(i any) int64 {
			return int64(i.(int))
		}
	case int8:
		return func(i any) int64 {
			return int64(i.(int8))
		}
	case int16:
		return func(i any) int64 {
			return int64(i.(int16))
		}
	case int32:
		return func(i any) int64 {
			return int64(i.(int32))
		}
	}

	return func(i any) int64 {
		return i.(int64)
	}
}

type uintConverter func(any) uint64

func getUIntConverter(data any) uintConverter {
	switch data.(type) {
	case uint:
		return func(i any) uint64 {
			return uint64(i.(uint))
		}
	case uint8:
		return func(i any) uint64 {
			return uint64(i.(uint8))
		}
	case uint16:
		return func(i any) uint64 {
			return uint64(i.(uint16))
		}
	case uint32:
		return func(i any) uint64 {
			return uint64(i.(uint32))
		}
	}

	return func(i any) uint64 {
		return i.(uint64)
	}
}

type floatConverter func(any) float64

func getFloatConverter(data any) floatConverter {
	switch data.(type) {
	case float32:
		return func(i any) float64 {
			return float64(i.(float32))
		}
	}

	return func(i any) float64 {
		return i.(float64)
	}
}
