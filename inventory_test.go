package inventory

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestReserve_Concurrent(t *testing.T) {
	const (
		initialStock = 100
		goroutines   = 200
		reserveQty   = 1
	)

	stock := &InventoryService{
		products: map[string]*Product{
			"someValue": {ID: "someValue", Name: "stockName", Stock: initialStock},
		},
	}

	var (
		wg           sync.WaitGroup
		successCount atomic.Int64
		failCount    atomic.Int64
	)

	ready := make(chan struct{})

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ready
			err := stock.Reserve("someValue", reserveQty)
			if err == nil {
				successCount.Add(1)
			} else {
				failCount.Add(1)
			}
		}()
	}

	close(ready)
	wg.Wait()

	t.Logf("success=%d fail=%d", successCount.Load(), failCount.Load())

	if got := successCount.Load(); got != initialStock {
		t.Errorf("expected %d successes, got %d", initialStock, got)
	}
	if got := failCount.Load(); got != goroutines-initialStock {
		t.Errorf("expected %d failures, got %d", goroutines-initialStock, got)
	}

	stock.mu.Lock()
	finalStock := stock.products["someValue"].Stock
	stock.mu.Unlock()

	if finalStock != 0 {
		t.Errorf("expected final stock=0, got %d", finalStock)
	}
}

func TestReserveMultiple_Atomicity(t *testing.T) {
	const (
		stockA     = 10
		stockB     = 5
		reserveA   = 8
		reserveB   = 8
		goroutines = 100
	)

	stock := &InventoryService{
		products: map[string]*Product{
			"A": {ID: "A", Name: "Product A", Stock: stockA},
			"B": {ID: "B", Name: "Product B", Stock: stockB},
		},
	}

	items := []Product{
		{ID: "A", Stock: reserveA},
		{ID: "B", Stock: reserveB},
	}

	var (
		wg           sync.WaitGroup
		successCount atomic.Int64
		failCount    atomic.Int64
	)

	ready := make(chan struct{})

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ready
			err := stock.ReserveMultiple(items)
			if err == nil {
				successCount.Add(1)
			} else {
				failCount.Add(1)
			}
		}()
	}

	close(ready)
	wg.Wait()

	t.Logf("success=%d fail=%d", successCount.Load(), failCount.Load())

	if got := successCount.Load(); got != 0 {
		t.Errorf("expected 0 successes, got %d — partial reservation occurred", got)
	}
	if got := failCount.Load(); got != goroutines {
		t.Errorf("expected %d failures, got %d", goroutines, got)
	}

	stock.mu.Lock()
	finalStockA := stock.products["A"].Stock
	finalStockB := stock.products["B"].Stock
	stock.mu.Unlock()

	if finalStockA != stockA {
		t.Errorf("product A: expected stock=%d, got=%d — A was partially deducted", stockA, finalStockA)
	}
	if finalStockB != stockB {
		t.Errorf("product B: expected stock=%d, got=%d — B was modified", stockB, finalStockB)
	}
}
