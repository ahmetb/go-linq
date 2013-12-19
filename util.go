package linq

func toInts(in []interface{}) ([]int, int, int, error) {
	dst := make([]int, len(in))
	var ok bool

	var minVal, maxVal int
	var min, max int
	var minSet, maxSet bool
	for i, v := range in {
		var r int
		if r, ok = v.(int); !ok {
			return nil, -1, -1, ErrUnsupportedType
		}
		dst[i] = r
		if r < minVal || !minSet {
			minVal = r
			min = i
			minSet = true
		}
		if r > maxVal || !maxSet {
			maxVal = r
			max = i
			maxSet = true
		}
	}
	return dst, min, max, nil
}

func toStrings(in []interface{}) ([]string, int, int, error) {
	dst := make([]string, len(in))
	var ok bool

	var minVal, maxVal string
	var min, max int
	var minSet, maxSet bool
	for i, v := range in {
		var r string
		if r, ok = v.(string); !ok {
			return nil, -1, -1, ErrUnsupportedType
		}
		dst[i] = r
		if r < minVal || !minSet {
			minVal = r
			min = i
			minSet = true
		}
		if r > maxVal || !maxSet {
			maxVal = r
			max = i
			maxSet = true
		}
	}
	return dst, min, max, nil
}

func toFloat64s(in []interface{}) ([]float64, int, int, error) {
	dst := make([]float64, len(in))
	var ok bool

	var minVal, maxVal float64
	var min, max int
	var minSet, maxSet bool
	for i, v := range in {
		var r float64
		if r, ok = v.(float64); !ok {
			return nil, -1, -1, ErrUnsupportedType
		}
		dst[i] = r
		if r < minVal || !minSet {
			minVal = r
			min = i
			minSet = true
		}
		if r > maxVal || !maxSet {
			maxVal = r
			max = i
			maxSet = true
		}
	}
	return dst, min, max, nil
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
