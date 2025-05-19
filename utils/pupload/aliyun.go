// 阿里云对象存储
package pupload

import (
	"context"
	"fmt"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/perpower-team/gotools/utils/paliyun"
)

type aliyunOss struct {
	config *paliyun.OssConfig
}

// 返回配置信息
func (t *aliyunOss) Config(ctx context.Context) paliyun.OssConfig {
	return *t.config
}

// 简单上传
func (t *aliyunOss) Put(ctx context.Context, objectKey, localFile string) (err error) {
	paliyun.Oss.NewClient(t.config)
	// 创建上传对象的请求
	putRequest := &oss.PutObjectRequest{
		Bucket:       oss.Ptr(t.config.Bucket), // 存储空间名称
		Key:          oss.Ptr(objectKey),       // 对象名称
		StorageClass: oss.StorageClassStandard, // 指定对象的存储类型为标准存储
		Acl:          oss.ObjectACLPublicRead,  // 指定对象的访问权限为公共读
	}

	// 执行上传对象的请求
	_, err = paliyun.Oss.OssClient.PutObjectFromFile(ctx, putRequest, localFile)

	return
}

// 返回完整的文件地址
func (t *aliyunOss) GetFullObjectKey(objectKey string) (fullObjectKey string) {
	if len(t.config.CdnUrl) > 0 {
		fullObjectKey = fmt.Sprintf("%s/%s", t.config.CdnUrl, objectKey)
	} else {
		fullObjectKey = fmt.Sprintf("%s/%s", t.config.Host, objectKey)
	}
	return
}
