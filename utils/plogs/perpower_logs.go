// 自建日志上报平台
package plogs

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
)

var PerLogs = sPerLogs{}

type sPerLogs struct {
	Config *PerLogsConfig
}

type PerLogsConfig struct {
	AppId     string // 应用ID
	AppKey    string // 应用密钥
	SecretKey string
	Debug     bool
}

const (
	apiUrl = "https://perlogs.leadmea.com/api/v1/comn-receiver" // 接口地址
)

// 实例化
func (s *sPerLogs) New(config *PerLogsConfig) *sPerLogs {
	return &sPerLogs{Config: config}
}

// 日志上报
func (s *sPerLogs) Report(ctx context.Context, logTime int64, module string, logData map[string]string) (err error) {
	_, err = resty.New().
		SetDebug(s.Config.Debug).
		R().
		SetHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("%s %s", s.Config.AppId, s.Config.AppKey),
		}).
		SetHeaders(s.GenerateSign(ctx)).
		SetBody(logData).
		Post(apiUrl)
	return
}

// 签名
func (s *sPerLogs) GenerateSign(ctx context.Context) (headers map[string]string) {
	var (
		timestamp = gtime.TimestampMilliStr()
		nonce     = grand.S(10, false)
	)

	sign := gmd5.MustEncrypt(fmt.Sprintf("nonce=%s&timestamp=%s&secretkey=%s", nonce, timestamp, s.Config.SecretKey))
	sign = gstr.ToUpper(sign)
	sign = base64.StdEncoding.EncodeToString(gconv.Bytes(sign))

	headers = map[string]string{
		"Timestamp": timestamp,
		"Nonce":     nonce,
		"Sign":      sign,
	}

	return
}
