package product_test

import (
	"coretrix/internal/mocks"
	"coretrix/internal/product"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_Search_Successful(t *testing.T) {
	repo := new(mocks.ProductRepository)
	searchService := new(mocks.ProductSearchService)
	searchResult := []product.ProductSearchResult{
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
	searchService.On("SearchProductByTitle", "someTitle", 0, 2).Once().Return(searchResult, 2, nil)
	configs := new(mocks.Configs)
	configs.On("GetDomain").Once().Return("localhost.test")
	logger := new(mocks.Logger)

	productService := product.NewService(repo, searchService, configs, logger)
	params := product.SearchParams{
		Title:    "someTitle",
		Page:     0,
		PageSize: 2,
	}
	res, httpStatus := productService.Search(params)
	assert.Equal(t, http.StatusOK, httpStatus)

	result, ok := res.Data.(product.SearchResponse)
	if !ok {
		t.Log("can not cast result to search response")
		t.Fail()

	}
	assert.Equal(t, 2, result.Total)
	assert.Equal(t, 2, len(result.Products))

	p1 := result.Products[0]
	assert.Equal(t, 1, p1.ID)
	assert.Equal(t, "someTitle1", p1.Title)
	assert.Equal(t, 1.99, p1.Price)
	assert.Equal(t, "this is title 1", p1.Description)
	assert.Equal(t, "localhost.test/images/product/1.jpg", p1.Image)

	p2 := result.Products[1]
	assert.Equal(t, 2, p2.ID)
	assert.Equal(t, "someTitle2", p2.Title)
	assert.Equal(t, 2.99, p2.Price)
	assert.Equal(t, "this is title 2", p2.Description)
	assert.Equal(t, "localhost.test/images/product/2.jpg", p2.Image)

	searchService.AssertExpectations(t)
	configs.AssertExpectations(t)
}

func TestService_Search_Error(t *testing.T) {
	repo := new(mocks.ProductRepository)
	searchService := new(mocks.ProductSearchService)
	searchResult := []product.ProductSearchResult{}
	searchService.On("SearchProductByTitle", "someTitle", 0, 2).Once().Return(searchResult, 0, errors.New("someErrorRelatedToSearch"))
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return()

	productService := product.NewService(repo, searchService, configs, logger)
	params := product.SearchParams{
		Title:    "someTitle",
		Page:     0,
		PageSize: 2,
	}
	res, httpStatus := productService.Search(params)
	assert.Equal(t, http.StatusInternalServerError, httpStatus)
	assert.Equal(t, "something went wrong", res.Message)

	searchService.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestService_All(t *testing.T) {
	repo := new(mocks.ProductRepository)
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
	repo.On("GetAll").Once().Return(products, nil)
	searchService := new(mocks.ProductSearchService)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	productService := product.NewService(repo, searchService, configs, logger)
	products, err := productService.All()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(products))

	p1 := products[0]
	assert.Equal(t, 1, p1.ID)
	assert.Equal(t, "someTitle1", p1.Title)
	assert.Equal(t, 1.99, p1.Price)
	assert.Equal(t, "this is title 1", p1.Description)
	assert.Equal(t, "/images/product/1.jpg", p1.Image)

	p2 := products[1]
	assert.Equal(t, 2, p2.ID)
	assert.Equal(t, "someTitle2", p2.Title)
	assert.Equal(t, 2.99, p2.Price)
	assert.Equal(t, "this is title 2", p2.Description)
	assert.Equal(t, "/images/product/2.jpg", p2.Image)

	repo.AssertExpectations(t)
}
