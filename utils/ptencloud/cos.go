package ptencloud

import (
	"net/http"
	"net/url"

	"github.com/tencentyun/cos-go-sdk-v5"
)

var Cos = sCos{}

type sCos struct {
	Config    *CosConfig
	CosClient *cos.Client
}

// 腾讯云COS对象存储配置结构体
type CosConfig struct {
	AppId      string                            // 应用ID
	SecretId   string                            // 秘钥ID
	SecretKey  string                            // 秘钥值
	Bucket     string                            // 存储桶名称
	Region     string                            // 指定地域
	PartSize   int64                             // 分片上传块大小，单位MB
	DefaultUrl string                            // 默认访问地址
	CdnUrl     string                            // 自定义域名地址
	Folder     string                            // 虚拟路径
	Action     []string                          // 允许的操作权限
	Resource   []string                          // 允许的路径
	Condition  map[string]map[string]interface{} // 生效条件
}

// 初始化COS client
func (s *sCos) NewClient(config *CosConfig) *cos.Client {
	u, _ := url.Parse(config.DefaultUrl)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.SecretId,
			SecretKey: config.SecretKey,
		},
	})

	s.Config = config
	s.CosClient = client

	return client
}
