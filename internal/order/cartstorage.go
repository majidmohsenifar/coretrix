package order

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type CartItem struct {
	ProductID    int     `json:"productId"`
	ProductTitle string  `json:"productTitle"`
	Quantity     int     `json:"quantity"`
	Price        float64 `json:"price"`
}
type Cart struct {
	Items []CartItem `json:"items"`
	Total float64    `json:"total"`
}

type CartStorage interface {
	GetCart(userID int) (Cart, error)
	UpdateCart(userID int, cart Cart) error
	DeleteCart(userID int) error
}

type cartStorage struct {
	rc *redis.Client
}

func (c *cartStorage) GetCart(userID int) (Cart, error) {
	cart := Cart{}
	key := c.getKey(userID)
	data, err := c.rc.Get(context.Background(), key).Result()
	if err != nil && err != redis.Nil {
		return cart, err
	}

	if err == redis.Nil {
		err := c.updateByKey(key, cart)
		if err != nil {
			return cart, err
		}
		return cart, nil
	}

	//here it means we have no errors
	err = json.Unmarshal([]byte(data), &cart)
	return cart, err
}

func (c *cartStorage) UpdateCart(userID int, cart Cart) error {
	key := c.getKey(userID)
	return c.updateByKey(key, cart)
}

func (c *cartStorage) DeleteCart(userID int) error {
	key := c.getKey(userID)
	err := c.rc.Del(context.Background(), key).Err()
	if err != nil && err != redis.Nil {
		return err
	}
	return nil
}
func (c *cartStorage) updateByKey(key string, cart Cart) error {
	data, err := json.Marshal(cart)
	if err != nil {
		return err
	}
	return c.rc.Set(context.Background(), key, string(data), time.Duration(0)).Err()
}

func (c *cartStorage) getKey(userID int) string {
	return fmt.Sprintf("cart:%d", userID)
}

func NewCartStorage(rc *redis.Client) CartStorage {
	return &cartStorage{
		rc: rc,
	}
}
