package plogs

import (
	"context"

	"github.com/huaweicloud/huaweicloud-lts-sdk-go/producer"
	"github.com/perpower-team/gotools/v2/utils/phuaweicloud"
)

var LtsLogs = sLtsLogs{}

type sLtsLogs struct {
	Config *LtsLogsConfig
}

type LtsLogsConfig struct {
	AccessKey       string            // AccessKeyId
	SecretAccessKey string            // AccessKeySecret
	Endpoint        string            // 节点
	Region          string            // 地域
	ProjectId       string            // 项目ID
	GroupId         string            // 日志组ID
	StreamId        string            // 日志流ID
	Retries         int               // 失败重试次数
	Callback        producer.CallBack // 发送结果回调
}

type LtsProducerResult = producer.Result

// 发送结果回调，如果要获取发送结果请实现该接口
type LtsCallback interface {
	Success(result *LtsProducerResult)
	Fail(result *LtsProducerResult)
}

// 实例化
func (s *sLtsLogs) New(config *LtsLogsConfig) *sLtsLogs {
	return &sLtsLogs{Config: config}
}

// 日志上报
func (s *sLtsLogs) Report(ctx context.Context, logTime int64, module string, logData map[string]string) (err error) {
	producerInstance, err := phuaweicloud.Lts.NewProducer(
		&phuaweicloud.LtsConfig{
			AccessKey:       s.Config.AccessKey,
			SecretAccessKey: s.Config.SecretAccessKey,
			Endpoint:        s.Config.Endpoint,
			Region:          s.Config.Region,
			ProjectId:       s.Config.ProjectId,
			Retries:         s.Config.Retries,
		},
	)
	if err != nil {
		return
	}

	producerInstance.Start() // 启动producer实例

	ParseLogMessage(&logData)

	err = s.SendLogWithCallBack(producerInstance, logTime, logData)

	// 关闭发送实例
	producerInstance.Close(60000)

	return
}

// 发送日志
func (s *sLtsLogs) SendLogWithCallBack(producerInstance *producer.Producer, logTime int64, logData map[string]string) (err error) {
	ParseLogMessage(&logData)

	log := producer.GenerateLog([]string{logData["message"]}, logData)
	realLogTime := uint32(logTime)
	log.Time = &realLogTime
	err = producerInstance.SendLogWithCallBack(s.Config.GroupId, s.Config.StreamId, log, s.Config.Callback)
	return
}

// 发送结构化日志
func (s *sLtsLogs) SendLogStructWithCallBack(producerInstance *producer.Producer, logTime int64, logData map[string]string) (err error) {
	ParseLogMessage(&logData)

	log := producer.StructLog{
		Time:     logTime,
		Contents: make([]map[string]string, 0),
	}

	log.Contents = append(log.Contents, logData)
	err = producerInstance.SendLogStructWithCallBack(s.Config.GroupId, s.Config.StreamId, &log, s.Config.Callback)

	return
}
