package readmetrics

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type TestTS struct {
	LTime time.Time
	TTime time.Time
}

const (
	LTimeLayout   = "15:04:05.999999"
	TTimeLayout   = "15:04:05.999"
	LogFilePath   = "tmp/FIX.4.4-CUST2_Order-ANCHORAGE.messages.current.log"
	MetricsFolder = "loadtest_metrics"
)

// Execute reads and processes metrics from a FIX log file
func Execute() error {
	// Create metrics folder if it doesn't exist
	if _, err := os.Stat(MetricsFolder); os.IsNotExist(err) {
		err = os.Mkdir(MetricsFolder, 0755)
		if err != nil {
			return fmt.Errorf("error creating metrics folder: %v", err)
		}
	}

	// Define metrics log file path with a timestamp for uniqueness
	timestamp := time.Now().Format("20060102_150405")
	metricsLogFile := filepath.Join(MetricsFolder, fmt.Sprintf("metrics_%s.log", timestamp))
	metricsLog, err := os.Create(metricsLogFile)
	if err != nil {
		return fmt.Errorf("error creating metrics log file: %v", err)
	}
	defer metricsLog.Close()

	// Open the log file to read
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting working directory: %v", err)
	}

	logFile, err := os.Open(filepath.Join(dir, LogFilePath))
	if err != nil {
		return fmt.Errorf("error opening log file: %v", err)
	}
	defer logFile.Close()

	// Initialize a scanner to read the file line by line
	scanner := bufio.NewScanner(logFile)
	times := make([]TestTS, 0)

	// Process each line in the log file
	for scanner.Scan() {
		line := scanner.Text()

		// Check for specific message type "35=D" to filter relevant messages
		if strings.Contains(line, "35=D") {
			sub1 := strings.Split(line, " ")
			if len(sub1) > 2 {
				// Extract and parse the local timestamp (LTime)
				localTime := sub1[1]
				parsedLTime, err := time.Parse(LTimeLayout, localTime)
				if err != nil {
					return fmt.Errorf("error parsing local time: %v", err)
				}

				// Process message body to isolate field-value pairs
				sub2 := strings.Split(sub1[2], "\u0001")
				if len(sub2) > 5 && strings.Contains(sub2[5], "52=") {
					// Extract and parse the timestamp (TTime) from "52="
					sub3 := strings.Split(sub2[5], "-")
					if len(sub3) > 1 {
						orderTimestamp := sub3[1]
						parsedTTime, err := time.Parse(TTimeLayout, orderTimestamp)
						if err != nil {
							return fmt.Errorf("error parsing time: %v", err)
						}

						// Store parsed timestamps in TestTS struct
						times = append(times, TestTS{LTime: parsedLTime, TTime: parsedTTime})
					}
				}
			}
		}
	}

	// Error check for scanner
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading log file: %v", err)
	}

	// Calculate metrics first
	totalMessages := len(times)
	if totalMessages == 0 {
		return fmt.Errorf("no order messages found")
	}

	var totalInterval float64
	var totalLatency float64
	latencies := make([]float64, 0, totalMessages-1)

	// Calculate individual message intervals and latencies
	for i := 1; i < len(times); i++ {
		t := times[i]
		t0 := times[i-1]

		latency := float64(t.LTime.Sub(t.TTime)) / float64(time.Millisecond)
		interval := float64(t.TTime.Sub(t0.TTime)) / float64(time.Millisecond)

		totalLatency += latency
		totalInterval += interval
		latencies = append(latencies, latency) // Store latencies for later calculations
	}

	// Calculate final metrics
	totalDuration := times[len(times)-1].TTime.Sub(times[0].TTime).Seconds()
	throughput := float64(totalMessages) / totalDuration
	avgLatency := totalLatency / float64(len(latencies))
	avgInterval := totalInterval / float64(len(latencies))

	// Calculate min and max latencies
	minLatency := math.MaxFloat64
	maxLatency := 0.0
	for _, latency := range latencies {
		if latency < minLatency {
			minLatency = latency
		}
		if latency > maxLatency {
			maxLatency = latency
		}
	}

	// Print metrics
	fmt.Fprintf(metricsLog, "Throughput: %.2f messages per second\n", throughput)
	fmt.Fprintf(metricsLog, "Average Interval = %.3fms, Average Latency = %.3fms\n", avgInterval, avgLatency)
	fmt.Fprintf(metricsLog, "Message Metrics: Min Latency = %.3fms, Max Latency = %.3fms\n",
		minLatency, maxLatency)

	fmt.Printf("Metrics logged to %s\n", metricsLogFile)
	return nil
}
