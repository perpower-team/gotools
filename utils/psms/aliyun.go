package psms

import (
	"context"
	"encoding/json"
	"fmt"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v5/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	credentials "github.com/aliyun/credentials-go/credentials"
)

type AliyunSms struct {
	Sms
	config *AliyunSmsConfig
}

// 定义传参结构体
type AliyunSmsConfig struct {
	AccessKeyId     string
	AccessKeysecret string
	Endpoint        string // 地域
	SignName        string //短信签名
}

// 实例化
func NewAliyunSms(config AliyunSmsConfig) *AliyunSms {
	return &AliyunSms{
		config: &config,
	}
}

// 发送短信
func (s *AliyunSms) Send(ctx context.Context, mobile []string, params map[string]any) (res bool, response any, err error) {
	if len(mobile) == 0 {
		return false, nil, fmt.Errorf("号码不能为空")
	}

	// 初始化client
	credentialConfig := new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId(s.config.AccessKeyId).
		SetAccessKeySecret(s.config.AccessKeysecret)

	credential, err := credentials.NewCredential(credentialConfig)
	if err != nil {
		return
	}

	config := &openapi.Config{
		Credential: credential,
		Endpoint:   tea.String(s.config.Endpoint),
	}

	client := &dysmsapi20170525.Client{}
	client, err = dysmsapi20170525.NewClient(config)

	if len(mobile) == 1 {
		resp, err2 := s.sendSingle(client, params["templateId"].(string), params["templateParam"].(string), mobile[0])
		if err2 != nil {
			return false, nil, err2
		}

		// 返回 OK 代表请求成功
		if resp != nil && resp.Body.Code == tea.String("OK") {
			return true, resp, nil
		}
		return false, resp, nil
	} else {
		var (
			jsonByte          []byte
			mobileJson        string
			templateParamJson string
			signNameJson      string
			signNameArr       = make([]string, len(mobile))
		)

		// mobie json
		jsonByte, err = json.Marshal(mobile)
		if err != nil {
			return
		}
		mobileJson = string(jsonByte)

		// signName json
		for i := 0; i < len(mobile); i++ {
			signNameArr[i] = s.config.SignName
		}
		jsonByte, err = json.Marshal(signNameArr)
		if err != nil {
			return
		}
		signNameJson = string(jsonByte)

		// templateParam json
		jsonByte, err = json.Marshal(params["templateParam"])
		if err != nil {
			return
		}
		templateParamJson = string(jsonByte)

		resp, err2 := s.sendBatch(client, params["templateId"].(string), templateParamJson, mobileJson, signNameJson)
		if err2 != nil {
			return false, nil, err2
		}

		// 返回 OK 代表请求成功
		if resp != nil && resp.Body.Code == tea.String("OK") {
			return true, resp, nil
		}
		return false, resp, nil
	}
}

// 发送单个号码
// mobile: 手机号
// templateId: 模板ID
// templateParam: 模板参数
func (s *AliyunSms) sendSingle(client *dysmsapi20170525.Client, templateId, templateParam, mobile string) (*dysmsapi20170525.SendSmsResponse, error) {
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  tea.String(mobile),
		SignName:      tea.String(s.config.SignName),
		TemplateCode:  tea.String(templateId),
		TemplateParam: tea.String(templateParam),
	}

	runtime := &util.RuntimeOptions{}
	return client.SendSmsWithOptions(sendSmsRequest, runtime)
}

// 批量发送
// mobile: 手机号
// templateId: 模板ID
// templateParam: 模板参数
func (s *AliyunSms) sendBatch(client *dysmsapi20170525.Client, templateId, templateParamJson, mobileJson, signNameJson string) (*dysmsapi20170525.SendBatchSmsResponse, error) {
	sendSmsRequest := &dysmsapi20170525.SendBatchSmsRequest{
		PhoneNumberJson:   tea.String(mobileJson),
		SignNameJson:      tea.String(signNameJson),
		TemplateCode:      tea.String(templateId),
		TemplateParamJson: tea.String(templateParamJson),
	}

	runtime := &util.RuntimeOptions{}
	return client.SendBatchSmsWithOptions(sendSmsRequest, runtime)
}
