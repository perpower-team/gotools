// 腾讯云COS对象存储操作
package pupload

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/gogf/gf/v2/os/gfile"
	perpowerFuncs "github.com/perpower-team/gotools/v2/funcs"
	"github.com/perpower-team/gotools/v2/utils/perrors"
	"github.com/perpower-team/gotools/v2/utils/ptencloud"
	"github.com/tencentyun/cos-go-sdk-v5"
	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
)

type Object = cos.Object
type ObjectTag = cos.ObjectTaggingTag
type CredentialResult = sts.CredentialResult

type TencentCos struct {
	Uploader
	config *ptencloud.CosConfig
}

func NewTencentCos(config ptencloud.CosConfig) *TencentCos {
	return &TencentCos{
		config: &config,
	}
}

// 返回配置信息
func (t *TencentCos) Config(ctx context.Context) any {
	return *t.config
}

// 返回完整的对象文件地址
// objectKey: 对象地址
func (t *TencentCos) GetFullObjectKey(objectKey string) string {
	if len(t.config.CdnUrl) > 0 {
		return fmt.Sprintf("%s/%s", t.config.CdnUrl, objectKey)
	}
	return fmt.Sprintf("%s/%s", t.config.DefaultUrl, objectKey)
}

// 上传本地文件
// objectKey: string  文件对象
// fileName: string 本地文件路径
// hasHost: bool 返回objectKey时候加上域名地址
func (t *TencentCos) DoUpload(ctx context.Context, objectKey, filePath string, hasHost ...bool) (string, error) {
	if !gfile.Exists(filePath) {
		return "", perrors.Throwf("上传文件: %s不存在", filePath)
	}

	client := ptencloud.Cos.NewClient(t.config)

	_, _, err := client.Object.Upload(
		ctx, objectKey, filePath, &cos.MultiUploadOptions{
			PartSize:       t.config.PartSize,
			ThreadPoolSize: 3,
		},
	)

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

// 使用高级上传接口通过表单方式上传文件，上传接口根据用户文件的长度，自动切分数据
// objectKey: string  文件对象
// file: *multipart.FileHeader 本地文件
// hasHost: bool 返回objectKey时候加上域名地址
func (t *TencentCos) UploadForm(ctx context.Context, objectKey string, file *multipart.FileHeader, hasHost ...bool) (string, error) {
	client := ptencloud.Cos.NewClient(t.config)

	filePath, err := perpowerFuncs.CreateTempPath(file)
	if err != nil {
		return "", err
	}

	_, _, err = client.Object.Upload(
		ctx, objectKey, filePath, &cos.MultiUploadOptions{
			PartSize:       t.config.PartSize,
			ThreadPoolSize: 3,
		},
	)

	if err != nil {
		return "", err
	}

	if len(hasHost) > 0 && hasHost[0] {
		return t.GetFullObjectKey(objectKey), nil
	} else {
		return objectKey, nil
	}

}

// 使用简单上传接口
// objectKey: string  文件对象
// file: *multipart.FileHeader 本地文件
// hasHost: bool 返回objectKey时候加上域名地址
func (t *TencentCos) Put(ctx context.Context, objectKey string, file *multipart.FileHeader, hasHost ...bool) (string, error) {
	client := ptencloud.Cos.NewClient(t.config)

	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close() // 创建文件 defer 关闭

	_, err1 := client.Object.Put(ctx, objectKey, f, nil)
	if err1 != nil {
		return "", err1
	}

	if len(hasHost) > 0 && hasHost[0] {
		return t.GetFullObjectKey(objectKey), nil
	} else {
		return objectKey, nil
	}
}

// Delete 删除一个或多个文件对象
// objectKeys: 文件对象
func (t *TencentCos) Delete(ctx context.Context, objectKeys ...string) (err error) {
	if len(objectKeys) < 1 {
		return
	}

	client := ptencloud.Cos.NewClient(t.config)

	obs := []cos.Object{}
	for _, v := range objectKeys {
		obs = append(obs, cos.Object{Key: strings.TrimSpace(v)})
	}
	opt := &cos.ObjectDeleteMultiOptions{
		Objects: obs,
		// 布尔值，这个值决定了是否启动 Quiet 模式
		// 值为 true 启动 Quiet 模式，值为 false 则启动 Verbose 模式，默认值为 false
		// Quiet: true,
	}

	_, _, err = client.Object.DeleteMulti(ctx, opt)
	return err
}

// 下载对象到本地目录
// objectKey: string 文件对象
// filepath: string 本地文件路径
func (t *TencentCos) Download(ctx context.Context, objectKey, filepath string) error {
	client := ptencloud.Cos.NewClient(t.config)

	if ok, err := t.IsExist(ctx, objectKey); !ok || err != nil {
		return errors.New("文件不存在")
	}

	opt := &cos.MultiDownloadOptions{
		ThreadPoolSize: 5,
	}
	_, err := client.Object.Download(ctx, objectKey, filepath, opt)
	return err
}

// 下载对象到浏览器弹窗下载
// objectKey: string 文件对象
// fileName: string 下载保存的文件名
func (t *TencentCos) DownloadWeb(ctx context.Context, w http.ResponseWriter, objectKey, fileName string) {
	client := ptencloud.Cos.NewClient(t.config)

	if ok, err := t.IsExist(ctx, objectKey); !ok || err != nil {
		http.Error(w, "文件不存在", http.StatusNotFound)
		return
	}

	resp, err := client.Object.Get(ctx, objectKey, nil)
	if err != nil {
		http.Error(w, "服务器异常", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "禁止的请求", http.StatusForbidden)
		return
	}

	// 流式下载
	t.DownloadStreamFile(w, resp.Body, fileName)
}

// 初始化分片上传, 并返回uploadID
// objectKey: string  文件对象
func (t *TencentCos) InitiateMultipartUpload(ctx context.Context, objectKey string) (string, error) {
	client := ptencloud.Cos.NewClient(t.config)
	v, _, err := client.Object.InitiateMultipartUpload(ctx, objectKey, nil)
	if err != nil {
		return "", err
	}
	return v.UploadID, nil
}

// 分片上传(注意：块有最小1MB的限制，小文件不要用分块上传)
// objectKey: string  文件对象
// uploadID: string
// partNumber: int 分片号
// data: []byte 文件数据
// return: partETag, error
func (t *TencentCos) UploadPart(ctx context.Context, objectKey, uploadID string, partNumber int, data []byte) (string, error) {
	client := ptencloud.Cos.NewClient(t.config)

	// 转换数据为bytes.Reader
	byteReader := bytes.NewReader(data)
	resp, err := client.Object.UploadPart(ctx, objectKey, uploadID, partNumber, byteReader, &cos.ObjectUploadPartOptions{
		ContentLength: int64(len(data)),
	})

	partETag := resp.Header.Get("ETag")
	return partETag, err
}

// 完成分片上传
// objectKey: string  文件对象
// uploadID: string
// objectParts: []Object
// objectTags: []ObjectTag 对象标签
func (t *TencentCos) CompleteMultipartUpload(ctx context.Context, objectKey, uploadID string, objectParts []Object, objectTags ...ObjectTag) error {
	client := ptencloud.Cos.NewClient(t.config)

	opt := &cos.CompleteMultipartUploadOptions{}
	opt.Parts = objectParts
	if len(objectTags) > 0 {
		tagStr := ""
		for key, val := range objectTags {
			if key > 0 {
				tagStr += ("&" + val.Key + "=" + val.Value)
			} else {
				tagStr += (val.Key + "=" + val.Value)
			}
		}
		opt.XOptionHeader.Add("x-cos-tagging", tagStr)
	}
	_, _, err := client.Object.CompleteMultipartUpload(ctx, objectKey, uploadID, opt)
	return err
}

// 终止分片上传并删除已上传的块
// objectKey: string  文件对象
// uploadID: string
func (t *TencentCos) AbortMultipartUpload(ctx context.Context, objectKey, uploadID string) error {
	client := ptencloud.Cos.NewClient(t.config)

	_, err := client.Object.AbortMultipartUpload(ctx, objectKey, uploadID)
	return err
}

// 判断指定对象是否存在
// objectKey: string  文件对象
// return: isExist, err
func (t *TencentCos) IsExist(ctx context.Context, objectKey string) (bool, error) {
	client := ptencloud.Cos.NewClient(t.config)
	return client.Object.IsExist(ctx, objectKey)
}

// 给对象设置标签
// objectKey: string  文件对象
// objectTags: []ObjectTag 对象标签
func (t *TencentCos) PutTagging(ctx context.Context, objectKey string, objectTags ...ObjectTag) error {
	client := ptencloud.Cos.NewClient(t.config)
	opt := &cos.ObjectPutTaggingOptions{
		TagSet: objectTags,
	}
	_, err := client.Object.PutTagging(ctx, objectKey, opt)
	return err
}

// 删除对象标签
// objectKey: string  文件对象
func (t *TencentCos) DeleteTagging(ctx context.Context, objectKey string) error {
	client := ptencloud.Cos.NewClient(t.config)
	_, err := client.Object.DeleteTagging(ctx, objectKey)
	return err
}

// 查询对象标签
// objectKey: string  文件对象
func (t *TencentCos) GetTagging(ctx context.Context, objectKey string) ([]ObjectTag, error) {
	client := ptencloud.Cos.NewClient(t.config)
	resp, _, err := client.Object.GetTagging(ctx, objectKey)
	return resp.TagSet, err
}

// 获取预签名URL
// objectKey: string  文件对象
// expired: time.Duration URL有效期
// return: presignedUrl, err
func (t *TencentCos) GetPresignedURL(ctx context.Context, objectKey string, expired time.Duration) (string, error) {
	client := ptencloud.Cos.NewClient(t.config)

	if ok, err := t.IsExist(ctx, objectKey); !ok || err != nil {
		return "", errors.New("文件不存在")
	}

	presignedURL, err := client.Object.GetPresignedURL(ctx, http.MethodGet, objectKey, t.config.SecretId, t.config.SecretKey, expired, nil)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

// 根据预签名URL下载对象
// presignedUrl: string 预签名URL
// fileName: string 下载保存文件名
func (t *TencentCos) DownloadPresignedObject(w http.ResponseWriter, presignedUrl, fileName string) {
	resp, err := http.Get(presignedUrl)
	if err != nil {
		http.Error(w, "下载失败", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "禁止的请求", http.StatusForbidden)
		return
	}
	// 流式下载
	t.DownloadStreamFile(w, resp.Body, fileName)
}

// 流式下载文件
func (t *TencentCos) DownloadStreamFile(w http.ResponseWriter, respBody io.ReadCloser, fileName string) {
	w.Header().Set("Content-Type", "application/octet-stream") // 默认让浏览器下载文件
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Cache-Control", "no-cache")

	// 使用chunked分块编码
	w.Header().Set("Transfer-Encoding", "chunked")
	w.WriteHeader(http.StatusOK)

	// 此处调用Flush()确保第一个数据块能够直接发送给客户端
	w.(http.Flusher).Flush()

	buf := make([]byte, 1024*1024*5) // 设置5MB的缓冲区
	for {
		n, err := respBody.Read(buf)
		if err != nil && err != io.EOF {
			http.Error(w, "服务器异常", http.StatusInternalServerError)
			return
		}

		if n == 0 {
			// 文件已经读取完毕
			break
		}

		// 将块写入响应体
		if _, err := w.Write(buf[:n]); err != nil {
			http.Error(w, "服务器异常", http.StatusInternalServerError)
			return
		}
		// 刷新缓冲区，确保数据被立即发送
		w.(http.Flusher).Flush()
	}
}
