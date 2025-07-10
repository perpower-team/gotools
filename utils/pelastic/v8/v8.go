package v8

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	es8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/perpower-team/gotools/v2/utils/pelastic/config"
)

type V8 struct {
	client *es8.Client
}

// NewClient 实例化客户端
func (v8 *V8) NewClient(conf config.Config) (*V8, error) {
	cfg := es8.Config{
		Addresses:  conf.Nodes,
		Username:   conf.Username,
		Password:   conf.Password,
		MaxRetries: conf.MaxRetries,
	}

	if len(conf.ApiKey) > 0 {
		cfg.APIKey = conf.ApiKey
	}

	if conf.Transport != nil {
		cfg.Transport = conf.Transport
	}

	client, err := es8.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	v8.client = client

	return v8, nil
}

// ExistIndex 检测指定index 是否存在
// indexName: 索引名
func (v8 *V8) ExistIndex(ctx context.Context, indexName string) (bool, error) {
	// 创建请求对象
	req := esapi.IndicesExistsRequest{
		Index:  []string{indexName},                                       // 指定索引名称
		Header: map[string][]string{"Content-Type": {"application/json"}}, // 显式声明类型
	}

	// 执行请求
	res, err := req.Do(ctx, v8.client)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	if res.IsError() {
		// 404 表示索引不存在
		if res.StatusCode == http.StatusNotFound {
			return false, nil
		}

		return false, fmt.Errorf("错误响应: %s", res.String())
	}

	return true, nil
}

// CreateIndex 创建索引
// indexName:  索引名
// mappings:  指定索引映射,json格式
func (v8 *V8) CreateIndex(ctx context.Context, indexName string, mappings string) error {
	// 创建请求对象
	req := esapi.IndicesCreateRequest{
		Index:  indexName, // 指定索引名称
		Body:   strings.NewReader(mappings),
		Header: map[string][]string{"Content-Type": {"application/json"}}, // 显式声明类型
	}

	// 执行请求
	res, err := req.Do(ctx, v8.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// 检查响应
	if res.IsError() {
		return fmt.Errorf("错误响应: %s", res.String())
	}

	// 解析成功响应
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}

	// 检查是否确认创建
	acknowledged, ok := result["acknowledged"].(bool)
	if !ok || !acknowledged {
		return fmt.Errorf("索引创建未确认")
	}

	return nil
}

// DeleteIndex 删除索引
// indexNames:  索引集
func (v8 *V8) DeleteIndex(ctx context.Context, indexNames []string) error {
	// 构造 DELETE 请求
	req := esapi.IndicesDeleteRequest{
		Index:  indexNames,
		Header: map[string][]string{"Content-Type": {"application/json"}}, // 显式声明类型
	}

	// 发送请求
	res, err := req.Do(ctx, v8.client)
	if err != nil {
		return fmt.Errorf("请求失败: %v", err)
	}
	defer res.Body.Close()

	// 处理响应
	if res.IsError() {
		return fmt.Errorf("错误响应: %s", res.String())
	}

	return nil
}

// Write 索引文档
// indexName: 索引名
// bodyContent:  索引数据
func (v8 *V8) Write(ctx context.Context, indexName string, bodyContent interface{}, _id ...string) (string, error) {
	var (
		docId string
	)
	// 指定文档ID
	if len(_id) > 0 {
		docId = _id[0]
	}

	newBody, err := json.Marshal(bodyContent)
	if err != nil {
		return "", fmt.Errorf("json序列化文档失败: %v", err)
	}

	// 构造请求
	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: docId,
		Body:       bytes.NewReader(newBody),
		Header:     map[string][]string{"Content-Type": {"application/json"}}, // 显式声明类型
		Refresh:    "true",                                                    // 立即刷新可见
	}

	// 发送请求
	res, err := req.Do(ctx, v8.client)
	if err != nil {
		return "", fmt.Errorf("请求失败: %v", err)
	}
	defer res.Body.Close()

	// 检查响应
	if res.IsError() {
		return docId, fmt.Errorf("错误响应: %s", res.String())
	}

	// 解析结果
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return docId, fmt.Errorf("解析响应失败: %v", err)
	}

	return result["_id"].(string), nil
}

// Update 更新文档
// indexName:  索引名
// _id: 文档ID
// bodyContent:  索引数据
func (v8 *V8) Update(ctx context.Context, indexName string, _id string, bodyContent interface{}) error {
	retry := 3
	newBody, err := json.Marshal(map[string]interface{}{
		"doc": bodyContent,
	})
	if err != nil {
		return fmt.Errorf("json序列化文档失败: %v", err)
	}

	req := esapi.UpdateRequest{
		Index:           indexName,
		DocumentID:      _id,
		Body:            bytes.NewReader(newBody),
		Header:          map[string][]string{"Content-Type": {"application/json"}}, // 显式声明类型
		RetryOnConflict: &retry,                                                    // 冲突重试次数
		Refresh:         "true",                                                    // 立即刷新可见
	}

	res, err := req.Do(ctx, v8.client)
	if err != nil {
		return fmt.Errorf("请求失败: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		if res.StatusCode == http.StatusNotFound {
			return fmt.Errorf("文档%s不存在", _id)
		}
		return fmt.Errorf("错误响应: %s", res.String())
	}

	return nil
}

// Delete 删除文档
// indexName:  索引名
// _id: 文档ID
func (v8 *V8) Delete(ctx context.Context, indexName string, _id string) error {
	req := esapi.DeleteRequest{
		Index:      indexName,
		DocumentID: _id,
		Header:     map[string][]string{"Content-Type": {"application/json"}}, // 显式声明类型
		Refresh:    "true",                                                    // 立即刷新可见
	}

	res, err := req.Do(ctx, v8.client)
	if err != nil {
		return fmt.Errorf("请求失败: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		if res.StatusCode == http.StatusNotFound {
			return fmt.Errorf("文档%s不存在", _id)
		}
		return fmt.Errorf("错误响应: %s", res.String())
	}

	return nil
}

// Search 搜索文档
// indexName:  索引名
// query:  查询条件,json格式
func (v8 *V8) Search(ctx context.Context, indexName []string, query string) (map[string]interface{}, error) {
	req := esapi.SearchRequest{
		Index:  indexName,
		Header: map[string][]string{"Content-Type": {"application/json"}}, // 显式声明类型
		Query:  query,
	}

	res, err := req.Do(ctx, v8.client)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("错误响应: %s", res.String())
	}

	// 解析结果
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析失败: %v", err)
	}

	return result, nil
}
