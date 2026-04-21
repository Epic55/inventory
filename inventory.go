package inventory

import "sync"

type Product struct {
	ID    string
	Name  string
	Stock int
}

type SafeInventoryService struct {
	mu       sync.RWMutex
	products map[string]*Product
}

func (s *SafeInventoryService) GetStock(productID string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	product := s.products[productID]
	if product == nil {
		return 0
	}
	return product.Stock
}

func (s *SafeInventoryService) Reserve(productID string, quantity int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	product, ok := s.products[productID]
	if !ok || product == nil {
		return ErrProductNotFound
	}

	if product.Stock < quantity {
		return ErrInsufficientStock
	}

	product.Stock -= quantity
	return nil
}

func (s *SafeInventoryService) ReserveMultiple(items []ReserveItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, item := range items {
		product, ok := s.products[item.ProductID]
		if !ok || product == nil {
			return ErrProductNotFound
		}
		if product.Stock < item.Quantity {
			return ErrInsufficientStock
		}
	}

	for _, item := range items {
		s.products[item.ProductID].Stock -= item.Quantity
	}

	return nil
}
