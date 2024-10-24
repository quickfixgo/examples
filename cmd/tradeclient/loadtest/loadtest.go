package loadtest

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/quickfixgo/examples/cmd/tradeclient/internal"
)

// LoadTestConfig holds configuration for the load test.
type LoadTestConfig struct {
	OrdersPerSecond int
	TotalOrders     int
	SenderCompID    string
	TargetCompID    string
}

// OrderTimestamp holds the sent, response, and local arrival time of an order.
type OrderTimestamp struct {
	SentTime     time.Time
	ResponseTime time.Time
	LocalArrival time.Time     // Time when the response is received
	Latency      time.Duration // Latency calculated
}

// RunLoadTest runs the load test based on the provided configuration.
func RunLoadTest(cfg LoadTestConfig) {
	outputFile, err := os.OpenFile("output.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening output.log: %v", err)
	}
	defer outputFile.Close()

	var wg sync.WaitGroup
	successCount := 0
	failureCount := 0
	timestamps := make([]OrderTimestamp, 0, cfg.TotalOrders)

	startTime := time.Now()

	// Launch goroutines to send orders
	for i := 0; i < cfg.TotalOrders; i++ {
		wg.Add(1)
		go func(orderID int) {
			defer wg.Done()
			sentTime := time.Now()
			err := internal.QueryEnterOrder(cfg.SenderCompID, cfg.TargetCompID)
			responseTime := time.Now()

			// Simulate local arrival time (could be the same as response time for simplicity)
			localArrival := responseTime

			// Calculate latency
			latency := localArrival.Sub(sentTime)

			timestamps = append(timestamps, OrderTimestamp{
				SentTime:     sentTime,
				ResponseTime: responseTime,
				LocalArrival: localArrival,
				Latency:      latency,
			})

			if err != nil {
				failureCount++
				return
			}
			successCount++
		}(i)

		time.Sleep(time.Second / time.Duration(cfg.OrdersPerSecond))
	}

	// Wait for all goroutines to finish
	wg.Wait()
	totalTime := time.Since(startTime)

	// Log results only after all orders are processed
	successRate := float64(successCount) / float64(cfg.TotalOrders) * 100
	failRate := float64(failureCount) / float64(cfg.TotalOrders) * 100

	// Prepare detailed results for logging
	resultSummary := fmt.Sprintf("Sent %d orders in %s\nSuccess Rate: %.2f%%\nFailure Rate: %.2f%%\n",
		cfg.TotalOrders, totalTime, successRate, failRate)

	// Log the results to output.log
	if _, err := outputFile.WriteString(resultSummary); err != nil {
		log.Fatalf("error writing to output.log: %v", err)
	}

	// Print only a simple message to the console
	fmt.Println("Load test complete.")

	// Analyze the timestamps and log intervals
	for i, ts := range timestamps {
		if i > 0 {
			interval := ts.ResponseTime.Sub(timestamps[i-1].ResponseTime)
			if _, err := outputFile.WriteString(fmt.Sprintf("Order %d - Interval from previous response: %v\n", i+1, interval)); err != nil {
				log.Fatalf("error writing to output.log: %v", err)
			}
		}
		responseTime := ts.ResponseTime.Sub(ts.SentTime)
		if _, err := outputFile.WriteString(fmt.Sprintf("Order %d - Time taken to process: %v\n", i+1, responseTime)); err != nil {
			log.Fatalf("error writing to output.log: %v", err)
		}
		if _, err := outputFile.WriteString(fmt.Sprintf("Order %d - Latency: %v\n", i+1, ts.Latency)); err != nil {
			log.Fatalf("error writing to output.log: %v", err)
		}
	}
}
