package funcs

// MaxMin 返回切片中的最大值和最小值。
// 如果切片为空，则返回 (0, 0, false)，其中最后一个布尔值表示结果是否有效。
func MaxMinInt(nums []int) (max int, min int, ok bool) {
	if len(nums) == 0 {
		return 0, 0, false
	}

	max, min = nums[0], nums[0]
	for _, num := range nums[1:] {
		if num > max {
			max = num
		}
		if num < min {
			min = num
		}
	}

	return max, min, true
}

func MaxMinUint(nums []uint) (max uint, min uint, ok bool) {
	if len(nums) == 0 {
		return 0, 0, false
	}

	max, min = nums[0], nums[0]
	for _, num := range nums[1:] {
		if num > max {
			max = num
		}
		if num < min {
			min = num
		}
	}

	return max, min, true
}

func MaxMinInt64(nums []int64) (max int64, min int64, ok bool) {
	if len(nums) == 0 {
		return 0, 0, false
	}

	max, min = nums[0], nums[0]
	for _, num := range nums[1:] {
		if num > max {
			max = num
		}
		if num < min {
			min = num
		}
	}

	return max, min, true
}

func MaxMinUint64(nums []uint64) (max uint64, min uint64, ok bool) {
	if len(nums) == 0 {
		return 0, 0, false
	}

	max, min = nums[0], nums[0]
	for _, num := range nums[1:] {
		if num > max {
			max = num
		}
		if num < min {
			min = num
		}
	}

	return max, min, true
}
