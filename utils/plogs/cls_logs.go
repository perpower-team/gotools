// 腾讯云CLS日志服务
package plogs

import (
	"context"

	"github.com/perpower-team/gotools/v2/utils/ptencloud"
	cls "github.com/tencentcloud/tencentcloud-cls-sdk-go"
)

var ClsLogs = sClsLogs{}

type sClsLogs struct {
	Config *ClsLogsConfig
}

type ClsLogsConfig struct {
	SecretId  string       // 秘钥ID
	SecretKey string       // 秘钥值
	Endpoint  string       // 地域
	TopicId   string       // 日志主题ID
	Retries   int          // 失败重试次数
	Callback  cls.CallBack // 发送结果回调
}

type ClsProducerResult = cls.Result

// 发送结果回调，如果要获取发送结果请实现该接口
type ClsCallback interface {
	Success(result *ClsProducerResult)
	Fail(result *ClsProducerResult)
}

// 实例化
func (s *sClsLogs) New(config *ClsLogsConfig) *sClsLogs {
	return &sClsLogs{Config: config}
}

// 日志上报(异步)
func (s *sClsLogs) Report(ctx context.Context, logTime int64, module string, logData map[string]string) (err error) {
	producerInstance, err := ptencloud.Cls.NewAsyncProducerClient(
		&ptencloud.ClsProducerConfig{
			Endpoint:  s.Config.Endpoint,
			SecretId:  s.Config.SecretId,
			SecretKey: s.Config.SecretKey,
		},
	)

	// 启动异步发送程序
	producerInstance.Start()

	// 构建日志
	ParseLogMessage(&logData)
	log := cls.NewCLSLog(logTime, logData)

	producerInstance.SendLog(s.Config.TopicId, log, s.Config.Callback)

	producerInstance.Close(60000)

	return
}
