package product

import "coretrix/internal/platform"

var allProducts = []Product{
	{
		ID:          1,
		Title:       "product1",
		Price:       1.99,
		Description: "this is a product 1",
		Image:       "/images/products/1.jpg",
	},
	{
		ID:          2,
		Title:       "product2",
		Price:       2.99,
		Description: "this is a product 2",
		Image:       "/images/products/2.jpg",
	},
	{
		ID:          3,
		Title:       "product3",
		Price:       2.99,
		Description: "this is a product 3",
		Image:       "/images/products/3.jpg",
	},
	{
		ID:          4,
		Title:       "product4",
		Price:       1.99,
		Description: "this is a product 4",
		Image:       "/images/products/4.jpg",
	},
	{
		ID:          5,
		Title:       "product5",
		Price:       3.99,
		Description: "this is a product 5",
		Image:       "/images/products/5.jpg",
	},
}

type Product struct {
	ID          int
	Title       string
	Price       float64
	Description string
	Image       string
}

type Repository interface {
	GetAll() ([]Product, error)
	FindByID(ID int) (Product, error)
}

type repository struct {
}

func (r *repository) FindByID(ID int) (Product, error) {
	for _, p := range allProducts {
		if p.ID == ID {
			return p, nil
		}
	}
	return Product{}, platform.ErrProductNotFound
}

func (r *repository) GetAll() ([]Product, error) {
	return allProducts, nil
}

func NewRepository() Repository {
	return &repository{}
}
