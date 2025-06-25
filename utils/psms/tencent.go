// 腾讯云短信
package psms

import (
	"context"
	"fmt"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type TencentSms struct {
	Sms
	config *TsmsConfig
}

// 定义传参结构体
type TsmsConfig struct {
	SecretId   string
	SecretKey  string
	Endpoint   string // 地域
	AppId      string //短信应用appId
	SignName   string //短信签名
	ExtendCode string //短信码号扩展号
	SenderId   string //国际/港澳台短信 SenderId
}

// 实例化
func NewTencentSms(config TsmsConfig) *TencentSms {
	return &TencentSms{
		config: &config,
	}
}

// 发送短信
func (s *TencentSms) Send(ctx context.Context, mobile []string, params map[string]any) (bool, any, error) {
	if len(mobile) == 0 {
		return false, nil, fmt.Errorf("号码不能为空")
	}

	//实例化一个认证对象
	credential := common.NewCredential(
		s.config.SecretId,
		s.config.SecretKey,
	)

	/* 非必要步骤:
	 * 实例化一个客户端配置对象，可以指定超时时间等配置 */
	cpf := profile.NewClientProfile()

	cpf.HttpProfile.ReqMethod = "POST"

	/* 指定接入地域域名，默认就近地域接入域名为 sms.tencentcloudapi.com ，也支持指定地域域名访问，例如广州地域的域名为 sms.ap-guangzhou.tencentcloudapi.com */
	cpf.HttpProfile.Endpoint = s.config.Endpoint

	/* SDK默认用TC3-HMAC-SHA256进行签名，非必要请不要修改这个字段 */
	cpf.SignMethod = "HmacSHA1"

	/* 实例化要请求产品(以sms为例)的client对象
	 * 第二个参数是地域信息，可以直接填写字符串ap-guangzhou，支持的地域列表参考 https://cloud.tencent.com/document/api/382/52071#.E5.9C.B0.E5.9F.9F.E5.88.97.E8.A1.A8 */
	client, _ := sms.NewClient(credential, "ap-guangzhou", cpf)

	// 实例化一个请求对象，根据调用的接口和实际情况，可以进一步设置请求参数
	request := sms.NewSendSmsRequest()

	/* 短信应用ID: 短信SdkAppId在 [短信控制台] 添加应用后生成的实际SdkAppId，示例如1400006666 */
	// 应用 ID 可前往 [短信控制台](https://console.cloud.tencent.com/smsv2/app-manage) 查看
	request.SmsSdkAppId = common.StringPtr(s.config.AppId)

	/* 短信签名内容: 使用 UTF-8 编码，必须填写已审核通过的签名 */
	// 签名信息可前往 [国内短信](https://console.cloud.tencent.com/smsv2/csms-sign) 或 [国际/港澳台短信](https://console.cloud.tencent.com/smsv2/isms-sign) 的签名管理查看
	request.SignName = common.StringPtr(s.config.SignName)

	/* 模板 ID: 必须填写已审核通过的模板 ID */
	// 模板 ID 可前往 [国内短信](https://console.cloud.tencent.com/smsv2/csms-template) 或 [国际/港澳台短信](https://console.cloud.tencent.com/smsv2/isms-template) 的正文模板管理查看
	request.TemplateId = common.StringPtr(params["templateId"].(string))

	/* 模板参数: 模板参数的个数需要与 TemplateId 对应模板的变量个数保持一致，若无模板参数，则设置为空*/
	request.TemplateParamSet = common.StringPtrs(params["templateParamSet"].([]string))

	/* 下发手机号码，采用 E.164 标准，+[国家或地区码][手机号]
	 * 示例如：+8613711112222， 其中前面有一个+号 ，86为国家码，13711112222为手机号，最多不要超过200个手机号*/
	request.PhoneNumberSet = common.StringPtrs(mobile)

	/* 短信码号扩展号（无需要可忽略）: 默认未开通，如需开通请联系 [腾讯云短信小助手] */
	if len(s.config.ExtendCode) > 0 {
		request.ExtendCode = common.StringPtr(s.config.ExtendCode)
	}

	/* 国际/港澳台短信 SenderId（无需要可忽略）: 国内短信填空，默认未开通，如需开通请联系 [腾讯云短信小助手] */
	if len(s.config.SenderId) > 0 {
		request.SenderId = common.StringPtr(s.config.SenderId)
	}

	// 通过client对象调用想要访问的接口，需要传入请求对象
	resp, err := client.SendSms(request)

	if err != nil {
		return false, resp.Response, err
	}

	sendStatus := resp.Response.SendStatusSet
	if *sendStatus[0].Code != "Ok" {
		return false, resp.Response, nil
	}

	return true, resp.Response, nil
}
