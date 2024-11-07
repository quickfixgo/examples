package readmetrics

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock function to simulate reading from a fixed logfile and generating expected latencies
func TestCalculateLatenciesToFile(t *testing.T) {
	// Prepare paths for log and temp directory
	logFilePath := "test_data/test.log"
	tmpDir := "test_data"

	// Ensure the tmpDir exists or create it for this test
	err := os.MkdirAll(tmpDir, os.ModePerm)
	require.NoError(t, err)

	// Call the function under test
	err = CalculateLatenciesToFile(logFilePath, tmpDir)
	require.NoError(t, err)

	// Read the generated latencies.txt file
	latenciesFilePath := filepath.Join(tmpDir, "latencies.txt")
	latenciesData, err := os.ReadFile(latenciesFilePath)
	require.NoError(t, err)

	// Define the expected latencies based on the mock data
	// Latency 1: Time difference between first D message and first 8 message
	// Latency 2: Time difference between second D message and second 8 message
	expectedLatencies := []string{
		"Latency 1: 151 ms\n",
		"Latency 2: 156 ms\n",
	}

	// Assert that the latencies file contains the expected latencies
	for i, expected := range expectedLatencies {
		assert.Contains(t, string(latenciesData), expected, "Missing or incorrect latency at index %d", i+1)
	}

	// Read the generated metrics.txt file to ensure it contains the expected average latency and throughput
	metricsFilePath := filepath.Join(tmpDir, "metrics.txt")
	metricsData, err := os.ReadFile(metricsFilePath)
	require.NoError(t, err)

	// Verify that the average latency is present and matches the expected value
	assert.Contains(t, string(metricsData), "Average Latency: 153.50 ms\n", "Average latency is incorrect")

	// Verify that the throughput data is correct
	assert.Contains(t, string(metricsData), "Minute: 2024-11-06 13:48, Throughput: 2 orders/min", "Throughput data is incorrect")
}
