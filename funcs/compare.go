package funcs

import "golang.org/x/exp/constraints"

// 返回多个integer类型数值中最大的那个
func Max[T constraints.Integer](nums ...T) (max T) {
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

// MaxMin 返回切片中的最大值和最小值。
// 如果切片为空，则返回 (0, 0, false)，其中最后一个布尔值表示结果是否有效。
func MaxMin[T constraints.Ordered](nums []T) (max T, min T, ok bool) {
	if len(nums) == 0 {
		return
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

// Clamp 返回 v 在 min 和 max 之间的值。如果 v 小于 min，则返回 min；如果 v 大于 max，则返回 max。
func Clamp[T constraints.Ordered](v, min, max T) T {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
