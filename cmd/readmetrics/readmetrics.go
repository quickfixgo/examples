package readmetrics

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Constants for file paths
const LogFilePath = "tmp/FIX.4.4-CUST2_Order-ANCHORAGE.messages.current.log"
const OutputFilePath = "tmp/log_data.json"

// LogEntry represents a structure for the relevant log information
type LogEntry struct {
	MessageType string            `json:"message_type"`
	Timestamp   string            `json:"timestamp"`
	Fields      map[string]string `json:"fields"`
}

// Execute reads the log file, extracts relevant information, and saves it as JSON
func Execute() error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting working directory: %v", err)
	}

	logFile, err := os.Open(filepath.Join(dir, LogFilePath))
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

	if err := saveToJSON(entries); err != nil {
		return fmt.Errorf("error saving to JSON: %v", err)
	}

	fmt.Printf("Data saved to %s\n", OutputFilePath)
	return nil
}

// saveToJSON converts entries to JSON format and saves to a file
func saveToJSON(entries []LogEntry) error {
	jsonData, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("error converting to JSON: %v", err)
	}

	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting working directory: %v", err)
	}
	outputFile, err := os.Create(filepath.Join(dir, OutputFilePath))
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
