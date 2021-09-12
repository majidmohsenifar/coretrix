package platform

import (
	"github.com/RediSearch/redisearch-go/redisearch"
	"go.uber.org/zap"
)

type RedisSearchClient interface {
	Index(docs ...redisearch.Document) error
	Search(text string, offset, limit int) (docs []redisearch.Document, total int, err error)
	Delete(docID string) error
}

type redisSearchClient struct {
	client *redisearch.Client
	logger Logger
}

func (r *redisSearchClient) Index(docs ...redisearch.Document) error {
	return r.client.Index(docs...)
}

func (r *redisSearchClient) Delete(docID string) error {
	return r.client.DeleteDocument(docID)
}

func (r *redisSearchClient) Search(text string, offset, limit int) (docs []redisearch.Document, total int, err error) {
	query := redisearch.NewQuery(text).Limit(offset, limit).SetInFields("title").SetReturnFields("id", "title", "price", "description", "image")
	return r.client.Search(query)
}

func (r *redisSearchClient) createIndex(index string) {
	sc := redisearch.NewSchema(redisearch.DefaultOptions).
		AddField(redisearch.NewNumericField("id")).
		AddField(redisearch.NewNumericField("price")).
		AddField(redisearch.NewTextField("title")).
		AddField(redisearch.NewTextField("description")).
		AddField(redisearch.NewTextField("image"))

	_, err := r.client.Info()
	if err == nil {
		//here it means the index exists so we do nothing
		return

	}
	// Create the index with the given schema
	if err := r.client.CreateIndex(sc); err != nil {
		r.logger.Fatal("can not create index",
			zap.Error(err),
			zap.String("service", "redisSearchClient"),
			zap.String("method", "createIndex"),
		)
	}

}

func NewRedisSearchClient(configs Configs, logger Logger) RedisSearchClient {
	addr := configs.GetString("redissearch.dsn")
	index := configs.GetString("redisearch.index")
	client := redisearch.NewClient(addr, index)
	s := &redisSearchClient{
		client: client,
	}
	s.createIndex(index)

	return s
}
