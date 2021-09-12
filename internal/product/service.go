package product

import (
	"coretrix/internal/platform"
	"coretrix/internal/response"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type SearchParams struct {
	Title    string `form:"title"`
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
}

type SearchResponseProduct struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
}
type SearchResponse struct {
	Products []SearchResponseProduct `json:"products"`
	Total    int                     `json:"total"`
}

type Service interface {
	Search(params SearchParams) (apiResponse response.ApiResponse, statusCode int)
	All() ([]Product, error)
	FindByID(ID int) (Product, error)
}

type service struct {
	repo          Repository
	searchService SearchService
	configs       platform.Configs
	logger        platform.Logger
}

func (s *service) Search(params SearchParams) (apiResponse response.ApiResponse, statusCode int) {
	title := strings.Trim(params.Title, "")
	page := params.Page
	if page < 0 {
		page = 0
	}
	pageSize := params.PageSize
	if pageSize <= 0 || pageSize > 10 {
		pageSize = 10
	}

	result, total, err := s.searchService.SearchProductByTitle(title, page, pageSize)
	if err != nil {
		s.logger.Error("failed to search product by title", err,
			zap.String("service", "productService"),
			zap.String("method", "Search"),
			zap.String("title", title),
		)
		return response.Error("something went wrong", http.StatusInternalServerError, nil)
	}

	products := make([]SearchResponseProduct, 0)
	domain := s.configs.GetDomain()
	for _, r := range result {
		imageURL := domain + r.Image
		p := SearchResponseProduct{
			ID:          r.ID,
			Title:       r.Title,
			Price:       r.Price,
			Description: r.Description,
			Image:       imageURL,
		}
		products = append(products, p)

	}
	res := SearchResponse{
		Products: products,
		Total:    total,
	}
	return response.Success(res, "")
}

func (s *service) All() ([]Product, error) {
	return s.repo.GetAll()
}

func (s *service) FindByID(ID int) (Product, error) {
	return s.repo.FindByID(ID)
}

func NewService(repo Repository, searchService SearchService, configs platform.Configs, logger platform.Logger) Service {
	return &service{
		repo:          repo,
		searchService: searchService,
		configs:       configs,
		logger:        logger,
	}

}
