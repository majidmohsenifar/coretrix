package di

import (
	"coretrix/config"
	"coretrix/internal/api"
	"coretrix/internal/auth"
	"coretrix/internal/command"
	"coretrix/internal/order"
	"coretrix/internal/platform"
	"coretrix/internal/product"

	"github.com/go-redis/redis/v8"
)

type Container interface {
	GetHttpServer() api.HttpServer
	GetIndexProductCommand() command.ConsoleCommand
	GetRedisClient() *redis.Client
	GetRedisSearchClient() platform.RedisSearchClient
}

type container struct {
	httpServer           api.HttpServer
	logger               platform.Logger
	configs              platform.Configs
	productService       product.Service
	productRepository    product.Repository
	productSearchService product.SearchService
	redisSearchClient    platform.RedisSearchClient
	cartService          order.CartService
	authService          auth.Service
	indexProductCommand  command.ConsoleCommand
	cartStorage          order.CartStorage
	redisClient          *redis.Client
}

func (c *container) getLogger() platform.Logger {
	if c.logger == nil {
		configs := c.getConfigs()
		logger := platform.NewLogger(configs)
		c.logger = logger
	}
	return c.logger
}

func (c *container) getConfigs() platform.Configs {
	if c.configs == nil {
		viper := config.SetConfigs()
		configs := platform.NewConfigs(viper)
		c.configs = configs
	}
	return c.configs
}

func (c *container) GetHttpServer() api.HttpServer {
	if c.httpServer == nil {
		logger := c.getLogger()
		services := c.getHttpServerServices()
		httpServer := api.NewHttpServer(services, logger)
		c.httpServer = httpServer
	}
	return c.httpServer
}

func (c *container) getProductService() product.Service {
	if c.productService == nil {
		productRepo := c.getProductRepository()
		searchService := c.getProductSearchService()
		configs := c.getConfigs()
		logger := c.getLogger()
		productService := product.NewService(productRepo, searchService, configs, logger)
		c.productService = productService
	}
	return c.productService
}

func (c *container) getProductRepository() product.Repository {
	if c.productRepository == nil {
		productRepo := product.NewRepository()
		c.productRepository = productRepo
	}
	return c.productRepository
}

func (c *container) getProductSearchService() product.SearchService {
	if c.productSearchService == nil {
		redisSearchClient := c.GetRedisSearchClient()
		searchService := product.NewSearchService(redisSearchClient)
		c.productSearchService = searchService
	}
	return c.productSearchService
}

func (c *container) getHttpServerServices() api.Services {
	return api.Services{
		Configs:        c.getConfigs(),
		ProductService: c.getProductService(),
		CartService:    c.getCartService(),
		AuthService:    c.getAuthService(),
	}
}

func (c *container) GetRedisSearchClient() platform.RedisSearchClient {
	if c.redisSearchClient == nil {
		configs := c.getConfigs()
		logger := c.getLogger()
		redisSearchClient := platform.NewRedisSearchClient(configs, logger)
		c.redisSearchClient = redisSearchClient
	}
	return c.redisSearchClient
}

func (c *container) getCartService() order.CartService {
	if c.cartService == nil {
		cartStorage := c.getCartStorage()
		productService := c.getProductService()
		logger := c.getLogger()
		cartService := order.NewCartService(cartStorage, productService, logger)
		c.cartService = cartService
	}
	return c.cartService
}

func (c *container) GetIndexProductCommand() command.ConsoleCommand {
	if c.indexProductCommand == nil {
		productService := c.getProductService()
		productSearchService := c.getProductSearchService()
		logger := c.getLogger()
		cmd := command.NewIndexProductCmd(productService, productSearchService, logger)
		c.indexProductCommand = cmd
	}
	return c.indexProductCommand
}

func (c *container) getCartStorage() order.CartStorage {
	if c.cartStorage == nil {
		redisClient := c.GetRedisClient()
		cartStorage := order.NewCartStorage(redisClient)
		c.cartStorage = cartStorage
	}
	return c.cartStorage
}

func (c *container) GetRedisClient() *redis.Client {
	if c.redisClient == nil {
		configs := c.getConfigs()
		redisClient := platform.NewRedisClient(configs)
		c.redisClient = redisClient
	}
	return c.redisClient
}

func (c *container) getAuthService() auth.Service {
	if c.authService == nil {
		configs := c.getConfigs()
		authService := auth.NewAuthService(configs)
		c.authService = authService
	}
	return c.authService
}

func NewContainer() Container {
	return &container{}
}
