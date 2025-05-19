// 文件上传
package pupload

import (
	"github.com/perpower-team/gotools/utils/paliyun"
	"github.com/perpower-team/gotools/utils/pcos"
)

type Uploader struct {
	TencentCos *tencentCos
	AliyunOss  *aliyunOss
}

// Instance
// conf: interface{} 上传方式配置
func Instance(conf interface{}) (c Uploader) {
	switch confType := conf.(type) { // 考虑到switch类型断言的问题，将结果分配给一个变量，否则可能会触发panic
	case pcos.CosConfig:
		cfg := conf.(pcos.CosConfig)
		c.TencentCos = &tencentCos{
			config: &cfg,
		}
		_ = confType
	case paliyun.OssConfig:
		cfg := conf.(paliyun.OssConfig)
		c.AliyunOss = &aliyunOss{
			config: &cfg,
		}
		_ = confType
	}

	return c
}
