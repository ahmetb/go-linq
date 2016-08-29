package linq

type comparer func(interface{}, interface{}) int

// Comparable is an interface that has to be implemented by a
// custom collection elememts in order to work with linq.
//
// Example:
//		func (f foo) CompareTo(c Comparable) int {
//			a, b := f.f1, c.(foo).f1
//
//			if a < b {
//				return -1
//			} else if a > b {
//				return 1
//			}
//
//			return 0
//		}
type Comparable interface {
	CompareTo(Comparable) int
}

func getComparer(data interface{}) comparer {
	switch data.(type) {
	case int:
		return func(x, y interface{}) int {
			a, b := x.(int), y.(int)
			if a < b {
				return -1
			} else if a > b {
				return 1
			}

			return 0
		}
	case int8:
		return func(x, y interface{}) int {
			a, b := x.(int8), y.(int8)
			if a < b {
				return -1
			} else if a > b {
				return 1
			}

			return 0
		}
	case int16:
		return func(x, y interface{}) int {
			a, b := x.(int16), y.(int16)
			if a < b {
				return -1
			} else if a > b {
				return 1
			}

			return 0
		}
	case int32:
		return func(x, y interface{}) int {
			a, b := x.(int32), y.(int32)
			if a < b {
				return -1
			} else if a > b {
				return 1
			}

			return 0
		}
	case int64:
		return func(x, y interface{}) int {
			a, b := x.(int64), y.(int64)
			if a < b {
				return -1
			} else if a > b {
				return 1
			}

			return 0
		}
	case uint:
		return func(x, y interface{}) int {
			a, b := x.(uint), y.(uint)
			if a < b {
				return -1
			} else if a > b {
				return 1
			}

			return 0
		}
	case uint8:
		return func(x, y interface{}) int {
			a, b := x.(uint8), y.(uint8)
			if a < b {
				return -1
			} else if a > b {
				return 1
			}

			return 0
		}
	case uint16:
		return func(x, y interface{}) int {
			a, b := x.(uint16), y.(uint16)
			if a < b {
				return -1
			} else if a > b {
				return 1
			}

			return 0
		}
	case uint32:
		return func(x, y interface{}) int {
			a, b := x.(uint32), y.(uint32)
			if a < b {
				return -1
			} else if a > b {
				return 1
			}

			return 0
		}
	case uint64:
		return func(x, y interface{}) int {
			a, b := x.(uint64), y.(uint64)
			if a < b {
				return -1
			} else if a > b {
				return 1
			}

			return 0
		}
	case float32:
		return func(x, y interface{}) int {
			a, b := x.(float32), y.(float32)
			if a < b {
				return -1
			} else if a > b {
				return 1
			}

			return 0
		}
	case float64:
		return func(x, y interface{}) int {
			a, b := x.(float64), y.(float64)
			if a < b {
				return -1
			} else if a > b {
				return 1
			}

			return 0
		}
	case string:
		return func(x, y interface{}) int {
			a, b := x.(string), y.(string)
			if a == b {
				return 0
			} else if a > b {
				return 1
			} else {
				return -1
			}
		}
	case bool:
		return func(x, y interface{}) int {
			a, b := x.(bool), y.(bool)
			if a == b {
				return 0
			} else if a {
				return 1
			}

			return -1
		}
	default:
		return func(x, y interface{}) int {
			a, b := x.(Comparable), y.(Comparable)
			return a.CompareTo(b)
		}
	}
}
