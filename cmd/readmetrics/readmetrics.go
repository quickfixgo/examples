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

// LogEntry represents the structure of a log entry parsed from the raw log.
type LogEntry struct {
	MessageType string            `json:"message_type"` // Type of message (e.g., "D", "8").
	Timestamp   string            `json:"timestamp"`    // Timestamp of the log entry.
	Fields      map[string]string `json:"fields"`       // Additional fields in the log.
}

// LogMetricsEntry stores parsed information for latency and throughput calculations.
type LogMetricsEntry struct {
	timestamp time.Time // Timestamp of the message.
	msgType   string    // Type of message (e.g., "D", "8").
	clOrdID   string    // Client Order ID.
}

// Execute processes the log file, calculates metrics, and saves them to output files.
func Execute(logFilePath, outputFilePath, tmpDir string) error {
	// Get the current working directory
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting working directory: %v", err)
	}

	// Open the log file
	logFile, err := os.Open(filepath.Join(dir, logFilePath))
	if err != nil {
		return fmt.Errorf("error opening log file: %v", err)
	}
	defer logFile.Close()

	// Prepare a scanner to read the log file line by line
	scanner := bufio.NewScanner(logFile)
	entries := make([]LogEntry, 0)

	// Read each line in the log file and parse relevant entries
	for scanner.Scan() {
		line := scanner.Text()

		// Filter lines that are message type "D" or "8"
		if strings.Contains(line, "35=D") || strings.Contains(line, "35=8") {
			entry := LogEntry{
				Fields: make(map[string]string),
			}

			// Split the line by spaces and process the parts
			parts := strings.Split(line, " ")
			if len(parts) > 2 {
				// Extract message type and timestamp
				entry.MessageType = strings.Split(parts[2], "\u0001")[0]
				entry.Timestamp = parts[1]

				// Extract fields (key-value pairs) from the log
				for _, part := range parts {
					if strings.Contains(part, "=") {
						keyValue := strings.SplitN(part, "=", 2)
						if len(keyValue) == 2 {
							entry.Fields[keyValue[0]] = keyValue[1]
						}
					}
				}
			}

			entries = append(entries, entry) // Add the entry to the list
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading log file: %v", err)
	}

	// Save the parsed entries to a JSON file
	if err := saveToJSON(entries, outputFilePath); err != nil {
		return fmt.Errorf("error saving to JSON: %v", err)
	}

	// Calculate latencies and save them
	if err := CalculateLatenciesToFile(logFilePath, tmpDir); err != nil {
		return fmt.Errorf("error calculating latencies: %v", err)
	}

	// Calculate success rates for orders
	filledCount, newOrderCount, successRate, err := countFilledOrders(logFilePath)
	if err != nil {
		return fmt.Errorf("error calculating success percentages: %v", err)
	}

	// Write the calculated metrics to file
	if err := writeMetricsToFile(tmpDir, filledCount, newOrderCount, successRate); err != nil {
		return fmt.Errorf("error writing metrics to file: %v", err)
	}

	fmt.Printf("Raw Data saved to %s\n", outputFilePath)
	return nil
}

// saveToJSON saves the parsed log entries to a JSON file.
func saveToJSON(entries []LogEntry, outputFilePath string) error {
	// Marshal entries into JSON format
	jsonData, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("error converting to JSON: %v", err)
	}

	// Get the current working directory to create the output file path
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting working directory: %v", err)
	}

	// Create and open the output file
	outputFile, err := os.Create(filepath.Join(dir, outputFilePath))
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outputFile.Close()

	// Write the JSON data to the file
	_, err = outputFile.Write(jsonData)
	if err != nil {
		return fmt.Errorf("error writing to output file: %v", err)
	}

	return nil
}

