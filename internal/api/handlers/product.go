package handlers

import (
	"coretrix/internal/product"

	"github.com/gin-gonic/gin"
)

func SearchProduct(s product.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := product.SearchParams{}
		err := c.ShouldBindQuery(&p)
		if err != nil {
			errorResponse, statusCode := getValidationError(err)
			c.AbortWithStatusJSON(statusCode, errorResponse)
			return
		}
		resp, statusCode := s.Search(p)
		c.JSON(statusCode, resp)
	}
}
