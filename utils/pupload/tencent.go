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
	"time"

	"github.com/gogf/gf/v2/frame/g"
	perpowerFuncs "github.com/perpower-team/gotools/funcs"
	"github.com/perpower-team/gotools/utils/pcos"
	"github.com/tencentyun/cos-go-sdk-v5"
	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
)

type Object = cos.Object
type ObjectTag = cos.ObjectTaggingTag
type CredentialResult = sts.CredentialResult

type tencentCos struct {
	config *pcos.CosConfig
}

// 返回配置信息
func (t *tencentCos) Config(ctx context.Context) pcos.CosConfig {
	return *t.config
}

// 返回完整的对象文件地址
// objectKey: 对象地址
func (t *tencentCos) GetFullObjectKey(objectKey string) string {
	if len(t.config.CdnUrl) > 0 {
		return t.config.CdnUrl + objectKey
	}
	return t.config.DefaultUrl + objectKey
}

// 使用高级上传接口上传本地文件，上传接口根据用户文件的长度，自动切分数据
// objectKey: string  文件对象
// fileName: string 本地文件路径
// hasHost: bool 返回objectKey时候加上域名地址
func (t *tencentCos) UploadLocal(ctx context.Context, objectKey, fileName string, hasHost ...bool) (string, error) {
	client := pcos.NewClient(t.config)

	_, _, err := client.Object.Upload(
		ctx, objectKey, fileName, &cos.MultiUploadOptions{
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

// 使用高级上传接口通过表单方式上传文件，上传接口根据用户文件的长度，自动切分数据
// objectKey: string  文件对象
// file: *multipart.FileHeader 本地文件
// hasHost: bool 返回objectKey时候加上域名地址
func (t *tencentCos) UploadForm(ctx context.Context, objectKey string, file *multipart.FileHeader, hasHost ...bool) (string, error) {
	client := pcos.NewClient(t.config)

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
func (t *tencentCos) Put(ctx context.Context, objectKey string, file *multipart.FileHeader, hasHost ...bool) (string, error) {
	client := pcos.NewClient(t.config)

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

// Delete 删除文件对象
// objectKey: string 文件对象
func (t *tencentCos) Delete(ctx context.Context, objectKey string) error {
	client := pcos.NewClient(t.config)

	_, err := client.Object.Delete(ctx, objectKey)
	return err
}

// 下载对象到本地目录
// objectKey: string 文件对象
// filepath: string 本地文件路径
func (t *tencentCos) Download(ctx context.Context, objectKey, filepath string) error {
	client := pcos.NewClient(t.config)

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
func (t *tencentCos) DownloadWeb(ctx context.Context, w http.ResponseWriter, objectKey, fileName string) {
	client := pcos.NewClient(t.config)

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
func (t *tencentCos) InitiateMultipartUpload(ctx context.Context, objectKey string) (string, error) {
	client := pcos.NewClient(t.config)
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
func (t *tencentCos) UploadPart(ctx context.Context, objectKey, uploadID string, partNumber int, data []byte) (string, error) {
	client := pcos.NewClient(t.config)

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
func (t *tencentCos) CompleteMultipartUpload(ctx context.Context, objectKey, uploadID string, objectParts []Object, objectTags ...ObjectTag) error {
	client := pcos.NewClient(t.config)

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
func (t *tencentCos) AbortMultipartUpload(ctx context.Context, objectKey, uploadID string) error {
	client := pcos.NewClient(t.config)

	_, err := client.Object.AbortMultipartUpload(ctx, objectKey, uploadID)
	return err
}

// 判断指定对象是否存在
// objectKey: string  文件对象
// return: isExist, err
func (t *tencentCos) IsExist(ctx context.Context, objectKey string) (bool, error) {
	client := pcos.NewClient(t.config)

	ok, err := client.Object.IsExist(ctx, objectKey)
	return ok, err
}

// 给对象设置标签
// objectKey: string  文件对象
// objectTags: []ObjectTag 对象标签
func (t *tencentCos) PutTagging(ctx context.Context, objectKey string, objectTags ...ObjectTag) error {
	client := pcos.NewClient(t.config)
	opt := &cos.ObjectPutTaggingOptions{
		TagSet: objectTags,
	}
	_, err := client.Object.PutTagging(ctx, objectKey, opt)
	return err
}

// 删除对象标签
// objectKey: string  文件对象
func (t *tencentCos) DeleteTagging(ctx context.Context, objectKey string) error {
	client := pcos.NewClient(t.config)
	_, err := client.Object.DeleteTagging(ctx, objectKey)
	return err
}

// 查询对象标签
// objectKey: string  文件对象
func (t *tencentCos) GetTagging(ctx context.Context, objectKey string) ([]ObjectTag, error) {
	client := pcos.NewClient(t.config)
	resp, _, err := client.Object.GetTagging(ctx, objectKey)
	return resp.TagSet, err
}

// 获取预签名URL
// objectKey: string  文件对象
// expired: time.Duration URL有效期
// return: presignedUrl, err
func (t *tencentCos) GetPresignedURL(ctx context.Context, objectKey string, expired time.Duration) (string, error) {
	client := pcos.NewClient(t.config)

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
func (t *tencentCos) DownloadPresignedObject(w http.ResponseWriter, presignedUrl, fileName string) {
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
func (t *tencentCos) DownloadStreamFile(w http.ResponseWriter, respBody io.ReadCloser, fileName string) {
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

// 获取临时密钥
func (t *tencentCos) GetTempKey(ctx context.Context) (resp *CredentialResult, err error) {
	stsClient := sts.NewClient(
		t.config.SecretId,  // 用户的 SecretId
		t.config.SecretKey, // 用户的 SecretKey
		nil,
		sts.Host("sts.tencentcloudapi.com"), // 设置域名, 默认域名sts.tencentcloudapi.com
		sts.Scheme("https"),                 // 设置协议, 默认为https，公有云sts获取临时密钥不允许走http，特殊场景才需要设置http
	)
	// 策略概述 https://cloud.tencent.com/document/product/436/18023
	policyDefault := sts.CredentialPolicyStatement{
		Action: []string{
			// 简单上传
			"name/cos:PostObject",
			"name/cos:PutObject",
			// 分片上传
			"name/cos:InitiateMultipartUpload",
			"name/cos:ListMultipartUploads",
			"name/cos:ListParts",
			"name/cos:UploadPart",
			"name/cos:CompleteMultipartUpload",
		},
		Effect: "allow",
		Resource: []string{
			"qcs::cos:" + t.config.Region + ":uid/" + t.config.AppId + ":" + t.config.Bucket + "/*",
		},
		Condition: map[string]map[string]interface{}{
			"string_like_if_exist": {
				"cos:content-type": []string{"image/*", "application/xml", "application/octet-stream", "video/*"},
			},
		},
	}

	if !g.IsEmpty(t.config.Action) {
		policyDefault.Action = t.config.Action
	}

	if !g.IsEmpty(t.config.Resource) {
		resource := make([]string, 0)
		for _, v := range t.config.Resource {
			resource = append(resource, fmt.Sprintf("qcs::cos:%s:uid/%s:%s/%s", t.config.Region, t.config.AppId, t.config.Bucket, v))
		}
		policyDefault.Resource = resource
	}

	if !g.IsEmpty(t.config.Condition) {
		policyDefault.Condition = t.config.Condition
	}

	opt := &sts.CredentialOptions{
		DurationSeconds: int64(time.Hour.Seconds()),
		Region:          t.config.Region,
		Policy: &sts.CredentialPolicy{
			Statement: []sts.CredentialPolicyStatement{
				policyDefault,
			},
		},
	}

	// 请求临时密钥
	resp, err = stsClient.GetCredential(opt)
	return
}
