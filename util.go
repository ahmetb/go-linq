package linq

func toInts(in []interface{}) ([]int, error) {
	dst := make([]int, len(in))
	var ok bool
	for i, v := range in {
		var r int
		if r, ok = v.(int); !ok {
			return nil, ErrTypeMismatch
		}
		dst[i] = r
	}
	return dst, nil
}

func toStrings(in []interface{}) ([]string, error) {
	dst := make([]string, len(in))
	var ok bool
	for i, v := range in {
		var r string
		if r, ok = v.(string); !ok {
			return nil, ErrTypeMismatch
		}
		dst[i] = r
	}
	return dst, nil
}

func toFloat64s(in []interface{}) ([]float64, error) {
	dst := make([]float64, len(in))
	var ok bool
	for i, v := range in {
		var r float64
		if r, ok = v.(float64); !ok {
			return nil, ErrTypeMismatch
		}
		dst[i] = r
	}
	return dst, nil
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

func minMaxInts(in []interface{}) (int, int, error) {
	var ok bool
	var minVal, maxVal int
	var min, max int
	var minSet, maxSet bool
	for i, v := range in {
		var r int
		if r, ok = v.(int); !ok {
			return -1, -1, ErrTypeMismatch
		}
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
	return min, max, nil
}

func minMaxUints(in []interface{}) (int, int, error) {
	var ok bool
	var minVal, maxVal uint
	var min, max int
	var minSet, maxSet bool
	for i, v := range in {
		var r uint
		if r, ok = v.(uint); !ok {
			return -1, -1, ErrTypeMismatch
		}
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
	return min, max, nil
}

func minMaxFloat64s(in []interface{}) (int, int, error) {
	var ok bool
	var minVal, maxVal float64
	var min, max int
	var minSet, maxSet bool
	for i, v := range in {
		var r float64
		if r, ok = v.(float64); !ok {
			return -1, -1, ErrTypeMismatch
		}
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
	return min, max, nil
}
