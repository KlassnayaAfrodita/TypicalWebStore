package storage

import (
	"errors"
	"sync"
)

var NotFound = errors.New("not found")

type Comment struct {
	ID     int    `json:"comment_id"`
	Rating int    `json:"camment_rating"`
	Review string `json:"comment_review"`
}

type Product struct {
	ID       int        `json:"product_id"`
	Name     string     `json:"product_name"`
	Price    float32    `json:"product_price"`
	Quantity int        `json:"product_quantity"`
	About    string     `json:"product_about"`
	Comments []*Comment `json:"product_comments`
}

type ProductStorage struct {
	products map[int]Product
	mu       *sync.RWMutex
	nextID   int
}

func NewProductStorage() *ProductStorage {
	return &ProductStorage{
		products: map[int]Product{
			1: Product{
				ID:       1,
				Name:     "laptop",
				Price:    1000.,
				Quantity: 5,
				About:    "laptop",
			},
		},
		mu: &sync.RWMutex{},
	}
}

func (ps *ProductStorage) AddProduct(product Product) (Product, error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.nextID++
	product.ID = ps.nextID
	ps.products[ps.nextID] = product

	return product, nil
}

func (ps *ProductStorage) GetProduct(id int) (Product, error) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	product, ok := ps.products[id]
	if !ok {
		return Product{}, NotFound
	}

	return product, nil
}

func (ps *ProductStorage) GetProducts() ([]Product, error) {
	result := make([]Product, 0, len(ps.products))

	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, value := range ps.products {
		result = append(result, value)
	}

	return result, nil
}

func (ps *ProductStorage) ChangeProduct(in Product) (Product, error) {
	_, ok := ps.products[in.ID]

	if !ok {
		return Product{}, NotFound
	}

	ps.products[in.ID] = in
	return in, nil
}

func (ps *ProductStorage) DeleteProduct(in Product) (Product, error) {
	product, ok := ps.products[in.ID]

	if !ok {
		return Product{}, NotFound
	}

	delete(ps.products, in.ID)
	return product, nil
}
