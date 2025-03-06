package pcos

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/tencentyun/cos-go-sdk-v5"
)

type cWatermark struct{}

var Watermark = cWatermark{}

// 文字水印
// conf: cos配置
// text: 水印文字
// objectKey: 对象key
// filePath: 本地文件路径
func (w *cWatermark) Text(ctx context.Context, conf *CosConfig, text string, objectKey string, filePath string) (resp *cos.ImageProcessResult, err error) {
	client := NewClient(conf)

	opt := &cos.ObjectPutOptions{
		ACLHeaderOptions: nil,
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			XOptionHeader: &http.Header{},
		},
	}

	textEncode := base64.RawURLEncoding.EncodeToString([]byte(text))
	fillEncode := base64.RawURLEncoding.EncodeToString([]byte("#f6fcff"))

	pic := &cos.PicOperations{
		IsPicInfo: 1,
		Rules: []cos.PicOperationsRules{
			{
				Bucket: conf.Bucket,
				FileId: "/" + objectKey,
				Rule:   fmt.Sprintf("watermark/2/text/%s/fontsize/30/fill/%s/dissolve/80/dx/20/dy/20/batch/1/degree/315/spacing/65/shadow/40", textEncode, fillEncode),
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
func (w *cWatermark) Image(ctx context.Context, conf *CosConfig, image string, objectKey string, filePath string) (resp *cos.ImageProcessResult, err error) {
	client := NewClient(conf)

	opt := &cos.ObjectPutOptions{
		ACLHeaderOptions: nil,
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			XOptionHeader: &http.Header{},
		},
	}
	imageKey := base64.RawURLEncoding.EncodeToString([]byte(image))
	pic := &cos.PicOperations{
		IsPicInfo: 1,
		Rules: []cos.PicOperationsRules{
			{
				Bucket: conf.Bucket,
				FileId: "/" + objectKey,
				Rule:   fmt.Sprintf("watermark/1/image_key/%s/dissolve/80/batch/1/degree/315", imageKey),
			},
		},
	}

	opt.XOptionHeader.Add("Pic-Operations", cos.EncodePicOperations(pic))
	resp, _, err = client.CI.PutFromFile(ctx, objectKey, filePath, opt)
	return
}
