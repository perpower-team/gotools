package pelastic

import (
	"context"
	"fmt"

	"github.com/perpower-team/gotools/v2/utils/pelastic/config"
	"github.com/perpower-team/gotools/v2/utils/pelastic/v7"
	v8 "github.com/perpower-team/gotools/v2/utils/pelastic/v8"
)

// Elastic 接口定义
type Elastic interface {

	// ExistIndex 检测指定index 是否存在
	// indexName: 索引名
	ExistIndex(ctx context.Context, indexName string) (bool, error)

	// CreateIndex 创建索引
	// indexName:  索引名
	// mappings:  指定索引映射,json格式
	CreateIndex(ctx context.Context, indexName string, mappings string) error

	// DeleteIndex 删除索引
	// indexNames:  索引集
	DeleteIndex(ctx context.Context, indexNames []string) error

	// Write 索引文档
	// indexName: 索引名
	// bodyContent:  索引数据
	Write(ctx context.Context, indexName string, bodyContent interface{}, _id ...string) (string, error)

	// Update 更新文档
	// indexName:  索引名
	// _id: 文档ID
	// bodyContent:  索引数据
	Update(ctx context.Context, indexName string, _id string, bodyContent interface{}) error

	// Delete 删除文档
	// indexName:  索引名
	// _id: 文档ID
	Delete(ctx context.Context, indexName string, _id string) error

	// Search 搜索文档
	// indexName:  索引名
	// query:  查询条件,json格式
	Search(ctx context.Context, indexName []string, query string) (map[string]interface{}, error)
}

func New(conf config.Config) (Elastic, error) {
	switch conf.Version {
	case config.Elasticsearch7:
		return new(v7.V7).NewClient(conf)
	case config.Elasticsearch8:
		return new(v8.V8).NewClient(conf)
	default:
		return nil, fmt.Errorf("%s", "Elasticsearch版本错误")
	}
}
