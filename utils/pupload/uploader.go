// 多平台文件上传
package pupload

import (
	"context"

	"github.com/perpower-team/gotools/v2/utils/paliyun"
	"github.com/perpower-team/gotools/v2/utils/phuaweicloud"
	"github.com/perpower-team/gotools/v2/utils/ptencloud"
)

const (
	UPLOAD_DRIVER_TENCLOUD    = "cos"   // 腾讯云cos存储
	UPLOAD_DRIVER_ALIYUN      = "oss"   // 阿里云oss存储
	UPLOAD_DRIVER_HUAWEICLOUD = "obs"   // 华为云obs存储
	UPLOAD_DRIVER_LOCAL       = "local" // 本地存储
)

// 定义接口类
type Uploader interface {
	Config(ctx context.Context) any                                                            // 获取当前上传方式配置
	DoUpload(ctx context.Context, objectKey, filePath string, hasHost ...bool) (string, error) // 上传本地文件
	GetFullObjectKey(objectKey string) string                                                  // 获取完整文件url地址
	Delete(ctx context.Context, objectKeys ...string) error                                    // 删除一个或多个指定文件
	IsExist(ctx context.Context, objectKey string) (bool, error)                               // 判断文件对象是否存在
}

// 实例化
// driver: 上传平台
// config: 配置
func NewAdapter(driver string, config any) Uploader {
	switch driver {
	case UPLOAD_DRIVER_TENCLOUD:
		cfg := config.(ptencloud.CosConfig)
		return NewTencentCos(cfg)
	case UPLOAD_DRIVER_ALIYUN:
		cfg := config.(paliyun.OssConfig)
		return NewAliyunOss(cfg)
	case UPLOAD_DRIVER_HUAWEICLOUD:
		cfg := config.(phuaweicloud.ObsConfig)
		return NewHuaweiObs(cfg)
	default:
		cfg := config.(LocalConfig)
		return NewLocal(cfg)
	}
}
