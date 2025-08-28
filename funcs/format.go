package funcs

import (
	"fmt"
	"math"
	"strconv"

	"github.com/gogf/gf/v2/util/gconv"
)

// 格式化数字为两位小数,并四舍五入
func RoundFormat(amount float64) float64 {
	return math.Round(amount*100) / 100
}

// 格式化数字为两位小数的金额格式
func MoneyFormat(amount interface{}) string {
	return strconv.FormatFloat(gconv.Float64(amount), 'f', 2, 64)
}

// 数字简化处理
func ConvertNumber(num int) string {
	if num < 1000 {
		return strconv.Itoa(num)
	}

	var suffix string
	var value float64

	if num < 10000 {
		suffix = "k"
		value = float64(num) / 1000
	} else {
		suffix = "w"
		value = float64(num) / 10000
	}

	// 格式化输出，保留一位小数
	return fmt.Sprintf("%.1f%s", value, suffix)
}
