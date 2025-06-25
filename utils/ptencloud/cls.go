package ptencloud

import (
	cls "github.com/tencentcloud/tencentcloud-cls-sdk-go"
)

var Cls = sCls{}

type sCls struct {
	Config         *ClsProducerConfig
	ClsAsyncClient *cls.AsyncProducerClient // 异步
	ClsSyncClient  *cls.SyncProducerClient  // 同步
}

type ClsProducerConfig struct {
	SecretId  string // 秘钥ID
	SecretKey string // 秘钥值
	Endpoint  string // 地域
	Retries   int    // 失败重试次数
}

// 创建异步生产者客户端实例
func (s *sCls) NewAsyncProducerClient(config *ClsProducerConfig) (producerInstance *cls.AsyncProducerClient, err error) {
	producerConfig := cls.GetDefaultAsyncProducerClientConfig()
	// 填入域名信息，填写指引：https://cloud.tencent.com/document/product/614/18940#.E5.9F.9F.E5.90.8D，请参见链接中 API 上传日志 Tab 中的域名
	producerConfig.Endpoint = config.Endpoint

	// 并请确保密钥关联的账号具有相应的日志上传权限，权限配置指引：https://cloud.tencent.com/document/product/614/68374#.E4.BD.BF.E7.94.A8-api-.E4.B8.8A.E4.BC.A0.E6.95.B0.E6.8D.AE
	producerConfig.AccessKeyID = config.SecretId
	producerConfig.AccessKeySecret = config.SecretKey

	// 失败重试次数
	producerConfig.Retries = config.Retries

	// 创建异步生产者客户端实例
	producerInstance, err = cls.NewAsyncProducerClient(producerConfig)

	s.Config = config
	s.ClsAsyncClient = producerInstance

	return
}

// 创建同步生产者客户端实例
func (s *sCls) NewSyncProducerClient(config *ClsProducerConfig) (client *cls.SyncProducerClient, err error) {
	producerConfig := cls.GetDefaultSyncProducerClientConfig()
	// 填入域名信息，填写指引：https://cloud.tencent.com/document/product/614/18940#.E5.9F.9F.E5.90.8D，请参见链接中 API 上传日志 Tab 中的域名
	producerConfig.Endpoint = config.Endpoint

	// 并请确保密钥关联的账号具有相应的日志上传权限，权限配置指引：https://cloud.tencent.com/document/product/614/68374#.E4.BD.BF.E7.94.A8-api-.E4.B8.8A.E4.BC.A0.E6.95.B0.E6.8D.AE
	producerConfig.AccessKeyID = config.SecretId
	producerConfig.AccessKeySecret = config.SecretKey

	// 创建异步生产者客户端实例
	client, err = cls.NewSyncProducerClient(producerConfig)

	s.Config = config
	s.ClsSyncClient = client

	return
}
