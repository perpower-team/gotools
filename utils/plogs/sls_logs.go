// 阿里云SLS日志服务
package plogs

import (
	"context"
	"time"

	"github.com/aliyun/aliyun-log-go-sdk/producer"
	"github.com/perpower-team/gotools/v2/funcs"
	"github.com/perpower-team/gotools/v2/utils/paliyun"
)

var SlsLogs = sSlsLogs{}

type sSlsLogs struct {
	Config *SlsLogsConfig
}

type SlsLogsConfig struct {
	AccessKeyId     string // AccessKeyId
	AccessKeySecret string // AccessKeySecret
	Endpoint        string // 地域
	ProjectName     string
	LogstorName     string
	Retries         int               // 失败重试次数
	Callback        producer.CallBack // 发送结果回调
}

type SlsProducerResult = producer.Result

// 发送结果回调，如果要获取发送结果请实现该接口
type SlsCallback interface {
	Success(result *SlsProducerResult)
	Fail(result *SlsProducerResult)
}

// 实例化
func (s *sSlsLogs) New(config *SlsLogsConfig) *sSlsLogs {
	return &sSlsLogs{Config: config}
}

// 日志上报(异步)
func (s *sSlsLogs) Report(ctx context.Context, logTime int64, module string, logData map[string]string) (err error) {
	producerInstance, err := paliyun.Sls.NewProducer(
		&paliyun.SlsConfig{
			AccessKeyId:     s.Config.AccessKeyId,
			AccessKeySecret: s.Config.AccessKeySecret,
			Endpoint:        s.Config.Endpoint,
			Retries:         s.Config.Retries,
		},
	)
	if err != nil {
		return
	}

	producerInstance.Start() // 启动producer实例

	// 发送日志
	ParseLogMessage(&logData)
	log := producer.GenerateLog(uint32(time.UnixMilli(logTime).Unix()), logData)
	source, _ := funcs.GetLocalIP()
	producerInstance.SendLogWithCallBack(s.Config.ProjectName, s.Config.LogstorName, module, source, log, s.Config.Callback)

	producerInstance.SafeClose() // 安全关闭

	return
}
