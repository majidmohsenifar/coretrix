package order

import (
	"coretrix/internal/platform"
	"coretrix/internal/product"
	"coretrix/internal/response"
	"coretrix/internal/user"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

type AddParams struct {
	ID       int `json="id"`
	Quantity int `json="quantity"`
}

type RemoveParams struct {
	ID int `json="id"`
}

type CartResponse struct {
	Cart Cart `json:"cart"`
}

type CartService interface {
	Add(u *user.AutheticatedUser, params AddParams) (apiResponse response.ApiResponse, statusCode int)
	Remove(u *user.AutheticatedUser, params RemoveParams) (apiResponse response.ApiResponse, statusCode int)
	Get(u *user.AutheticatedUser) (apiResponse response.ApiResponse, statusCode int)
}

type cartService struct {
	cartStorage    CartStorage
	productService product.Service
	logger         platform.Logger
}

func (c *cartService) Add(u *user.AutheticatedUser, params AddParams) (apiResponse response.ApiResponse, statusCode int) {
	if params.Quantity <= 0 {
		return response.Error("quantity must be more than zero", http.StatusUnprocessableEntity, nil)
	}
	product, err := c.productService.FindByID(params.ID)
	if err != nil && !errors.Is(err, platform.ErrProductNotFound) {
		c.logger.Error("can not find product", err,
			zap.String("service", "cartService"),
			zap.String("method", "Add"),
			zap.Int("productID", params.ID),
			zap.Int("userID", u.ID),
			zap.Bool("isGuest", u.IsGuest),
		)
		return response.Error("something went wrong", http.StatusInternalServerError, nil)
	}
	if errors.Is(err, platform.ErrProductNotFound) {
		return response.Error(err.Error(), http.StatusUnprocessableEntity, nil)
	}
	cart, err := c.cartStorage.GetCart(u.ID)
	if err != nil {
		c.logger.Error("can not get cart", err,
			zap.String("service", "cartService"),
			zap.String("method", "Add"),
			zap.Int("productID", params.ID),
			zap.Int("userID", u.ID),
			zap.Bool("isGuest", u.IsGuest),
		)
		return response.Error("something went wrong", http.StatusInternalServerError, nil)
	}
	itemFound := false
	for i, item := range cart.Items {
		if item.ProductID == product.ID {
			cart.Items[i].Quantity = params.Quantity
			itemFound = true
			break
		}
	}
	if !itemFound {
		item := CartItem{
			ProductID:    params.ID,
			ProductTitle: product.Title,
			Quantity:     params.Quantity,
			Price:        product.Price,
		}
		cart.Items = append(cart.Items, item)
	}
	c.calculateTotal(&cart)
	err = c.cartStorage.UpdateCart(u.ID, cart)
	if err != nil {
		c.logger.Error("can not update cart", err,
			zap.String("service", "cartService"),
			zap.String("method", "Add"),
			zap.Int("productID", params.ID),
			zap.Int("userID", u.ID),
			zap.Bool("isGuest", u.IsGuest),
		)
		return response.Error("something went wrong", http.StatusInternalServerError, nil)
	}

	res := CartResponse{
		Cart: cart,
	}
	return response.Success(res, "")
}

func (c *cartService) Remove(u *user.AutheticatedUser, params RemoveParams) (apiResponse response.ApiResponse, statusCode int) {
	product, err := c.productService.FindByID(params.ID)
	if err != nil && !errors.Is(err, platform.ErrProductNotFound) {
		c.logger.Error("can not find product", err,
			zap.String("service", "cartService"),
			zap.String("method", "Remove"),
			zap.Int("productID", params.ID),
			zap.Int("userID", u.ID),
			zap.Bool("isGuest", u.IsGuest),
		)
		return response.Error("something went wrong", http.StatusInternalServerError, nil)
	}
	if errors.Is(err, platform.ErrProductNotFound) {
		return response.Error(err.Error(), http.StatusUnprocessableEntity, nil)
	}
	cart, err := c.cartStorage.GetCart(u.ID)
	if err != nil {
		c.logger.Error("can not get cart", err,
			zap.String("service", "cartService"),
			zap.String("method", "Remove"),
			zap.Int("productID", params.ID),
			zap.Int("userID", u.ID),
			zap.Bool("isGuest", u.IsGuest),
		)
		return response.Error("something went wrong", http.StatusInternalServerError, nil)
	}

	itemFound := false
	for i, item := range cart.Items {
		if item.ProductID == product.ID {
			cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			itemFound = true
			break
		}
	}
	if itemFound {
		c.calculateTotal(&cart)
		err = c.cartStorage.UpdateCart(u.ID, cart)
		if err != nil {
			c.logger.Error("can not update cart", err,
				zap.String("service", "cartService"),
				zap.String("method", "Remove"),
				zap.Int("productID", params.ID),
				zap.Int("userID", u.ID),
				zap.Bool("isGuest", u.IsGuest),
			)
			return response.Error("something went wrong", http.StatusInternalServerError, nil)
		}

	}
	res := CartResponse{
		Cart: cart,
	}
	return response.Success(res, "")
}

func (c *cartService) Get(u *user.AutheticatedUser) (apiResponse response.ApiResponse, statusCode int) {
	cart, err := c.cartStorage.GetCart(u.ID)
	if err != nil {
		c.logger.Error("can not get cart", err,
			zap.String("service", "cartService"),
			zap.String("method", "Get"),
			zap.Int("userID", u.ID),
			zap.Bool("isGuest", u.IsGuest),
		)
		return response.Error("something went wrong", http.StatusInternalServerError, nil)
	}
	res := CartResponse{
		Cart: cart,
	}
	return response.Success(res, "")
}

func (c *cartService) calculateTotal(cart *Cart) {
	var total float64
	for _, item := range cart.Items {
		total += float64(item.Quantity) * item.Price

	}
	cart.Total = total
}

func NewCartService(cartStorage CartStorage, productService product.Service, logger platform.Logger) CartService {
	return &cartService{
		cartStorage:    cartStorage,
		productService: productService,
		logger:         logger,
	}
}
