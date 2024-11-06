package readmetrics

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LogEntry represents a structure for the relevant log information
type LogEntry struct {
	MessageType string            `json:"message_type"`
	Timestamp   string            `json:"timestamp"`
	Fields      map[string]string `json:"fields"`
}

// Struct to store log entries
type LogMetricsEntry struct {
	timestamp time.Time
	msgType   string
	clOrdID   string
}

// Execute reads the log file, extracts relevant information, and saves it as JSON
func Execute(logFilePath, outputFilePath, tmpDir string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting working directory: %v", err)
	}

	logFile, err := os.Open(filepath.Join(dir, logFilePath))
	if err != nil {
		return fmt.Errorf("error opening log file: %v", err)
	}
	defer logFile.Close()

	scanner := bufio.NewScanner(logFile)
	entries := make([]LogEntry, 0)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "35=D") || strings.Contains(line, "35=8") {
			entry := LogEntry{
				Fields: make(map[string]string),
			}

			parts := strings.Split(line, " ")
			if len(parts) > 2 {
				entry.MessageType = strings.Split(parts[2], "\u0001")[0]
				entry.Timestamp = parts[1]

				// Extract fields
				for _, part := range parts {
					if strings.Contains(part, "=") {
						keyValue := strings.SplitN(part, "=", 2)
						if len(keyValue) == 2 {
							entry.Fields[keyValue[0]] = keyValue[1]
						}
					}
				}
			}

			entries = append(entries, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading log file: %v", err)
	}

	if err := saveToJSON(entries, outputFilePath); err != nil {
		return fmt.Errorf("error saving to JSON: %v", err)
	}

	if err := CalculateLatenciesToFile(logFilePath, tmpDir); err != nil {
		return fmt.Errorf("error calculating latencies: %v", err)
	}

	// Count filled and new orders and calculate success rate
	filledCount, newOrderCount, successRate, err := countFilledOrders(logFilePath)
	if err != nil {
		return fmt.Errorf("error calculating success/failure percentages: %v", err)
	}

	// Write metrics (new orders, filled orders, success rate) to the metrics file
	if err := writeMetricsToFile(tmpDir, filledCount, newOrderCount, successRate); err != nil {
		return fmt.Errorf("error writing metrics to file: %v", err)
	}

	fmt.Printf("Raw Data saved to %s\n", outputFilePath)
	return nil
}

// saveToJSON converts entries to JSON format and saves to a file
func saveToJSON(entries []LogEntry, outputFilePath string) error {
	jsonData, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("error converting to JSON: %v", err)
	}

	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting working directory: %v", err)
	}
	outputFile, err := os.Create(filepath.Join(dir, outputFilePath))
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outputFile.Close()

	_, err = outputFile.Write(jsonData)
	if err != nil {
		return fmt.Errorf("error writing to output file: %v", err)
	}

	return nil
}

// parseFIXMessage parses a FIX message from a log line.
func parseFIXMessage(line string) (LogMetricsEntry, error) {
	fields := strings.Split(line, "")
	msg := LogMetricsEntry{}
	timestampStr := line[:26]
	timestamp, err := time.Parse("2006/01/02 15:04:05.000000", timestampStr)
	if err != nil {
		return msg, err
	}
	msg.timestamp = timestamp

	for _, field := range fields {
		if strings.HasPrefix(field, "35=") {
			msg.msgType = strings.TrimPrefix(field, "35=")
		} else if strings.HasPrefix(field, "11=") {
			msg.clOrdID = strings.TrimPrefix(field, "11=")
		}
	}
	return msg, nil
}

