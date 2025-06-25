package ptencloud

import (
	"context"

	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
)

var Sts = sSts{}

type CredentialResult = sts.CredentialResult

type sSts struct {
	Config    *StsConfig
	StsClient *sts.Client
}

type StsConfig struct {
	SecretId        string
	SecretKey       string
	Duration        int64             // 临时凭证有效期, 单位秒
	Region          string            // 地域
	Host            string            // 域名
	Scheme          string            // 协议，默认为https，公有云sts获取临时密钥不允许走http，特殊场景才需要设置http
	Policy          *StsConfig_Policy // 策略
	RoleArn         string            // 角色arn
	RoleSessionName string            // 角色session名称
	ExternalId      string            // 角色外部ID
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

// 实例化STS客户端
func (s *sSts) NewStsClient(config *StsConfig) {
	stsClient := sts.NewClient(
		config.SecretId,  // 用户的 SecretId
		config.SecretKey, // 用户的 SecretKey
		nil,
		sts.Host(config.Host),     // 设置域名
		sts.Scheme(config.Scheme), // 设置协议
	)

	s.StsClient = stsClient
	s.Config = config
}

// 获取临时密钥
func (s *sSts) GetTempKey(ctx context.Context) (resp *CredentialResult, err error) {
	opt := &sts.CredentialOptions{
		DurationSeconds: s.Config.Duration,
		Region:          s.Config.Region,
	}

	// 添加策略
	if s.Config.Policy != nil {
		statement := make([]sts.CredentialPolicyStatement, 0)
		for _, item := range s.Config.Policy.Statement {
			policyDefault := sts.CredentialPolicyStatement{
				Effect:   item.Effect,
				Action:   item.Action,
				Resource: item.Resource,
			}

			if item.Condition != nil {
				policyDefault.Condition = *item.Condition
			}

			statement = append(statement, policyDefault)
		}

		opt.Policy = &sts.CredentialPolicy{
			Version:   s.Config.Policy.Version,
			Statement: statement,
		}
	}

	// 请求临时密钥
	if len(s.Config.RoleArn) > 0 {
		opt.RoleArn = s.Config.RoleArn
		opt.RoleSessionName = s.Config.RoleSessionName
		if len(s.Config.ExternalId) > 0 {
			opt.ExternalId = s.Config.ExternalId
		}
		resp, err = s.StsClient.GetRoleCredential(opt)
	} else {
		resp, err = s.StsClient.GetCredential(opt)
	}
	return
}
