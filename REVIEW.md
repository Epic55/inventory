## Part 1

type InventoryService struct {
## ADDED LINE. 
## READ/WRITE OPERATIONS IN MAP ARE NOT SAFE FOR CONCURRENCY, BECAUSE OF THAT WE NEED TO ADD MUTEX
	mu       sync.RWMutex 
	products map[string]*Product
}

func (s *InventoryService) GetStock(productID string) int {
## Race Condition 1. ADDED 2 LINES BELOW:
	s.mu.RLock()
	defer s.mu.RUnlock()

	product := s.products[productID]
	if product == nil {
		return 0
	}
	return product.Stock
}

func (s *InventoryService) Reserve(productID string, quantity int) error {
## Race Condition 2. ADDED 2 LINES BELOW:
	s.mu.Lock()
	defer s.mu.Unlock()

## ADDED CHECKING FOR EXISTING KEY IN MAP
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

func (s *InventoryService) ReserveMultiple(items []ReserveItem) error {
## Race Condition 3. ADDED 2 LINES BELOW:
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check all first
	for _, item := range items {
## ADDED CHECKING FOR NIL POINTER AND FOR EXISTING KEY IN MAP
		product, ok := s.products[item.ProductID]
        if !ok || product == nil {
            return ErrProductNotFound
        }
		if product.Stock < item.Quantity {
			return ErrInsufficientStock
		}
	}

	// Then reserve all
	for _, item := range items {
		s.products[item.ProductID].Stock -= item.Quantity
	}

	return nil
}

func (s *InventoryService) SafeReserve(productID string, quantity int) error {
## Race Condition 4. REPLACED 3 OLD LINES WITH THESE 2 NEW LINES BELOW:
	s.mu.Lock() 
	defer s.mu.Unlock()

## ADDED CHECKING FOR NIL POINTER AND FOR EXISTING KEY IN MAP
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