// CalculateLatenciesToFile reads a log file, calculates latencies for 35=D messages,
// and writes the latencies and throughput to separate files in the specified directory.
func CalculateLatenciesToFile(logFilePath, tmpDir string) error {
	file, err := os.Open(logFilePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	dMessages := make(map[string]LogMetricsEntry)
	latencies := []int64{} // Store latencies in an array for average calculation
	throughputCounts := make(map[time.Time]int)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		msg, err := parseFIXMessage(line)
		if err != nil {
			fmt.Println("Error parsing line:", err)
			continue
		}

		// Track 35=D message timestamps for latency and throughput
		if msg.msgType == "D" {
			dMessages[msg.clOrdID] = msg

			// Round down timestamp to the nearest minute for throughput calculation
			minute := msg.timestamp.Truncate(time.Minute)
			throughputCounts[minute]++
		} else if msg.msgType == "8" && msg.clOrdID != "" {
			// Calculate latency
			if dMsg, found := dMessages[msg.clOrdID]; found {
				latency := msg.timestamp.Sub(dMsg.timestamp).Milliseconds()
				latencies = append(latencies, latency)
				delete(dMessages, msg.clOrdID) // Remove to avoid multiple calculations for same ClOrdID
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	// Write latencies to a separate file in the tmpDir directory
	latencyFile, err := os.Create(filepath.Join(tmpDir, "latencies.txt"))
	if err != nil {
		return fmt.Errorf("error creating latencies file: %v", err)
	}
	defer latencyFile.Close()

	writer := bufio.NewWriter(latencyFile)

	// Write latency data
	for index, latency := range latencies {
		_, err := writer.WriteString(fmt.Sprintf("Latency %d: %d ms\n", index+1, latency))
		if err != nil {
			return fmt.Errorf("error writing to latencies file: %v", err)
		}
	}

	// Calculate average latency
	averageLatency := float64(0)
	if len(latencies) > 0 {
		for _, latency := range latencies {
			averageLatency += float64(latency)
		}
		averageLatency /= float64(len(latencies))
	}

	// Write output for average latency and throughput to another file in tmpDir
	metricsFile, err := os.Create(filepath.Join(tmpDir, "metrics.txt"))
	if err != nil {
		return fmt.Errorf("error creating metrics file: %v", err)
	}
	defer metricsFile.Close()

	metricsWriter := bufio.NewWriter(metricsFile)

	// Write the average latency to the metrics file
	_, err = metricsWriter.WriteString(fmt.Sprintf("Average Latency: %.2f ms\n", averageLatency))
	if err != nil {
		return fmt.Errorf("error writing average latency to metrics file: %v", err)
	}

	// Write throughput data
	for minute, count := range throughputCounts {
		throughputStr := fmt.Sprintf("Minute: %s, Throughput: %d orders/min\n", minute.Format("2006-01-02 15:04"), count)
		_, err := metricsWriter.WriteString(throughputStr)
		if err != nil {
			return fmt.Errorf("error writing throughput to metrics file: %v", err)
		}
	}

	writer.Flush()
	metricsWriter.Flush()

	return nil
}

// countFilledOrders reads a FIX log file and counts how many orders were filled (150=F),
// as well as the number of new orders (35=D), and calculates the success rate.
func countFilledOrders(logFilePath string) (int, int, float64, error) {
	file, err := os.Open(logFilePath)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to open log file: %v", err)
	}
	defer file.Close()

	var filledCount, newOrderCount int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the line contains 35=D (new order)
		if strings.Contains(line, "35=D") {
			newOrderCount++
		}

		// Check if the line contains 150=F (filled order)
		if strings.Contains(line, "150=F") {
			filledCount++
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, 0, 0, fmt.Errorf("failed to scan log file: %v", err)
	}

	// Calculate success rate (if newOrderCount > 0 to avoid division by zero)
	var successRate float64
	if newOrderCount > 0 {
		successRate = float64(filledCount) / float64(newOrderCount) * 100
	}

	return filledCount, newOrderCount, successRate, nil
}

// writeMetricsToFile writes the filled and new orders count, and the success rate to the metrics file
func writeMetricsToFile(tmpDir string, filledCount, newOrderCount int, successRate float64) error {
	metricsFile, err := os.OpenFile(filepath.Join(tmpDir, "metrics.txt"), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening metrics file: %v", err)
	}
	defer metricsFile.Close()

	metricsWriter := bufio.NewWriter(metricsFile)

	// Write filled and new orders count, and success rate to the metrics file
	_, err = metricsWriter.WriteString(fmt.Sprintf("Total New Orders: %v\n", newOrderCount))
	if err != nil {
		return fmt.Errorf("error writing new orders count to metrics file: %v", err)
	}

	_, err = metricsWriter.WriteString(fmt.Sprintf("Total Orders Successfully Filled: %v\n", filledCount))
	if err != nil {
		return fmt.Errorf("error writing filled orders count to metrics file: %v", err)
	}

	_, err = metricsWriter.WriteString(fmt.Sprintf("Success Rate: %.2f%%\n", successRate))
	if err != nil {
		return fmt.Errorf("error writing success rate to metrics file: %v", err)
	}

	return metricsWriter.Flush()
}
