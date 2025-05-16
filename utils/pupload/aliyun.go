// 阿里云对象存储
package pupload

import (
	"context"

	"github.com/perpower-team/gotools/utils/paliyun"
)

type aliyunOss struct {
	config *paliyun.OssConfig
}

// 返回配置信息
func (t *aliyunOss) Config(ctx context.Context) paliyun.OssConfig {
	return *t.config
}
