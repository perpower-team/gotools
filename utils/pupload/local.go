// 本地文件存储
package pupload

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/perpower-team/gotools/v2/utils/perrors"
)

type Local struct {
	Uploader
	config *LocalConfig
}

type LocalConfig struct {
	Dir  string // 本地文件存储目录
	Host string // 文件地址host
}

// 实例化客户端
func NewLocal(cfg LocalConfig) *Local {
	return &Local{
		config: &cfg,
	}
}

// 返回配置信息
func (t *Local) Config(ctx context.Context) any {
	return *t.config
}

// 上传本地文件
// objectKey: string  文件对象
// fileName: string 本地文件路径
// hasHost: bool 返回objectKey时候加上域名地址
func (t *Local) DoUpload(ctx context.Context, objectKey, filePath string, hasHost ...bool) (string, error) {
	if ok, _ := t.IsExist(ctx, filePath); !ok {
		return "", perrors.Throwf("上传文件: %s不存在", filePath)
	}

	// 构建本地文件保存路径
	localPath := fmt.Sprintf("%s/%s", t.config.Dir, objectKey)
	if !gfile.Exists(gfile.Dir(localPath)) {
		err := gfile.Mkdir(gfile.Abs(gfile.Dir(localPath)))
		if err != nil {
			return "", perrors.Throw("创建目录失败")
		}
	}

	err := gfile.CopyFile(filePath, localPath)
	if err != nil {
		return "", err
	}

	if len(hasHost) > 0 && hasHost[0] {
		return t.GetFullObjectKey(localPath), nil
	} else {
		return localPath, nil
	}
}

// 返回完整的文件地址
func (t *Local) GetFullObjectKey(objectKey string) string {
	return fmt.Sprintf("%s/%s", t.config.Host, objectKey)
}

// 删除一个或多个文件
// objectKeys: 文件对象
func (t *Local) Delete(ctx context.Context, objectKeys ...string) (err error) {
	if len(objectKeys) < 1 {
		return
	}

	for _, name := range objectKeys {
		gfile.RemoveFile(name)
	}

	return
}

// 判断文件是否存在
// objectKey: 文件对象
func (t *Local) IsExist(ctx context.Context, filepath string) (bool, error) {
	return gfile.Exists(filepath), nil
}
