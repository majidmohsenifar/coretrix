package order_test

import (
	"bytes"
	"context"
	"coretrix/internal/di"
	"coretrix/internal/order"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CartTests struct {
	*suite.Suite
	httpServer  http.Handler
	redisClient *redis.Client
}

func (t *CartTests) SetupSuite() {
	container := di.NewContainer()
	server := container.GetHttpServer()
	t.httpServer = server.GetEngine()
	t.redisClient = container.GetRedisClient()

}

func (t *CartTests) SetupTest() {
	key := fmt.Sprintf("cart:%d", 12345)
	err := t.redisClient.Del(context.Background(), key).Err()
	if err != nil {
		t.Fail("can not delete cart from redis")
	}
}

func (t *CartTests) TearDownTest() {
	key := fmt.Sprintf("cart:%d", 12345)
	err := t.redisClient.Del(context.Background(), key).Err()
	if err != nil {
		t.Fail("can not delete cart from redis")
	}
}

func (t *CartTests) TearDownSuite() {}

func (t *CartTests) TestBasket_Add_Get_Delete() {
	res := httptest.NewRecorder()
	data := `{"id":1,"quantity":2}`
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/cart", bytes.NewReader(body))
	token := "Bearer " + "someToken"
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)

	response := struct {
		Status  bool
		Message string
		Data    order.CartResponse
	}{}
	err := json.Unmarshal(res.Body.Bytes(), &response)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), 3.98, response.Data.Cart.Total)
	assert.Equal(t.T(), 1, len(response.Data.Cart.Items))
	item1 := response.Data.Cart.Items[0]
	assert.Equal(t.T(), 1, item1.ProductID)
	assert.Equal(t.T(), 1.99, item1.Price)
	assert.Equal(t.T(), 2, item1.Quantity)

	//get the cart
	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/cart", nil)
	token = "Bearer " + "someToken"
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)

	response = struct {
		Status  bool
		Message string
		Data    order.CartResponse
	}{}
	err = json.Unmarshal(res.Body.Bytes(), &response)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), 3.98, response.Data.Cart.Total)
	assert.Equal(t.T(), 1, len(response.Data.Cart.Items))
	item1 = response.Data.Cart.Items[0]
	assert.Equal(t.T(), 1, item1.ProductID)
	assert.Equal(t.T(), 1.99, item1.Price)
	assert.Equal(t.T(), 2, item1.Quantity)

	res = httptest.NewRecorder()
	data = `{"id":1}`
	body = []byte(data)
	req = httptest.NewRequest(http.MethodDelete, "/api/v1/cart", bytes.NewReader(body))
	token = "Bearer " + "someToken"
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)

	response = struct {
		Status  bool
		Message string
		Data    order.CartResponse
	}{}
	err = json.Unmarshal(res.Body.Bytes(), &response)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), 0.0, response.Data.Cart.Total)
	assert.Equal(t.T(), 0, len(response.Data.Cart.Items))
}

func TestCart(t *testing.T) {
	suite.Run(t, &CartTests{
		Suite: new(suite.Suite),
	})

}
