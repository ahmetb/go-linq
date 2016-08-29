package linq

import "testing"

func TestIntConverter(t *testing.T) {
	tests := []struct {
		input interface{}
		want  int64
	}{
		{2, 2},
		{int8(-1), -1},
		{int16(0), 0},
		{int32(10), 10},
		{int64(5), 5},
	}

	for _, test := range tests {
		if conv := getIntConverter(test.input); conv(test.input) != test.want {
			t.Errorf("IntConverter for %v failed", test.input)
		}
	}
}

func TestUIntConverter(t *testing.T) {
	tests := []struct {
		input interface{}
		want  uint64
	}{
		{uint(2), 2},
		{uint8(1), 1},
		{uint16(0), 0},
		{uint32(10), 10},
		{uint64(5), 5},
	}

	for _, test := range tests {
		if conv := getUIntConverter(test.input); conv(test.input) != test.want {
			t.Errorf("UIntConverter for %v failed", test.input)
		}
	}
}

func TestFloatConverter(t *testing.T) {
	tests := []struct {
		input interface{}
		want  float64
	}{
		{float32(-1), -1},
		{float64(0), 0},
	}

	for _, test := range tests {
		if conv := getFloatConverter(test.input); conv(test.input) != test.want {
			t.Errorf("FloatConverter for %v failed", test.input)
		}
	}
}
