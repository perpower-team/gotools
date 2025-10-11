package phuaweicloud

import (
	"github.com/huaweicloud/huaweicloud-lts-sdk-go/producer"
)

var Lts = sLts{}

type sLts struct {
	Config            *LtsConfig
	LtsProducerClient *producer.Producer
}

type LtsConfig struct {
	AccessKey       string // AccessKeyId
	SecretAccessKey string // AccessKeySecret
	Endpoint        string // 节点
	Region          string // 地域
	ProjectId       string // 项目ID
	Retries         int    // 失败重试次数
}

// 创建producer
func (s *sLts) NewProducer(config *LtsConfig) (producerInstance *producer.Producer, err error) {
	producerConfig := producer.GetConfig()
	producerConfig.Endpoint = config.Endpoint
	producerConfig.AccessKeyID = config.AccessKey
	producerConfig.AccessKeySecret = config.SecretAccessKey
	producerConfig.RegionId = config.Region
	producerConfig.ProjectId = config.ProjectId

	producerInstance = producer.InitProducer(producerConfig)

	s.Config = config
	s.LtsProducerClient = producerInstance

	return
}
