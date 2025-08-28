package config

import "net/http"

const (
	Elasticsearch7 int = 7
	Elasticsearch8 int = 8
)

type Config struct {
	Nodes      []string        // 服务节点
	Username   string          // 用户名
	Password   string          // 密码
	ApiKey     string          // 秘钥
	Transport  *http.Transport // 连接池配置
	MaxRetries int             // 重试次数
	Version    int             // 主版本号
}
