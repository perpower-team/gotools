package paliyun

import (
	"context"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	sts "github.com/alibabacloud-go/sts-20150401/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

var Sts = sSts{}

type sSts struct {
	Config    *StsConfig
	StsClient *sts.Client
}

// STS配置
type StsConfig struct {
	AccessKeyId     string            // AccessKeyId
	AccessKeySecret string            // AccessKeySecret
	Duration        int64             // 临时凭证有效期, 单位秒
	Policy          *StsConfig_Policy // 策略
	RoleArn         string            // 角色arn
	RoleSessionName string            // 角色session名称
	Endpoint        string            // 访问域名
	Region          string            // 地域
}

type StsConfig_Policy struct {
	Statement []StsConfig_Policy_Statement
	Version   string // 策略版本
}

type StsConfig_Policy_Statement struct {
	Action    []string // 策略操作
	Effect    string   // 策略效果
	Resource  []string // 资源
	Condition *map[string]map[string]interface{}
}

type StsToken = sts.AssumeRoleResponseBodyCredentials
type CredentialTempKey = credentials.CredentialModel

// 实例化STS客户端
func (s *sSts) NewStsClientWithAK(config *StsConfig) error {
	cfg := &openapi.Config{
		AccessKeyId:     tea.String(config.AccessKeyId),
		AccessKeySecret: tea.String(config.AccessKeySecret),
		Endpoint:        tea.String(config.Endpoint),
		RegionId:        tea.String(config.Region),
	}

	client, err := sts.NewClient(cfg)
	if err != nil {
		return err
	}

	s.StsClient = client
	s.Config = config

	return nil
}

// 获取STS TOKEN
func (s *sSts) GetStsToken(ctx context.Context) (resp *StsToken, err error) {
	request := &sts.AssumeRoleRequest{
		DurationSeconds: tea.Int64(s.Config.Duration),
		RoleArn:         tea.String(s.Config.RoleArn),
		RoleSessionName: tea.String(s.Config.RoleSessionName),
	}

	// 添加策略
	if s.Config.Policy != nil {
		policy, _ := gjson.EncodeString(s.Config.Policy)

		request.Policy = tea.String(policy)
	}

	response, err := s.StsClient.AssumeRole(request)
	if err != nil || g.IsNil(response) || g.IsNil(response.Body) {
		return
	}

	resp = response.Body.Credentials

	return
}

// 获取访问凭证
func (s *sSts) GetTempKey(ctx context.Context) (resp *CredentialTempKey, err error) {
	var (
		stsToken *StsToken
		config   *credentials.Config
	)
	if len(s.Config.RoleArn) > 0 { // 角色扮演方式
		stsToken, err = s.GetStsToken(ctx)
		if err != nil || g.IsNil(stsToken) {
			return
		}

		config = new(credentials.Config).
			SetType("sts").
			// AccessKey Id。
			SetAccessKeyId(*stsToken.AccessKeyId).
			// AccessKey Secret
			SetAccessKeySecret(*stsToken.AccessKeySecret).
			// STS临时凭证。
			SetSecurityToken(*stsToken.SecurityToken)
	} else {
		config = new(credentials.Config).
			SetType("access_key").
			// AccessKey Id。
			SetAccessKeyId(s.Config.AccessKeyId).
			// AccessKey Secret
			SetAccessKeySecret(s.Config.AccessKeySecret)
	}

	stsCredential, err := credentials.NewCredential(config)
	if err != nil {
		return
	}

	resp, err = stsCredential.GetCredential()

	return
}
