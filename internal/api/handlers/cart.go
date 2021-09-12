package handlers

import (
	"coretrix/internal/api/middleware"
	"coretrix/internal/order"
	"coretrix/internal/user"

	"github.com/gin-gonic/gin"
)

func AddToCart(s order.CartService) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := order.AddParams{}
		err := c.ShouldBindJSON(&p)
		if err != nil {
			errorResponse, statusCode := getValidationError(err)
			c.AbortWithStatusJSON(statusCode, errorResponse)
			return
		}
		val, _ := c.Get(middleware.UserKey)
		u := val.(*user.AutheticatedUser)
		resp, statusCode := s.Add(u, p)
		c.JSON(statusCode, resp)
	}
}

func RemoveFromCart(s order.CartService) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := order.RemoveParams{}
		err := c.ShouldBindJSON(&p)
		if err != nil {
			errorResponse, statusCode := getValidationError(err)
			c.AbortWithStatusJSON(statusCode, errorResponse)
			return
		}
		val, _ := c.Get(middleware.UserKey)
		u := val.(*user.AutheticatedUser)
		resp, statusCode := s.Remove(u, p)
		c.JSON(statusCode, resp)
	}
}

func GetCart(s order.CartService) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, _ := c.Get(middleware.UserKey)
		u := val.(*user.AutheticatedUser)
		resp, statusCode := s.Get(u)
		c.JSON(statusCode, resp)
	}
}
