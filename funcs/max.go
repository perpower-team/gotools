package funcs

// 返回多个int类型数值中最大的那个
func MaxInt(nums ...int) (max int) {
	if len(nums) == 0 {
		return
	}

	max = nums[0]
	for _, num := range nums {
		if num > max {
			max = num
		}
	}
	return
}

// 返回多个int类型数值中最大的那个
func MaxUint(nums ...uint) (max uint) {
	if len(nums) == 0 {
		return
	}

	max = nums[0]
	for _, num := range nums {
		if num > max {
			max = num
		}
	}
	return
}

// 返回多个int64类型数值中最大的那个
func MaxInt64(nums ...int64) (max int64) {
	if len(nums) == 0 {
		return
	}

	max = nums[0]
	for _, num := range nums {
		if num > max {
			max = num
		}
	}
	return
}

// 返回多个uint64类型数值中最大的那个
func MaxUint64(nums ...uint64) (max uint64) {
	if len(nums) == 0 {
		return
	}

	max = nums[0]
	for _, num := range nums {
		if num > max {
			max = num
		}
	}
	return
}

// 返回多个int32类型数值中最大的那个
func MaxInt32(nums ...int32) (max int32) {
	if len(nums) == 0 {
		return
	}

	max = nums[0]
	for _, num := range nums {
		if num > max {
			max = num
		}
	}
	return
}

// 返回多个uint32类型数值中最大的那个
func MaxUint32(nums ...uint32) (max uint32) {
	if len(nums) == 0 {
		return
	}

	max = nums[0]
	for _, num := range nums {
		if num > max {
			max = num
		}
	}
	return
}
