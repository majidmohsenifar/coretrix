package product_test

import (
	"coretrix/internal/di"
	"coretrix/internal/platform"
	"coretrix/internal/product"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ProductTests struct {
	*suite.Suite
	httpServer        http.Handler
	redisSearchClient platform.RedisSearchClient
}

func (t *ProductTests) SetupSuite() {
	container := di.NewContainer()
	server := container.GetHttpServer()
	t.httpServer = server.GetEngine()
	t.redisSearchClient = container.GetRedisSearchClient()

}

func (t *ProductTests) SetupTest() {
	_ = t.redisSearchClient.Delete("idTest")
	doc := redisearch.NewDocument("idTest", 1)
	doc.Set("title", "test").
		Set("id", 12345).
		Set("description", "test description").
		Set("image", "/images/products/12345.jpg").
		Set("price", 0.99)
	err := t.redisSearchClient.Index(doc)
	if err != nil {
		t.Fail(err.Error())
	}
}

func (t *ProductTests) TearDownTest() {
}

func (t *ProductTests) TearDownSuite() {}

func (t *ProductTests) TestSearch() {
	//get the cart
	res := httptest.NewRecorder()
	params := url.Values{}
	params.Set("title", "test")
	params.Set("page", "0")
	params.Set("pageSize", "10")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/product/search?"+params.Encode(), nil)
	t.httpServer.ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)

	response := struct {
		Status  bool
		Message string
		Data    product.SearchResponse
	}{}
	err := json.Unmarshal(res.Body.Bytes(), &response)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), 1, response.Data.Total)
	product1 := response.Data.Products[0]
	assert.Equal(t.T(), 12345, product1.ID)
	assert.Equal(t.T(), "test", product1.Title)
	assert.Equal(t.T(), 0.99, product1.Price)
	assert.Equal(t.T(), "test description", product1.Description)
	assert.Equal(t.T(), "http://localhost:8000/images/products/12345.jpg", product1.Image)

}

func TestProduct(t *testing.T) {
	suite.Run(t, &ProductTests{
		Suite: new(suite.Suite),
	})

}
