// 华为云OBS对象存储
package pupload

import (
	"context"
	"fmt"
	"os"

	"github.com/gogf/gf/v2/os/gfile"
	obs "github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"github.com/perpower-team/gotools/v2/utils/perrors"
	"github.com/perpower-team/gotools/v2/utils/phuaweicloud"
)

type HuaweiObs struct {
	Uploader
	config *phuaweicloud.ObsConfig
}

// 实例化客户端
func NewHuaweiObs(cfg phuaweicloud.ObsConfig) *HuaweiObs {
	return &HuaweiObs{
		config: &cfg,
	}
}

// 返回配置信息
func (t *HuaweiObs) Config(ctx context.Context) any {
	return *t.config
}

// 上传本地文件
// objectKey: string  文件对象
// fileName: string 本地文件路径
// hasHost: bool 返回objectKey时候加上域名地址
func (t *HuaweiObs) DoUpload(ctx context.Context, objectKey, filePath string, hasHost ...bool) (string, error) {
	if !gfile.Exists(filePath) {
		return "", perrors.Throwf("上传文件: %s不存在", filePath)
	}

	// 实例化客户端
	_, err := phuaweicloud.Obs.NewClient(t.config)
	if err != nil {
		return "", err
	}

	input := &obs.PutObjectInput{}
	// 指定存储桶名称
	input.Bucket = t.config.Bucket
	// 指定上传对象
	input.Key = objectKey

	fd, _ := os.Open(filePath)
	input.Body = fd
	// 流式上传本地文件
	_, err = phuaweicloud.Obs.ObsClient.PutObject(input)
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
func (t *HuaweiObs) GetFullObjectKey(objectKey string) string {
	if len(t.config.CdnUrl) > 0 {
		return fmt.Sprintf("%s/%s", t.config.CdnUrl, objectKey)
	}
	return fmt.Sprintf("%s/%s", t.config.DefaultUrl, objectKey)
}

// 删除一个或多个文件
// objectKeys: 文件对象
func (t *HuaweiObs) Delete(ctx context.Context, objectKeys ...string) (err error) {
	if len(objectKeys) < 1 {
		return
	}

	// 实例化客户端
	_, err = phuaweicloud.Obs.NewClient(t.config)
	if err != nil {
		return
	}

	input := &obs.DeleteObjectsInput{}
	// 指定存储桶名称
	input.Bucket = t.config.Bucket
	// 指定删除列表
	objects := make([]obs.ObjectToDelete, 0)
	for _, name := range objectKeys {
		objects = append(objects, obs.ObjectToDelete{Key: name})
	}

	input.Objects = objects
	// 批量删除对象
	_, err = phuaweicloud.Obs.ObsClient.DeleteObjects(input)

	return
}

// 判断文件是否存在
// objectKey: 文件对象
func (t *HuaweiObs) IsExist(ctx context.Context, objectKey string) (exist bool, err error) {
	// 实例化客户端
	_, err = phuaweicloud.Obs.NewClient(t.config)
	if err != nil {
		return false, err
	}

	input := &obs.HeadObjectInput{}
	// 指定存储桶名称
	input.Bucket = t.config.Bucket
	// 指定对象
	input.Key = objectKey
	_, err = phuaweicloud.Obs.ObsClient.HeadObject(input)
	if err != nil {
		return
	}

	exist = true
	return
}
