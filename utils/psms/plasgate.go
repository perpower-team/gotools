// 柬埔寨短信平台：https://cloud.plasgate.com/
package psms

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

type PlasgateSms struct {
	Sms
	config *PlasgateConfig
}

// 定义传参结构体
type PlasgateConfig struct {
	PrivateKey string // private key
	SecretKey  string // secret key
	Sender     string //短信签名
}

// 单个号码参数结构体
type SingleSendParameter struct {
	Sender  string `json:"sender"`  // 必须的参数
	To      string `json:"to"`      // 必须的参数
	Content string `json:"content"` // 必须的参数
}

type SingleSendResponse struct {
	QueueId      *string      `json:"queue_id,omitempty"`
	MessageCount *int         `json:"message_count,omitempty"`
	Message      *interface{} `json:"message,omitempty"`
}

// 多个号码参数结构体
type BatchSendParameter struct {
	Globals     BatchSendParameter_Global      `json:"globals"`      // 必须的参数
	Messages    []BatchSendParameter_Messages  `json:"messages"`     // 必须的参数
	BatchConfig BatchSendParameter_BatchConfig `json:"batch_config"` // 额外的参数
}

type BatchSendParameter_Global struct {
	Sender string `json:"sender"`
}

type BatchSendParameter_Messages struct {
	To      []string `json:"to"`
	Content string   `json:"content"`
}

type BatchSendParameter_BatchConfig struct {
	CallbackUrl string `json:"callback_url"`
	ScheduleAt  string `json:"schedule_at"`
}

const (
	single_sender_url = "https://cloudapi.plasgate.com/rest/send"
	batch_sender_url  = "https://cloudapi.plasgate.com/rest/batch-send"
)

// 实例化
func NewPlasgateSms(config PlasgateConfig) *PlasgateSms {
	return &PlasgateSms{
		config: &config,
	}
}

// 发送短信
func (s *PlasgateSms) Send(ctx context.Context, mobile []string, params map[string]any) (res bool, response any, err error) {
	if len(mobile) == 0 {
		return false, nil, fmt.Errorf("号码不能为空")
	}

	httpClient := resty.New().R()
	if _, ok := params["debug"]; ok {
		httpClient = httpClient.SetDebug(params["debug"].(bool))
	}
	httpClient = httpClient.SetHeader("X-Secret", s.config.SecretKey)

	if len(mobile) == 1 {
		res, response, err = s.sendSingle(httpClient, params["content"].(string), mobile[0])
	} else {
		res, response, err = s.sendBatch(httpClient, params["content"].(string), mobile...)
	}

	return
}

// 发送单个手机号
func (s *PlasgateSms) sendSingle(client *resty.Request, content string, mobile string) (res bool, response *resty.Response, err error) {
	var (
		resp *SingleSendResponse
	)
	client = client.SetBody(SingleSendParameter{
		Sender:  s.config.Sender,
		To:      mobile,
		Content: content,
	})
	response, err = client.SetQueryString(fmt.Sprintf("private_key=%s", s.config.PrivateKey)).Post(single_sender_url)
	if err != nil || g.IsNil(response) {
		return
	}

	if err = gjson.Unmarshal(response.Body(), &resp); err != nil {
		return
	}

	if resp.MessageCount != nil && *resp.MessageCount > 0 {
		res = true
	}

	return
}

// 批量发送多个手机号
func (s *PlasgateSms) sendBatch(client *resty.Request, content string, mobile ...string) (res bool, response *resty.Response, err error) {
	client = client.SetBody(BatchSendParameter{
		Globals: BatchSendParameter_Global{
			Sender: s.config.Sender,
		},
		Messages: []BatchSendParameter_Messages{
			{
				To:      mobile,
				Content: content,
			},
		},
	})
	response, err = client.SetQueryString(fmt.Sprintf("private_key=%s", s.config.PrivateKey)).Post(batch_sender_url)
	if err == nil && !g.IsNil(response) {
		res = true
	}
	return
}
