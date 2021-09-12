package product

import (
	"coretrix/internal/platform"
	"strconv"

	"github.com/RediSearch/redisearch-go/redisearch"
)

const (
	productPrefix = "p:"
)

type ProductSearchResult struct {
	ID          int
	Title       string
	Price       float64
	Description string
	Image       string
}

type SearchService interface {
	SearchProductByTitle(title string, page, pageSize int) ([]ProductSearchResult, int, error)
	IndexProducts([]Product) error
}

type searchService struct {
	redisSearchClient platform.RedisSearchClient
}

func (s *searchService) SearchProductByTitle(title string, page, pageSize int) ([]ProductSearchResult, int, error) {
	result := make([]ProductSearchResult, 0)
	offset := page * pageSize
	text := title + "*"
	docs, total, err := s.redisSearchClient.Search(text, offset, pageSize)
	if err != nil {
		return result, 0, err
	}
	for _, d := range docs {
		ID, _ := strconv.Atoi(d.Properties["id"].(string))
		price, _ := strconv.ParseFloat(d.Properties["price"].(string), 64)
		res := ProductSearchResult{
			ID:          ID,
			Title:       d.Properties["title"].(string),
			Price:       price,
			Description: d.Properties["description"].(string),
			Image:       d.Properties["image"].(string),
		}
		result = append(result, res)
	}
	return result, total, nil
}

func (s *searchService) IndexProducts(products []Product) error {
	docs := make([]redisearch.Document, 0)
	for _, p := range products {
		IDString := productPrefix + strconv.Itoa(p.ID)
		doc := redisearch.NewDocument(IDString, 1)
		doc.Set("title", p.Title).
			Set("id", p.ID).
			Set("description", p.Description).
			Set("image", p.Image).
			Set("price", p.Price)
		docs = append(docs, doc)
	}
	return s.redisSearchClient.Index(docs...)
}

func NewSearchService(redisSearchClient platform.RedisSearchClient) SearchService {
	return &searchService{
		redisSearchClient: redisSearchClient,
	}
}
