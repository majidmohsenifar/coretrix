package order_test

import (
	"coretrix/internal/mocks"
	"coretrix/internal/order"
	"coretrix/internal/platform"
	"coretrix/internal/product"
	"coretrix/internal/user"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCartService_ADD_ProductNotFound(t *testing.T) {
	cartStorage := new(mocks.CartStorage)
	productService := new(mocks.ProductService)
	productService.On("FindByID", 1).Once().Return(product.Product{}, platform.ErrProductNotFound)
	logger := new(mocks.Logger)
	cartService := order.NewCartService(cartStorage, productService, logger)
	u := &user.AutheticatedUser{
		ID:       1,
		Username: "username",
		IsGuest:  false,
	}
	params := order.AddParams{
		ID:       1,
		Quantity: 2,
	}
	res, httpStatus := cartService.Add(u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, httpStatus)
	assert.Equal(t, "product not found", res.Message)

	productService.AssertExpectations(t)
}

func TestCartService_ADD_Successful(t *testing.T) {
	cartStorage := new(mocks.CartStorage)
	cart := order.Cart{
		Items: []order.CartItem{},
		Total: 0.0,
	}
	cartStorage.On("GetCart", 1).Once().Return(cart, nil)
	cartStorage.On("UpdateCart", 1, mock.Anything).Once().Return(nil)

	productService := new(mocks.ProductService)
	p := product.Product{
		ID:          1,
		Title:       "product1",
		Price:       1.99,
		Description: "",
		Image:       "",
	}
	productService.On("FindByID", 1).Once().Return(p, nil)
	logger := new(mocks.Logger)
	cartService := order.NewCartService(cartStorage, productService, logger)
	u := &user.AutheticatedUser{
		ID:       1,
		Username: "username",
		IsGuest:  false,
	}
	params := order.AddParams{
		ID:       1,
		Quantity: 2,
	}
	res, httpStatus := cartService.Add(u, params)
	assert.Equal(t, http.StatusOK, httpStatus)

	result, ok := res.Data.(order.CartResponse)
	if !ok {
		t.Log("can not cast result to cart response")
		t.Fail()
	}
	assert.Equal(t, 3.98, result.Cart.Total)
	assert.Equal(t, 1, len(result.Cart.Items))

	cartItem := result.Cart.Items[0]
	assert.Equal(t, 1, cartItem.ProductID)
	assert.Equal(t, "product1", cartItem.ProductTitle)
	assert.Equal(t, 1.99, cartItem.Price)
	assert.Equal(t, 2, cartItem.Quantity)

	productService.AssertExpectations(t)
	cartStorage.AssertExpectations(t)
}

func TestCartService_Remove_ProductNotFound(t *testing.T) {
	cartStorage := new(mocks.CartStorage)
	productService := new(mocks.ProductService)
	productService.On("FindByID", 1).Once().Return(product.Product{}, platform.ErrProductNotFound)
	logger := new(mocks.Logger)
	cartService := order.NewCartService(cartStorage, productService, logger)
	u := &user.AutheticatedUser{
		ID:       1,
		Username: "username",
		IsGuest:  false,
	}
	params := order.RemoveParams{
		ID: 1,
	}
	res, httpStatus := cartService.Remove(u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, httpStatus)
	assert.Equal(t, "product not found", res.Message)

	productService.AssertExpectations(t)
}

func TestCartService_Remove_Successful(t *testing.T) {
	cartStorage := new(mocks.CartStorage)
	cart := order.Cart{
		Items: []order.CartItem{
			{
				ProductID:    1,
				ProductTitle: "product1",
				Quantity:     1,
				Price:        1.99,
			},
		},
		Total: 1.99,
	}
	cartStorage.On("GetCart", 1).Once().Return(cart, nil)
	cartStorage.On("UpdateCart", 1, mock.Anything).Once().Return(nil)

	productService := new(mocks.ProductService)
	p := product.Product{
		ID:          1,
		Title:       "product1",
		Price:       1.99,
		Description: "",
		Image:       "",
	}
	productService.On("FindByID", 1).Once().Return(p, nil)
	logger := new(mocks.Logger)
	cartService := order.NewCartService(cartStorage, productService, logger)
	u := &user.AutheticatedUser{
		ID:       1,
		Username: "username",
		IsGuest:  false,
	}
	params := order.RemoveParams{
		ID: 1,
	}
	res, httpStatus := cartService.Remove(u, params)
	assert.Equal(t, http.StatusOK, httpStatus)

	result, ok := res.Data.(order.CartResponse)
	if !ok {
		t.Log("can not cast result to cart response")
		t.Fail()
	}
	assert.Equal(t, 0.0, result.Cart.Total)
	assert.Equal(t, 0, len(result.Cart.Items))

	productService.AssertExpectations(t)
	cartStorage.AssertExpectations(t)
}

func TestCartService_Get(t *testing.T) {
	cartStorage := new(mocks.CartStorage)
	cart := order.Cart{
		Items: []order.CartItem{
			{
				ProductID:    1,
				ProductTitle: "product1",
				Quantity:     1,
				Price:        1.99,
			},
		},
		Total: 1.99,
	}
	cartStorage.On("GetCart", 1).Once().Return(cart, nil)

	productService := new(mocks.ProductService)
	logger := new(mocks.Logger)
	cartService := order.NewCartService(cartStorage, productService, logger)
	u := &user.AutheticatedUser{
		ID:       1,
		Username: "username",
		IsGuest:  false,
	}
	res, httpStatus := cartService.Get(u)
	assert.Equal(t, http.StatusOK, httpStatus)

	result, ok := res.Data.(order.CartResponse)
	if !ok {
		t.Log("can not cast result to cart response")
		t.Fail()
	}
	assert.Equal(t, 1.99, result.Cart.Total)
	assert.Equal(t, 1, len(result.Cart.Items))

	cartItem := result.Cart.Items[0]
	assert.Equal(t, 1, cartItem.ProductID)
	assert.Equal(t, "product1", cartItem.ProductTitle)
	assert.Equal(t, 1.99, cartItem.Price)
	assert.Equal(t, 1, cartItem.Quantity)
}
