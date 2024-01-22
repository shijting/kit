package mapx

// Keys 返回map中所有key
// 注意：返回的切片中元素的顺序与map中元素的顺序不一定相同，这是因为map是无序的
func Keys[K comparable](m map[K]any) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values 返回map中所有value
// 注意：返回的切片中元素的顺序与map中元素的顺序不一定相同，这是因为map是无序的
func Values[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// SlicesToMap 将两个切片转换为map，第一个切片作为键，第二个切片作为值
func SlicesToMap[K comparable, V any](keys []K, values []V) map[K]V {
	result := make(map[K]V)

	// 确保键和值切片长度相等
	if len(keys) != len(values) {
		panic("键和值切片长度不一致")
	}

	for i := range keys {
		result[keys[i]] = values[i]
	}

	return result
}

// MapToSlice 将map转换为键值切片
func MapToSlice[K comparable, V any](m map[K]V) ([]K, []V) {
	keys := make([]K, 0, len(m))
	values := make([]V, 0, len(m))

	for k, v := range m {
		keys = append(keys, k)
		values = append(values, v)
	}

	return keys, values
}

// MapToSliceWithOrder 将map转换为键值切片，并按照指定的顺序返回
func MapToSliceWithOrder[K comparable, V any](m map[K]V, order []K) ([]K, []V) {
	keys := make([]K, 0, len(m))
	values := make([]V, 0, len(m))

	for _, k := range order {
		if v, ok := m[k]; ok {
			keys = append(keys, k)
			values = append(values, v)
		}
	}

	return keys, values
}

// MapToSliceWithOrderFunc 将map转换为键值切片，并按照指定的顺序返回
func MapToSliceWithOrderFunc[K comparable, V any](m map[K]V, order func(k K, v V) bool) ([]K, []V) {
	keys := make([]K, 0, len(m))
	values := make([]V, 0, len(m))

	for k, v := range m {
		if order(k, v) {
			keys = append(keys, k)
			values = append(values, v)
		}
	}

	return keys, values
}
