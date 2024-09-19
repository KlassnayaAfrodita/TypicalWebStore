package storage

import (
	"errors"
	"sync"
)

var NotFound = errors.New("not found")

type Product struct {
	ID       int     `json:"product_id"`
	Name     string  `json:"product_name"`
	Price    float32 `json:"product_price"`
	Quantity int     `json:"product_quantity"`
	About    string  `json:"product_about"`
}

type ProductStore struct {
	products map[int]Product
	mu       *sync.RWMutex
	nextID   int
}

func NewProductStore() *ProductStore {
	return &ProductStore{
		products: map[int]Product{},
		mu:       &sync.RWMutex{},
	}
}

func (ps *ProductStore) AddProduct(product Product) (Product, error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.nextID++
	product.ID = ps.nextID
	ps.products[ps.nextID] = product

	return product, nil
}

func (ps *ProductStore) GetProduct(id int) (Product, error) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	product, ok := ps.products[id]
	if !ok {
		return Product{}, NotFound
	}

	return product, nil
}

func (ps *ProductStore) GetProducts() ([]Product, error) {
	result := make([]Product, 0, len(ps.products))

	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, value := range ps.products {
		result = append(result, value)
	}

	return result, nil
}

func (ps *ProductStore) ChangeProduct(in Product) (Product, error) {
	_, ok := ps.products[in.ID]

	if !ok {
		return Product{}, NotFound
	}

	ps.products[in.ID] = in
	return in, nil
}

func (ps *ProductStore) DeleteProduct(in Product) (Product, error) {
	product, ok := ps.products[in.ID]

	if !ok {
		return Product{}, NotFound
	}

	delete(ps.products, in.ID)
	return product, nil
}
