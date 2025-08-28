// 多平台日志组件包
package plogs

import "context"

// 定义接口
type LoggerAdapter interface {
	// 日志上报
	// logTime: 日志时间戳
	// module: 模块名称
	// logData: 日志数据
	Report(ctx context.Context, logTime int64, module string, logData map[string]string) error
}