// parseFIXMessage parses a single FIX message and returns the relevant data.
func parseFIXMessage(line string) (LogMetricsEntry, error) {
	// Split the line by the FIX field delimiter
	fields := strings.Split(line, "")
	msg := LogMetricsEntry{}

	// Parse the timestamp from the first 26 characters of the line
	timestampStr := line[:26]
	timestamp, err := time.Parse("2006/01/02 15:04:05.000000", timestampStr)
	if err != nil {
		return msg, err
	}
	msg.timestamp = timestamp

	// Extract message type and client order ID
	for _, field := range fields {
		if strings.HasPrefix(field, "35=") {
			msg.msgType = strings.TrimPrefix(field, "35=")
		} else if strings.HasPrefix(field, "11=") {
			msg.clOrdID = strings.TrimPrefix(field, "11=")
		}
	}
	return msg, nil
}

// CalculateLatenciesToFile calculates the latencies between orders and saves the results to files.
func CalculateLatenciesToFile(logFilePath, tmpDir string) error {
	// Open the log file for reading
	file, err := os.Open(logFilePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Initialize variables for storing latency data
	dMessages := make(map[string]LogMetricsEntry)
	latencies := []int64{}
	throughputCounts := make(map[time.Time]int)

	// Read each line of the log and calculate latencies
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		msg, err := parseFIXMessage(line)
		if err != nil {
			fmt.Println("Error parsing line:", err)
			continue
		}

		// Track order creation ("D") messages and calculate latency for execution ("8") messages
		if msg.msgType == "D" {
			dMessages[msg.clOrdID] = msg
			minute := msg.timestamp.Truncate(time.Minute)
			throughputCounts[minute]++
		} else if msg.msgType == "8" && msg.clOrdID != "" {
			if dMsg, found := dMessages[msg.clOrdID]; found {
				latency := msg.timestamp.Sub(dMsg.timestamp).Milliseconds()
				latencies = append(latencies, latency)
				delete(dMessages, msg.clOrdID)
			}
		}
	}

	// Handle any errors encountered during scanning
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	// Save latencies to a file
	latencyFile, err := os.Create(filepath.Join(tmpDir, "latencies.txt"))
	if err != nil {
		return fmt.Errorf("error creating latencies file: %v", err)
	}
	defer latencyFile.Close()

	writer := bufio.NewWriter(latencyFile)
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

	// Save the metrics (average latency, throughput) to file
	metricsFile, err := os.Create(filepath.Join(tmpDir, "metrics.txt"))
	if err != nil {
		return fmt.Errorf("error creating metrics file: %v", err)
	}
	defer metricsFile.Close()

	metricsWriter := bufio.NewWriter(metricsFile)
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

// countFilledOrders counts the number of filled and new orders and calculates the success rate.
func countFilledOrders(logFilePath string) (int, int, float64, error) {
	// Open the log file for scanning
	file, err := os.Open(logFilePath)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to open log file: %v", err)
	}
	defer file.Close()

	var filledCount, newOrderCount int

	// Scan the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Count new orders (type "D")
		if strings.Contains(line, "35=D") {
			newOrderCount++
		}

		// Count filled orders (150=F)
		if strings.Contains(line, "150=F") {
			filledCount++
		}
	}

	// Handle any errors encountered during scanning
	if err := scanner.Err(); err != nil {
		return 0, 0, 0, fmt.Errorf("failed to scan log file: %v", err)
	}

	// Calculate the success rate
	var successRate float64
	if newOrderCount > 0 {
		successRate = float64(filledCount) / float64(newOrderCount) * 100
	}

	return filledCount, newOrderCount, successRate, nil
}

// writeMetricsToFile appends metrics data (new orders, filled orders, success rate) to a file.
func writeMetricsToFile(tmpDir string, filledCount, newOrderCount int, successRate float64) error {
	// Open the metrics file for appending
	metricsFile, err := os.OpenFile(filepath.Join(tmpDir, "metrics.txt"), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening metrics file: %v", err)
	}
	defer metricsFile.Close()

	metricsWriter := bufio.NewWriter(metricsFile)

	// Write the metrics data (new orders, filled orders, success rate) to the file
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

	// Flush the buffered writer to ensure data is written to file
	return metricsWriter.Flush()
}
