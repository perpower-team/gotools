// 文件上传
package pupload

import (
	"github.com/perpower-team/gotools/utils/pcos"
)

type Uploader struct {
	TencentCos *tencentCos
}

// Instance
// conf: interface{} 上传方式配置
func Instance(conf interface{}) (c Uploader) {
	switch confType := conf.(type) { // 考虑到switch类型断言的问题，将结果分配给一个变量，否则可能会触发panic
	case pcos.CosConfig:
		cosConf := conf.(pcos.CosConfig)
		c.TencentCos = &tencentCos{
			config: &cosConf,
		}
		_ = confType
	}
	return c
}
