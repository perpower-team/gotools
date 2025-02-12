package perpowerFuncs

// array diff int
func ArrayDiff_Int(slice1, slice2 []int) []int {
	diffMap := make(map[int]bool)
	var diffSlice []int

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

// array diff uint
func ArrayDiff_Uint(slice1, slice2 []uint) []uint {
	diffMap := make(map[uint]bool)
	var diffSlice []uint

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

// array diff int32
func ArrayDiff_Int32(slice1, slice2 []int32) []int32 {
	diffMap := make(map[int32]bool)
	var diffSlice []int32

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

// array diff uint32
func ArrayDiff_Uint32(slice1, slice2 []uint32) []uint32 {
	diffMap := make(map[uint32]bool)
	var diffSlice []uint32

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

// array diff int64
func ArrayDiff_Int64(slice1, slice2 []int64) []int64 {
	diffMap := make(map[int64]bool)
	var diffSlice []int64

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

// array diff uint64
func ArrayDiff_Uint64(slice1, slice2 []uint64) []uint64 {
	diffMap := make(map[uint64]bool)
	var diffSlice []uint64

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

// array diff string
func ArrayDiff_String(slice1, slice2 []string) []string {
	diffMap := make(map[string]bool)
	var diffSlice []string

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
