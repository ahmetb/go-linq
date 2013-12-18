package linq

func toInts(in []interface{}) []int {
	dst := make([]int, len(in))
	for i := 0; i < len(in); i++ {
		dst[i] = in[i].(int)
	}
	return dst
}

func toStrings(in []interface{}) []string {
	dst := make([]string, len(in))
	for i := 0; i < len(in); i++ {
		dst[i] = in[i].(string)
	}
	return dst
}

func toFloat64s(in []interface{}) []float64 {
	dst := make([]float64, len(in))
	for i := 0; i < len(in); i++ {
		dst[i] = in[i].(float64)
	}
	return dst
}

func intsToInterface(in []int) []interface{} {
	dst := make([]interface{}, len(in))
	for i := 0; i < len(in); i++ {
		dst[i] = in[i]
	}
	return dst
}

func float64sToInterface(in []float64) []interface{} {
	dst := make([]interface{}, len(in))
	for i := 0; i < len(in); i++ {
		dst[i] = in[i]
	}
	return dst
}

func stringsToInterface(in []string) []interface{} {
	dst := make([]interface{}, len(in))
	for i := 0; i < len(in); i++ {
		dst[i] = in[i]
	}
	return dst
}
