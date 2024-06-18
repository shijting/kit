package mapx

import "errors"

// Keys returns all keys in the map
func Keys[K comparable](m map[K]any) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values returns all values in the map
func Values[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// SlicesToMap returns a map from two slices
func SlicesToMap[K comparable, V any](keys []K, values []V) (map[K]V, error) {
	result := make(map[K]V)

	if len(keys) != len(values) {
		return result, errors.New("keys and values length not equal")
	}

	for i := range keys {
		result[keys[i]] = values[i]
	}

	return result, nil
}

// MapToSlice returns two slices from a map
func MapToSlice[K comparable, V any](m map[K]V) ([]K, []V) {
	keys := make([]K, 0, len(m))
	values := make([]V, 0, len(m))

	for k, v := range m {
		keys = append(keys, k)
		values = append(values, v)
	}

	return keys, values
}

// Filter returns a map containing only the key-value pairs that satisfy the filter function
func Filter[K comparable, V any](m map[K]V, filter func(key K, value V) bool) map[K]V {
	result := make(map[K]V, len(m))

	for k, v := range m {
		if filter(k, v) {
			result[k] = v
		}
	}

	return result
}

// MapToSliceWithOrder converts a map to key and value slices based on a comparator function
func MapToSliceWithOrder[K comparable, V any](m map[K]V, comparator func(key K, value V) bool) ([]K, []V) {
	keys := make([]K, 0, len(m))
	values := make([]V, 0, len(m))

	for key, value := range m {
		if comparator(key, value) {
			keys = append(keys, key)
			values = append(values, value)
		}
	}

	return keys, values
}
