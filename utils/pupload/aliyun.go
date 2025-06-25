// 阿里云OSS对象存储
package pupload

import (
	"context"
	"fmt"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/perpower-team/gotools/v2/utils/paliyun"
	"github.com/perpower-team/gotools/v2/utils/perrors"
)

type AliyunOss struct {
	Uploader
	config *paliyun.OssConfig
}

// 实例化客户端
func NewAliyunOss(cfg paliyun.OssConfig) *AliyunOss {
	return &AliyunOss{
		config: &cfg,
	}
}

// 返回配置信息
func (t *AliyunOss) Config(ctx context.Context) any {
	return *t.config
}

// 上传本地文件
// objectKey: string  文件对象
// fileName: string 本地文件路径
// hasHost: bool 返回objectKey时候加上域名地址
func (t *AliyunOss) DoUpload(ctx context.Context, objectKey, filePath string, hasHost ...bool) (string, error) {
	if !gfile.Exists(filePath) {
		return "", perrors.Throwf("上传文件: %s不存在", filePath)
	}

	// 实例化客户端
	paliyun.Oss.NewClient(t.config)

	// 创建上传对象的请求
	putRequest := &oss.PutObjectRequest{
		Bucket:       oss.Ptr(t.config.Bucket), // 存储空间名称
		Key:          oss.Ptr(objectKey),       // 对象名称
		StorageClass: oss.StorageClassStandard, // 指定对象的存储类型为标准存储
		Acl:          oss.ObjectACLPublicRead,  // 指定对象的访问权限为公共读
	}

	// 执行上传对象的请求
	_, err := paliyun.Oss.OssClient.PutObjectFromFile(ctx, putRequest, filePath)
	if err != nil {
		return "", err
	}

	needHost := true
	if len(hasHost) > 0 {
		needHost = hasHost[0]
	}

	if needHost {
		return t.GetFullObjectKey(objectKey), nil
	} else {
		return objectKey, nil
	}
}

// 返回完整的文件地址
func (t *AliyunOss) GetFullObjectKey(objectKey string) string {
	if len(t.config.CdnUrl) > 0 {
		return fmt.Sprintf("%s/%s", t.config.CdnUrl, objectKey)
	}
	return fmt.Sprintf("%s/%s", t.config.Host, objectKey)
}

// 删除一个或多个文件
// objectKeys: 文件对象
func (t *AliyunOss) Delete(ctx context.Context, objectKeys ...string) (err error) {
	if len(objectKeys) < 1 {
		return
	}

	// 实例化客户端
	paliyun.Oss.NewClient(t.config)

	// 将对象名称列表转换为切片
	DeleteObjects := make([]oss.DeleteObject, 0)
	for _, name := range objectKeys {
		DeleteObjects = append(DeleteObjects, oss.DeleteObject{Key: tea.String(strings.TrimSpace(name))})
	}

	// 创建删除多个对象的请求
	request := &oss.DeleteMultipleObjectsRequest{
		Bucket:  tea.String(t.config.Bucket), // 存储空间名称
		Objects: DeleteObjects,               // 要删除的对象列表
	}
	_, err = paliyun.Oss.OssClient.DeleteMultipleObjects(ctx, request)
	return
}

// 判断文件是否存在
// objectKey: 文件对象
func (t *AliyunOss) IsExist(ctx context.Context, objectKey string) (bool, error) {
	// 实例化客户端
	paliyun.Oss.NewClient(t.config)
	return paliyun.Oss.OssClient.IsObjectExist(ctx, t.config.Bucket, objectKey)
}
