package perpowerFuncs

// 判断某个字符串元素是否在数组内
func InArrayString(element string, arr []string) bool {
	check := make(map[string]struct{})

	for _, value := range arr {
		check[value] = struct{}{}
	}

	// 检查元素是否在map
	if _, ok := check[element]; ok {
		return true
	}

	return false
}

// 判断某个int元素是否在数组内
func InArrayInt(element int, arr []int) bool {
	check := make(map[int]struct{})

	for _, value := range arr {
		check[value] = struct{}{}
	}

	// 检查元素是否在map
	if _, ok := check[element]; ok {
		return true
	}

	return false
}

// 判断某个uint元素是否在数组内
func InArrayUint(element uint, arr []uint) bool {
	check := make(map[uint]struct{})

	for _, value := range arr {
		check[value] = struct{}{}
	}

	// 检查元素是否在map
	if _, ok := check[element]; ok {
		return true
	}

	return false
}

// 判断某个int32元素是否在数组内
func InArrayInt32(element int32, arr []int32) bool {
	check := make(map[int32]struct{})

	for _, value := range arr {
		check[value] = struct{}{}
	}

	// 检查元素是否在map
	if _, ok := check[element]; ok {
		return true
	}

	return false
}

// 判断某个uint32元素是否在数组内
func InArrayUint32(element uint32, arr []uint32) bool {
	check := make(map[uint32]struct{})

	for _, value := range arr {
		check[value] = struct{}{}
	}

	// 检查元素是否在map
	if _, ok := check[element]; ok {
		return true
	}

	return false
}

// 判断某个int64元素是否在数组内
func InArrayInt64(element int64, arr []int64) bool {
	check := make(map[int64]struct{})

	for _, value := range arr {
		check[value] = struct{}{}
	}

	// 检查元素是否在map
	if _, ok := check[element]; ok {
		return true
	}

	return false
}

// 判断某个uint64元素是否在数组内
func InArrayUint64(element uint64, arr []uint64) bool {
	check := make(map[uint64]struct{})

	for _, value := range arr {
		check[value] = struct{}{}
	}

	// 检查元素是否在map
	if _, ok := check[element]; ok {
		return true
	}

	return false
}
