// 多平台日志组件包
package plogs

import (
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
)

// 定义接口
type LoggerAdapter interface {
	// 日志上报
	// logTime: 日志时间戳
	// module: 模块名称
	// logData: 日志数据
	Report(ctx context.Context, logTime int64, module string, logData map[string]string) error
}

// 重写message字段
func ParseLogMessage(logData *map[string]string) {
	jsonMessage := make(map[string]any)
	err := gjson.DecodeTo((*logData)["message"], &jsonMessage)
	if err != nil {
		return
	}
	(*logData)["message"] = gjson.MustEncodeString(jsonMessage["message"])
}
