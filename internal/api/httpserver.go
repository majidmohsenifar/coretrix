package api

import (
	"coretrix/internal/api/handlers"
	"coretrix/internal/api/middleware"
	"coretrix/internal/auth"
	"coretrix/internal/order"
	"coretrix/internal/platform"
	"coretrix/internal/product"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	DefaultAddr  = "0.0.0.0:8000"
	ReadTimeout  = 30 * time.Second
	WriteTimeout = 30 * time.Second
)

type InternalServerErrorResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type HttpServer interface {
	ListenAndServe(address string) error
	GetEngine() http.Handler
}

type httpServer struct {
	server   *http.Server
	engine   *gin.Engine
	services Services
}

type Services struct {
	Configs        platform.Configs
	ProductService product.Service
	CartService    order.CartService
	AuthService    auth.Service
}

func (s *httpServer) ListenAndServe(address string) error {
	s.server.Addr = address
	return s.server.ListenAndServe()
}

func (s *httpServer) GetEngine() http.Handler {
	return s.engine
}

func NewHttpServer(services Services, logger platform.Logger) HttpServer {
	apiRouter := gin.New()
	apiRouter.Use(globalRecover(logger, services.Configs))
	env := services.Configs.GetEnv()
	if strings.ToUpper(env) == platform.EnvProd {
		gin.SetMode(gin.ReleaseMode)
	}

	server := &http.Server{
		Addr:           DefaultAddr,
		Handler:        apiRouter,
		ReadTimeout:    ReadTimeout,
		WriteTimeout:   WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	s := &httpServer{
		server:   server,
		engine:   apiRouter,
		services: services,
	}

	s.registerRoutes()
	return s

}

func (s *httpServer) registerRoutes() {
	r := s.engine
	v1 := r.Group("/api/v1")
	{
		productRoutes := v1.Group("/product")
		{
			productRoutes.GET("/search", handlers.SearchProduct(s.services.ProductService))
		}

		cartRoutes := v1.Group("/cart")
		cartRoutes.Use(middleware.AuthMiddleware(s.services.AuthService))
		{
			cartRoutes.PUT("", handlers.AddToCart(s.services.CartService))
			cartRoutes.DELETE("", handlers.RemoveFromCart(s.services.CartService))
			cartRoutes.GET("", handlers.GetCart(s.services.CartService))
		}
	}

}

func globalRecover(logger platform.Logger, configs platform.Configs) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func(c *gin.Context) {
			message := http.StatusText(http.StatusInternalServerError)
			if rec := recover(); rec != nil {
				if configs.GetEnv() != platform.EnvProd {
					fmt.Println("rec  =>", rec)
				}
				err := errors.New("error 500")
				logger.Error(fmt.Sprintf("error  500 in global recover %v", rec), err,
					zap.String("service", "httpServer"),
					zap.String("method", "globalRecover"),
				)
				response := InternalServerErrorResponse{
					Status:  false,
					Message: message,
				}
				c.AbortWithStatusJSON(http.StatusInternalServerError, response)
			}
		}(c)
		c.Next()
	}

}
