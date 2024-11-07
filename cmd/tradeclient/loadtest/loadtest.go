package loadtest

import (
	"fmt"
	"sync"
	"time"

	"github.com/quickfixgo/examples/cmd/tradeclient/internal"
)

// LoadTestConfig holds configuration for the load test.
type LoadTestConfig struct {
	OrdersPerSecond int // Rate of orders per second
	TotalOrders     int // Total number of orders to send
	SenderCompID    string
	TargetCompID    string
}

// sends orders based on the provided configuration.
func RunLoadTest(cfg LoadTestConfig) {
	var wg sync.WaitGroup

	// send orders at the specified rate
	for i := 0; i < cfg.TotalOrders; i++ {
		wg.Add(1)
		go func(orderID int) {
			defer wg.Done()
			err := internal.QueryEnterOrder(cfg.SenderCompID, cfg.TargetCompID)
			if err != nil {
				fmt.Printf("Order %d failed: %v\n", orderID, err)
			}
		}(i)

		// Delay to maintain order rate
		time.Sleep(time.Second / time.Duration(cfg.OrdersPerSecond))
	}

	// Wait for all goroutines to complete
	wg.Wait()

	fmt.Println("Load test finished, all orders have been processed.")
}
