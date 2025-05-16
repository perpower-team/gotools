package paliyun

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

var Oss = sOss{}

type sOss struct {
	config    *OssConfig
	ossClient *oss.Client
}

// OSS存储桶配置
type OssConfig struct {
	AccessKeyId     string // AccessKeyId
	AccessKeySecret string // AccessKeySecret
	SecurityToken   string // sts token
	Endpoint        string // 访问域名,默认域名||自定义域名||加速域名
	Bucket          string // 存储桶名称
	Region          string // 指定地域
	PartSize        int64  // 分片上传块大小，单位MB
	IsCname         bool   // 是否为自定义域名，如果是自定义域名，此项必须设置为true
	Folder          string // 虚拟路径
	Duration        int64  // policy有效期, 单位秒
	Host            string
	CdnUrl          string // 自定义域名或加速域名
}

type ConfigStruct struct {
	Expiration string     `json:"expiration"`
	Conditions [][]string `json:"conditions"`
}

type PolicyToken struct {
	Host      string `json:"host"`
	Signature string `json:"signature"`
	Policy    string `json:"policy"`
	Dir       string `json:"dir"`
}

// 实例化OSS客户端
func (s *sOss) NewClient(conf *OssConfig) {
	provider := credentials.NewStaticCredentialsProvider(conf.AccessKeyId, conf.AccessKeySecret, conf.SecurityToken)
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(provider).
		WithRegion(conf.Region)

	if len(conf.Endpoint) == 0 {
		cfg = cfg.WithUseInternalEndpoint(true)
	} else {
		cfg = cfg.WithEndpoint(conf.Endpoint)
	}

	if conf.IsCname { // 如果是自定义域名
		cfg = cfg.WithUseCName(true)
	}

	s.config = conf
	s.ossClient = oss.NewClient(cfg)
}

// 生成签名
func (s *sOss) GetPolicyToken(ctx context.Context, dir string) (policy *PolicyToken, err error) {
	now := time.Now().Unix()
	expireEnd := now + s.config.Duration
	tokenExpire := s.GetGMTISO8601(expireEnd)

	config := &ConfigStruct{
		Expiration: tokenExpire,
	}

	var (
		condition []string
		result    []byte
	)
	condition = append(condition, "starts-with")
	condition = append(condition, "$key")
	condition = append(condition, dir)
	config.Conditions = append(config.Conditions, condition)

	result, err = json.Marshal(config)
	if err != nil {
		return
	}

	encodedResult := base64.StdEncoding.EncodeToString(result)
	h := hmac.New(sha1.New, []byte(s.config.AccessKeySecret))
	io.WriteString(h, encodedResult)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	policy = &PolicyToken{
		Host:      s.config.Host,
		Signature: signedStr,
		Policy:    encodedResult,
		Dir:       dir,
	}

	return
}

func (s *sOss) GetGMTISO8601(expireEnd int64) string {
	return time.Unix(expireEnd, 0).UTC().Format("2006-01-02T15:04:05Z")
}

// 返回完整的文件地址
func (s *sOss) GetFullObjectKey(objectKey string) (fullObjectKey string) {
	if len(s.config.CdnUrl) > 0 {
		fullObjectKey = fmt.Sprintf("%s/%s", s.config.CdnUrl, objectKey)
	} else {
		fullObjectKey = fmt.Sprintf("%s/%s", s.config.Host, objectKey)
	}
	return
}
