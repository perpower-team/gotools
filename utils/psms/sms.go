// 多平台短信发送
package psms

import "context"

const (
	SMS_DRIVER_TENCENT  = "tencent"  // 腾讯云短信
	SMS_DRIVER_ALIYUN   = "aliyun"   // 阿里云短信
	SMS_DRIVER_PLASGATE = "plasgate" // plasgate 短信
)

// 定义接口
type Sms interface {
	Send(ctx context.Context, mobile []string, params map[string]any) (bool, any, error)
}

// 实例化
// driver: 发送平台
// config: 配置
func NewAdapter(driver string, config any) Sms {
	switch driver {
	case SMS_DRIVER_TENCENT:
		cfg := config.(TsmsConfig)
		return NewTencentSms(cfg)
	case SMS_DRIVER_ALIYUN:
		cfg := config.(AliyunSmsConfig)
		return NewAliyunSms(cfg)
	case SMS_DRIVER_PLASGATE:
		cfg := config.(PlasgateConfig)
		return NewPlasgateSms(cfg)
	}

	return nil
}
