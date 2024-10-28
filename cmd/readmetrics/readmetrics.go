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
)

// Execute reads metrics from a specified FIX log file.
func Execute(testFileArg string) error {
	// Create the metrics log file
	metricsLogFile := "metrics.log"
	metricsLog, err := os.Create(metricsLogFile)
	if err != nil {
		return fmt.Errorf("error creating metrics log file: %v", err)
	}
	defer metricsLog.Close()

	// Open the input log file
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting working directory: %v", err)
	}

	testFile, err := os.Open(filepath.Join(dir, testFileArg))
	if err != nil {
		return fmt.Errorf("error opening %v: %v", testFileArg, err)
	}
	defer testFile.Close()

	// Create a scanner to read the file
	scanner := bufio.NewScanner(testFile)

	// Read the file line by line and process the timestamps
	times := make([]TestTS, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "35=D") { // Modify this condition based on your message type
			sub1 := strings.Split(line, " ")
			if len(sub1) > 2 {
				localTime := sub1[1]
				parsedLTime, err := time.Parse(LTimeLayout, localTime)
				if err != nil {
					return fmt.Errorf("error parsing local time: %v", err)
				}

				sub2 := strings.Split(sub1[2], "\u0001")
				if len(sub2) > 5 {
					if strings.Contains(sub2[5], "52=") {
						sub3 := strings.Split(sub2[5], "-")
						if len(sub3) > 1 {
							talosTime := sub3[1]
							parsedTTime, err := time.Parse(TTimeLayout, talosTime)
							if err != nil {
								return fmt.Errorf("error parsing talos time: %v", err)
							}

							testTS := TestTS{
								LTime: parsedLTime,
								TTime: parsedTTime,
							}
							times = append(times, testTS)
						}
					}
				}
			}
		}
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	// Calculate latencies and write to the metrics log
	for i := 1; i < len(times); i++ {
		t := times[i]
		t0 := times[i-1]

		// Convert latency and interval to milliseconds
		latency := float64(t.LTime.Sub(t.TTime)) / float64(time.Millisecond)
		interval := float64(t.TTime.Sub(t0.TTime)) / float64(time.Millisecond)
		fmt.Fprintf(metricsLog, "Message %d: Interval = %.3fms, Latency = %.3fms\n", i, interval, latency)
	}

	fmt.Println("Latency metrics logged to metrics.log")
	return nil
}
