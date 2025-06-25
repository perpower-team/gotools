package phuaweicloud

import (
	"context"

	obs "github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
)

var Obs = sObs{}

type sObs struct {
	Config    *ObsConfig
	ObsClient *obs.ObsClient
}

type ObsConfig struct {
	AccessKey       string // 临时AK
	SecretAccessKey string // 临时SK
	SecurityToken   string // SecurityToken
	Endpoint        string
	Bucket          string // 存储桶名称
	DefaultUrl      string // 默认访问域名
	CdnUrl          string // 自定义域名/加速域名
}

// 实例化客户端
func (s *sObs) NewClient(config *ObsConfig) (client *obs.ObsClient, err error) {
	client, err = obs.New(
		s.Config.AccessKey,
		s.Config.SecretAccessKey,
		s.Config.Endpoint,
		obs.WithSecurityToken(s.Config.SecurityToken),
		obs.WithRequestContext(context.Background()),
	)
	if err != nil {
		return
	}

	s.ObsClient = client
	s.Config = config
	return
}
