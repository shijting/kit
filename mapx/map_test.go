package mapx

import (
	"reflect"
	"testing"
)

func TestKeys(t *testing.T) {
	tests := []struct {
		name string
		m    map[string]interface{}
		want []string
	}{
		{"EmptyMap", map[string]interface{}{}, []string{}},
		{"NonEmptyMap", map[string]interface{}{"a": 1, "b": 2, "c": 3}, []string{"a", "b", "c"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Keys(tt.m)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValues(t *testing.T) {
	tests := []struct {
		name string
		m    map[string]int
		want []int
	}{
		{"EmptyMap", map[string]int{}, []int{}},
		{"NonEmptyMap", map[string]int{"a": 1, "b": 2, "c": 3}, []int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Values(tt.m)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Values() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlicesToMap(t *testing.T) {
	tests := []struct {
		name  string
		keys  []string
		vals  []int
		panic bool
		want  map[string]int
	}{
		{"EmptySlices", []string{}, []int{}, false, map[string]int{}},
		{"EqualLengthSlices", []string{"a", "b", "c"}, []int{1, 2, 3}, false, map[string]int{"a": 1, "b": 2, "c": 3}},
		{"MismatchedLengthSlices", []string{"a", "b", "c"}, []int{1, 2}, true, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.panic {
					t.Errorf("SlicesToMap() recover = %v, want panic: %v", r, tt.panic)
				}
			}()

			got := SlicesToMap(tt.keys, tt.vals)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SlicesToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
