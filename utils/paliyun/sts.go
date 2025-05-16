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
	config    *StsConfig
	stsClient *sts.Client
}

// STS配置
type StsConfig struct {
	AccessKeyId     string           // AccessKeyId
	AccessKeySecret string           // AccessKeySecret
	Duration        int64            // 临时凭证有效期, 单位秒
	Arn             string           // 角色arn
	Policy          StsConfig_Policy // 策略
	RoleSessionName string           // 角色session名称
	Endpoint        string           // 访问域名
}

type StsConfig_Policy struct {
	Statement []StsConfig_Policy_Statement
	Version   string // 策略版本
}

type StsConfig_Policy_Statement struct {
	Action   []string // 策略操作
	Effect   string   // 策略效果
	Resource []string // 资源
}

type StsToken = sts.AssumeRoleResponseBodyCredentials
type CredentialTempKey = credentials.CredentialModel

// 实例化STS客户端
func (s *sSts) NewStsClientWithAK(config *StsConfig) error {
	cfg := &openapi.Config{
		AccessKeyId:     tea.String(config.AccessKeyId),
		AccessKeySecret: tea.String(config.AccessKeySecret),
		Endpoint:        tea.String(config.Endpoint),
	}

	client, err := sts.NewClient(cfg)
	if err != nil {
		return err
	}

	s.stsClient = client
	s.config = config

	return nil
}

// 获取STS TOKEN
func (s *sSts) GetStsToken(ctx context.Context) (resp *StsToken, err error) {
	policy, _ := gjson.EncodeString(s.config.Policy)
	request := &sts.AssumeRoleRequest{
		DurationSeconds: tea.Int64(s.config.Duration),
		RoleArn:         tea.String(s.config.Arn),
		Policy:          tea.String(policy),
		RoleSessionName: tea.String(s.config.RoleSessionName),
	}

	response, err := s.stsClient.AssumeRole(request)
	if err != nil || g.IsNil(response) || g.IsNil(response.Body) {
		return
	}

	resp = response.Body.Credentials

	return
}

// 获取访问凭证
func (s *sSts) GetTempKey(ctx context.Context, stsToken StsToken) (resp *CredentialTempKey, err error) {
	config := new(credentials.Config).
		SetType("sts").
		// 从环境变量中获取AccessKey Id。
		SetAccessKeyId(*stsToken.AccessKeyId).
		// 从环境变量中获取AccessKey Secret
		SetAccessKeySecret(*stsToken.AccessKeySecret).
		// 从环境变量中获取STS临时凭证。
		SetSecurityToken(*stsToken.SecurityToken)

	stsCredential, err := credentials.NewCredential(config)
	if err != nil {
		return
	}
	resp, err = stsCredential.GetCredential()

	return
}
