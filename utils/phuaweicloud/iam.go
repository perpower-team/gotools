package phuaweicloud

import (
	"context"
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/config"
	iam "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3"
	iamModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"
	iamRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/region"
)

var Iam = sIam{}

type sIam struct {
	Config    *IamConfig
	IamClient *iam.IamClient
}

type IamConfig struct {
	AccessKey       string            // AccessKeyId
	SecretAccessKey string            // AccessKeySecret
	ProjectId       string            // 项目ID
	Endpoint        string            // 节点
	Region          string            // 地域
	Duration        int32             // 临时凭证有效期, 单位秒, 有效时间为15min至24h
	Policy          *StsConfig_Policy // 策略信息
}

type StsConfig_Policy struct {
	Statement []StsConfig_Policy_Statement
	Version   string // 策略版本
}

type StsConfig_Policy_Statement struct {
	Action    []string // 策略操作
	Effect    string   // 策略效果
	Resource  []string // 资源
	Condition *map[string]map[string][]string
}

type IamToken = iamModel.Credential

// 实例化IAM客户端
func (s *sIam) NewIamClient(cfg *IamConfig) error {
	// 配置认证信息
	builder := basic.NewCredentialsBuilder().
		WithAk(cfg.AccessKey).
		WithSk(cfg.SecretAccessKey)

	// 如果未填写ProjectId，SDK会自动调用IAM服务查询所在region对应的项目id
	if len(cfg.ProjectId) > 0 {
		builder = builder.WithProjectId(cfg.ProjectId)
	}

	auth, err := builder.
		// 配置SDK内置的IAM服务地址
		WithIamEndpointOverride(cfg.Endpoint).
		SafeBuild()
	if err != nil {
		return err
	}

	// 使用默认配置
	httpConfig := config.DefaultHttpConfig()

	// 配置是否忽略SSL证书校验， 默认不忽略
	httpConfig.WithIgnoreSSLVerification(true)

	// 默认超时时间为120秒，可根据需要配置
	httpConfig.WithTimeout(120 * time.Second)

	// 获取可用地区
	region, err := iamRegion.SafeValueOf(cfg.Region)
	if err != nil {
		return err
	}

	// 创建服务客户端
	iamClient, err := iam.IamClientBuilder().
		// 配置地区, 如果地区不存在会导致panic
		WithRegion(region).
		// 配置认证信息
		WithCredential(auth).
		// HTTP配置
		WithHttpConfig(httpConfig).
		SafeBuild()
	if err != nil {
		return err
	}

	client := iam.NewIamClient(iamClient)

	s.Config = cfg
	s.IamClient = client

	return err
}

// 获得临时 AK、SK 和 SecurityToken
func (s *sIam) CreateTemporaryAccessKeyByToken(ctx context.Context) (*IamToken, error) {
	var (
		request *iamModel.CreateTemporaryAccessKeyByTokenRequest
	)

	// 添加策略
	if s.Config.Policy != nil {
		statement := make([]iamModel.ServiceStatement, 0)
		for _, v := range s.Config.Policy.Statement {
			effect := iamModel.GetServiceStatementEffectEnum().ALLOW
			if v.Effect == "Deny" {
				effect = iamModel.GetServiceStatementEffectEnum().DENY
			}

			stat := iamModel.ServiceStatement{
				Action:   v.Action,
				Effect:   effect,
				Resource: &v.Resource,
			}
			if v.Condition != nil {
				stat.Condition = *v.Condition
			}
			statement = append(statement, stat)
		}

		duration := s.Config.Duration
		request = &iamModel.CreateTemporaryAccessKeyByTokenRequest{
			Body: &iamModel.CreateTemporaryAccessKeyByTokenRequestBody{
				Auth: &iamModel.TokenAuth{
					Identity: &iamModel.TokenAuthIdentity{
						Methods: []iamModel.TokenAuthIdentityMethods{
							iamModel.GetTokenAuthIdentityMethodsEnum().TOKEN,
						},
						Token: &iamModel.IdentityToken{
							DurationSeconds: &duration,
						},
						Policy: &iamModel.ServicePolicy{
							Version:   s.Config.Policy.Version,
							Statement: statement,
						},
					},
				},
			},
		}
	}

	resp, err := s.IamClient.CreateTemporaryAccessKeyByToken(request)
	if err != nil || resp == nil {
		return nil, err
	}

	return resp.Credential, nil
}
