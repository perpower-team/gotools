package funcs

import (
	"fmt"
	"math"
	"strconv"

	"github.com/gogf/gf/v2/os/gtime"
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

// FormatTime 格式化显示时间
func FormatTime(gt *gtime.Time) string {
	if gt == nil {
		return ""
	}
	n := gtime.Now().Timestamp()
	t := gt.Timestamp()

	var ys int64 = 31536000
	var ds int64 = 86400
	var hs int64 = 3600
	var ms int64 = 60
	var ss int64 = 1

	var rs string

	d := n - t
	switch {
	case d > ys:
		rs = fmt.Sprintf("%d年前", int(d/ys))
	case d > ds:
		rs = fmt.Sprintf("%d天前", int(d/ds))
	case d > hs:
		rs = fmt.Sprintf("%d小时前", int(d/hs))
	case d > ms:
		rs = fmt.Sprintf("%d分钟前", int(d/ms))
	case d > ss:
		rs = fmt.Sprintf("%d秒前", int(d/ss))
	default:
		rs = "刚刚"
	}

	return rs
}
