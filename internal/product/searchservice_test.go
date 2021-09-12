package product_test

import (
	"coretrix/internal/mocks"
	"coretrix/internal/product"
	"testing"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSearchService_SearchProductByTitle(t *testing.T) {
	redisSearchClient := new(mocks.RedisSearchClient)
	docs := []redisearch.Document{
		{
			Id:      "product:1",
			Score:   1,
			Payload: []byte{},
			Properties: map[string]interface{}{
				"id":          "1",
				"title":       "product1",
				"price":       "1.99",
				"description": "description1",
				"image":       "/images/products/1.jpg",
			},
		},
		{
			Id:      "product:2",
			Score:   1,
			Payload: []byte{},
			Properties: map[string]interface{}{
				"id":          "2",
				"title":       "product2",
				"price":       "2.99",
				"description": "description2",
				"image":       "/images/products/2.jpg",
			},
		},
	}
	redisSearchClient.On("Search", "title", 0, 2).Once().Return(docs, 2, nil)
	searchService := product.NewSearchService(redisSearchClient)
	result, total, err := searchService.SearchProductByTitle("title", 0, 2)
	assert.Nil(t, err)
	assert.Equal(t, 2, total)
	assert.Equal(t, 2, len(result))

	p1 := result[0]
	assert.Equal(t, 1, p1.ID)
	assert.Equal(t, "product1", p1.Title)
	assert.Equal(t, 1.99, p1.Price)
	assert.Equal(t, "description1", p1.Description)
	assert.Equal(t, "/images/products/1.jpg", p1.Image)

	p2 := result[1]
	assert.Equal(t, 2, p2.ID)
	assert.Equal(t, "product2", p2.Title)
	assert.Equal(t, 2.99, p2.Price)
	assert.Equal(t, "description2", p2.Description)
	assert.Equal(t, "/images/products/2.jpg", p2.Image)

	redisSearchClient.AssertExpectations(t)
}

func TestSearchService_IndexProducts(t *testing.T) {
	redisSearchClient := new(mocks.RedisSearchClient)
	redisSearchClient.On("Index", mock.Anything, mock.Anything).Once().Return(nil)
	searchService := product.NewSearchService(redisSearchClient)

	products := []product.Product{
		{
			ID:          1,
			Title:       "someTitle1",
			Price:       1.99,
			Description: "this is title 1",
			Image:       "/images/product/1.jpg",
		},
		{
			ID:          2,
			Title:       "someTitle2",
			Price:       2.99,
			Description: "this is title 2",
			Image:       "/images/product/2.jpg",
		},
	}
	err := searchService.IndexProducts(products)
	assert.Nil(t, err)

	redisSearchClient.AssertExpectations(t)

}
