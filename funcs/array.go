package funcs

// 判断某个元素是否在数组内（泛型支持）
func InArray[T comparable](element T, arr []T) bool {
	check := make(map[T]struct{})

	for _, value := range arr {
		check[value] = struct{}{}
	}

	// 检查元素是否在map
	if _, ok := check[element]; ok {
		return true
	}

	return false
}

// 比较两个切片,并返回差异集（泛型支持）
func DiffArray[T comparable](slice1, slice2 []T) []T {
	diffMap := make(map[T]bool)
	var diffSlice []T

	// 将 slice2 的所有元素放入一个 map 中，便于查找
	for _, v := range slice2 {
		diffMap[v] = true
	}

	// 遍历 slice1，检查元素是否在 map 中不存在，如果不存在则添加到结果切片
	for _, v := range slice1 {
		if !diffMap[v] {
			diffSlice = append(diffSlice, v)
		}
	}

	return diffSlice
}
