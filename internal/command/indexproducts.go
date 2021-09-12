package command

import (
	"context"
	"coretrix/internal/platform"
	"coretrix/internal/product"

	"go.uber.org/zap"
)

type indexProductCmd struct {
	productService       product.Service
	productSearchService product.SearchService
	logger               platform.Logger
}

func (cmd *indexProductCmd) Run(ctx context.Context, flags []string) {
	products, err := cmd.productService.All()
	if err != nil {
		cmd.logger.Error("can not get all products", err,
			zap.String("service", "indexProductCmd"),
			zap.String("method", "Run"),
		)
		return
	}
	err = cmd.productSearchService.IndexProducts(products)
	if err != nil {
		cmd.logger.Error("can not get all products", err,
			zap.String("service", "indexProductCmd"),
			zap.String("method", "Run"),
		)
		return
	}
}

func NewIndexProductCmd(productService product.Service, productSearchService product.SearchService, logger platform.Logger) ConsoleCommand {
	return &indexProductCmd{
		productService:       productService,
		productSearchService: productSearchService,
		logger:               logger,
	}
}
