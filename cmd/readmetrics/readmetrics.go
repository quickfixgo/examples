package readmetrics

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Struct to hold timestamps
type TestTS struct {
	LTime time.Time
	TTime time.Time
}

const (
	LTimeLayout = "15:04:05.999999"
	TTimeLayout = "15:04:05.999"
	LogFilePath = "tmp/FIX.4.4-CUST2_Order-ANCHORAGE.messages.current.log"
)

// Execute reads and processes metrics from a FIX log file
func Execute() error {
	// Open the log file
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting working directory: %v", err)
	}

	logFile, err := os.Open(filepath.Join(dir, LogFilePath))
	if err != nil {
		return fmt.Errorf("error opening log file: %v", err)
	}
	defer logFile.Close()

	// Open metrics log file
	metricsLogFile := "metrics.log"
	metricsLog, err := os.Create(metricsLogFile)
	if err != nil {
		return fmt.Errorf("error creating metrics log file: %v", err)
	}
	defer metricsLog.Close()

	// Read log lines and parse timestamps
	scanner := bufio.NewScanner(logFile)
	times := make([]TestTS, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "35=D") { // Filter based on message type
			sub1 := strings.Split(line, " ")
			if len(sub1) > 2 {
				localTime := sub1[1]
				parsedLTime, err := time.Parse(LTimeLayout, localTime)
				if err != nil {
					return fmt.Errorf("error parsing local time: %v", err)
				}

				sub2 := strings.Split(sub1[2], "\u0001")
				if len(sub2) > 5 && strings.Contains(sub2[5], "52=") {
					sub3 := strings.Split(sub2[5], "-")
					if len(sub3) > 1 {
						talosTime := sub3[1]
						parsedTTime, err := time.Parse(TTimeLayout, talosTime)
						if err != nil {
							return fmt.Errorf("error parsing talos time: %v", err)
						}

						times = append(times, TestTS{LTime: parsedLTime, TTime: parsedTTime})
					}
				}
			}
		}
	}

	// Variables for calculating averages
	var totalInterval float64
	var totalLatency float64

	// Calculate latency and throughput, then write to metrics log
	for i := 1; i < len(times); i++ {
		t := times[i]
		t0 := times[i-1]

		latency := float64(t.LTime.Sub(t.TTime)) / float64(time.Millisecond)
		interval := float64(t.TTime.Sub(t0.TTime)) / float64(time.Millisecond)

		totalLatency += latency
		totalInterval += interval

		// Logging the individual message metrics
		fmt.Fprintf(metricsLog, "Message %d: Interval = %.3fms, Latency = %.3fms\n", i, interval, latency)
	}

	// Calculate averages
	avgLatency := totalLatency / float64(len(times)-1)
	avgInterval := totalInterval / float64(len(times)-1)

	// Throughput Measurement
	totalMessages := len(times) - 1 // Total messages should exclude the first (which has no interval)
	totalDuration := times[len(times)-1].TTime.Sub(times[0].TTime).Seconds()
	throughput := float64(totalMessages) / totalDuration

	// Write averages and throughput to metrics log
	fmt.Fprintf(metricsLog, "Average Interval = %.3fms, Average Latency = %.3fms\n", avgInterval, avgLatency)
	fmt.Fprintf(metricsLog, "Throughput: %.2f messages per second\n", throughput)

	// Print overall metrics to console
	fmt.Printf("Average Interval = %.3fms, Average Latency = %.3fms\n", avgInterval, avgLatency)
	fmt.Printf("Throughput: %.2f messages per second\n", throughput)
	fmt.Println("Latency and throughput metrics logged to metrics.log")
	return nil
}
