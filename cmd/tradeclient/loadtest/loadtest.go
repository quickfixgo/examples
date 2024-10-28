package loadtest

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/quickfixgo/examples/cmd/readmetrics"
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
	Status       string        // "success" or "failure"
	ErrorMessage string        // Detailed error message in case of failure
}

// RunLoadTest runs the load test based on the provided configuration.
func RunLoadTest(cfg LoadTestConfig) {
	outputFile, err := os.OpenFile("output.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening output.log: %v", err)
	}
	defer outputFile.Close()

	var wg sync.WaitGroup
	var mu sync.Mutex
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

			localArrival := responseTime
			latency := localArrival.Sub(sentTime)
			status := "success"
			errorMessage := ""

			if err != nil {
				status = "failure"
				errorMessage = err.Error()
				mu.Lock()
				failureCount++
				mu.Unlock()
			} else {
				mu.Lock()
				successCount++
				mu.Unlock()
			}

			mu.Lock()
			timestamps = append(timestamps, OrderTimestamp{
				SentTime:     sentTime,
				ResponseTime: responseTime,
				LocalArrival: localArrival,
				Latency:      latency,
				Status:       status,
				ErrorMessage: errorMessage,
			})
			mu.Unlock()
		}(i)

		time.Sleep(time.Second / time.Duration(cfg.OrdersPerSecond))
	}

	// Wait for all goroutines to finish
	wg.Wait()
	totalTime := time.Since(startTime)

	// Calculate success and failure rates
	successRate := float64(successCount) / float64(cfg.TotalOrders) * 100
	failRate := float64(failureCount) / float64(cfg.TotalOrders) * 100

	// Prepare result summary
	resultSummary := fmt.Sprintf("Sent %d orders in %s\nSuccess Rate: %.2f%%\nFailure Rate: %.2f%%\n",
		cfg.TotalOrders, totalTime, successRate, failRate)
	if _, err := outputFile.WriteString(resultSummary); err != nil {
		log.Fatalf("error writing to output.log: %v", err)
	}

	// Log detailed results in JSON format
	for _, ts := range timestamps {
		tsJson, _ := json.Marshal(ts)
		if _, err := outputFile.WriteString(fmt.Sprintf("%s\n", tsJson)); err != nil {
			log.Fatalf("error writing to output.log: %v", err)
		}
	}

	fmt.Println("Load test complete.")

	// Call readmetrics after the load test
	err = readmetrics.Execute()
	if err != nil {
		log.Fatalf("Error executing readmetrics: %v", err)
	}
}
