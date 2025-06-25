package paliyun

import (
	aliSls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-log-go-sdk/producer"
)

var Sls = sSls{}

type sSls struct {
	Config            *SlsConfig
	SlsProducerClient *producer.Producer
}

type SlsConfig struct {
	AccessKeyId     string // AccessKeyId
	AccessKeySecret string // AccessKeySecret
	SecurityToken   string // sts token
	Endpoint        string // 地域
	Retries         int    // 失败重试次数
}

// 创建producer
func (s *sSls) NewProducer(config *SlsConfig) (producerInstance *producer.Producer, err error) {
	producerConfig := producer.GetDefaultProducerConfig()
	producerConfig.Endpoint = config.Endpoint
	producerConfig.Retries = config.Retries

	credentialsProvider := aliSls.NewStaticCredentialsProvider(config.AccessKeyId, config.AccessKeySecret, config.SecurityToken)
	producerConfig.CredentialsProvider = credentialsProvider

	producerInstance, err = producer.NewProducer(producerConfig)

	s.Config = config
	s.SlsProducerClient = producerInstance

	return
}
