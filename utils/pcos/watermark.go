package pcos

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/tencentyun/cos-go-sdk-v5"
)

type Watermark struct{}

// 文字水印
// conf: cos配置
// text: 水印文字
// objectKey: 对象key
// filePath: 本地文件路径
func (w *Watermark) Text(ctx context.Context, conf *CosConfig, text string, objectKey string, filePath string) (resp *cos.ImageProcessResult, err error) {
	client := NewClient(conf)

	opt := &cos.ObjectPutOptions{}
	textEncode := base64.RawURLEncoding.EncodeToString([]byte(text))
	pic := &cos.PicOperations{
		IsPicInfo: 0,
		Rules: []cos.PicOperationsRules{
			{
				Bucket: conf.Bucket,
				FileId: objectKey,
				Rule:   fmt.Sprintf("watermark/2/text/%s/batch/1", textEncode),
			},
		},
	}

	opt.XOptionHeader.Add("Pic-Operations", cos.EncodePicOperations(pic))
	resp, _, err = client.CI.PutFromFile(ctx, objectKey, filePath, opt)
	return
}

// 图片水印
// conf: cos配置
// image: 水印图片对象,图片在本存储桶中的路径及名称
// objectKey: 对象key
// filePath: 本地文件路径
func (w *Watermark) Image(ctx context.Context, conf *CosConfig, image string, objectKey string, filePath string) (resp *cos.ImageProcessResult, err error) {
	client := NewClient(conf)

	opt := &cos.ObjectPutOptions{}
	imageKey := base64.RawURLEncoding.EncodeToString([]byte(image))
	pic := &cos.PicOperations{
		IsPicInfo: 0,
		Rules: []cos.PicOperationsRules{
			{
				Bucket: conf.Bucket,
				FileId: objectKey,
				Rule:   fmt.Sprintf("watermark/1/image_key/%s/batch/1/degree/315", imageKey),
			},
		},
	}

	opt.XOptionHeader.Add("Pic-Operations", cos.EncodePicOperations(pic))
	resp, _, err = client.CI.PutFromFile(ctx, objectKey, filePath, opt)
	return
}
