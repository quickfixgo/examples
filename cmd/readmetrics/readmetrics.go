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

	// Calculate success and failure percentages and write to metrics file
	filledPct, rejectedPct, err := calculateSuccessFailure(logFilePath)
	if err != nil {
		return fmt.Errorf("error calculating success/failure percentages: %v", err)
	}

	if err := writeMetricsToFile(tmpDir, filledPct, rejectedPct); err != nil {
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

// calculateSuccessFailure reads a FIX log file and calculates the success (filled) and failure (rejected) percentages
func calculateSuccessFailure(logFilePath string) (float64, float64, error) {
	file, err := os.Open(logFilePath)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to open log file: %v", err)
	}
	defer file.Close()

	var filledCount, rejectedCount, totalCount int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the message type is an Execution Report (35=8)
		if strings.Contains(line, "35=8") {
			totalCount++

			// Check for filled (150=F) or rejected (150=8) execution status
			if strings.Contains(line, "150=F") {
				filledCount++
			} else if strings.Contains(line, "150=8") {
				rejectedCount++
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, 0, fmt.Errorf("failed to scan log file: %v", err)
	}

	if totalCount == 0 {
		return 0, 0, fmt.Errorf("no execution reports found in log")
	}

	filledPercentage := (float64(filledCount) / float64(totalCount)) * 100
	rejectedPercentage := (float64(rejectedCount) / float64(totalCount)) * 100

	return filledPercentage, rejectedPercentage, nil
}

// writeMetricsToFile writes the filled and rejected percentages to the metrics file
func writeMetricsToFile(tmpDir string, filledPct, rejectedPct float64) error {
	metricsFile, err := os.OpenFile(filepath.Join(tmpDir, "metrics.txt"), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening metrics file: %v", err)
	}
	defer metricsFile.Close()

	metricsWriter := bufio.NewWriter(metricsFile)

	// Write filled and rejected percentages to the metrics file
	_, err = metricsWriter.WriteString(fmt.Sprintf("Filled Percentage: %.2f%%\n", filledPct))
	if err != nil {
		return fmt.Errorf("error writing filled percentage to metrics file: %v", err)
	}

	_, err = metricsWriter.WriteString(fmt.Sprintf("Rejected Percentage: %.2f%%\n", rejectedPct))
	if err != nil {
		return fmt.Errorf("error writing rejected percentage to metrics file: %v", err)
	}

	return metricsWriter.Flush()
}
