package src

import (
	"context"
	"os"
	"testing"
	"time"
)

// Helper function to create a temporary CSV file for testing
func createTempCSV(content string) (string, error) {
	file, err := os.CreateTemp("", "test_*.csv")
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}

// Test reading a valid CSV file
func TestReadCSVFile_Valid(t *testing.T) {
	filePath := "../testdata/valid.csv"
	urlChan := make(chan string, 50)
	metrics := &Metrics{}
	ctx := context.Background()
	go readCSVFile(filePath, urlChan, metrics, ctx)
	close(urlChan)
	var actualURLs []string
	for url := range urlChan {
		actualURLs = append(actualURLs, url)
	}
	expectedURLs := []string{"https://example.com", "https://google.com"}
	if len(actualURLs) != len(expectedURLs) {
		t.Errorf("Expected %d URLs, got %d", len(expectedURLs), len(actualURLs))
	}
	if metrics.TotalURLs.Load() != uint64(len(expectedURLs)) {
		t.Errorf("Expected TotalURLs=%d, got %d", len(expectedURLs), metrics.TotalURLs.Load())
	}
}

// Test reading an empty CSV file
func TestReadCSVFile_EmptyFile(t *testing.T) {
	filePath := "../testdata/empty.csv"
	urlChan := make(chan string, 50)
	metrics := &Metrics{}
	ctx := context.Background()
	go readCSVFile(filePath, urlChan, metrics, ctx)
	close(urlChan)
	var actualURLs []string
	for url := range urlChan {
		actualURLs = append(actualURLs, url)
	}
	// Verify results
	if len(actualURLs) != 0 {
		t.Errorf("Expected 0 URLs, got %d", len(actualURLs))
	}
	if metrics.TotalURLs.Load() != 0 {
		t.Errorf("Expected TotalURLs=0, got %d", metrics.TotalURLs.Load())
	}

}

// Test handling an invalid format (extra columns)
func TestReadCSVFile_InvalidFormat(t *testing.T) {
	filePath := "../testdata/invalid.csv"
	urlChan := make(chan string, 50)
	metrics := &Metrics{}
	ctx := context.Background()

	go readCSVFile(filePath, urlChan, metrics, ctx)
	close(urlChan)
	var actualURLs []string
	for url := range urlChan {
		actualURLs = append(actualURLs, url)
	}
	// Verify results
	if len(actualURLs) != 0 {
		t.Errorf("Expected 0 URLs, got %d", len(actualURLs))
	}
	if metrics.TotalURLs.Load() != 0 {
		t.Errorf("Expected TotalURLs=0, got %d", metrics.TotalURLs.Load())
	}

}

// Test reading with context cancellation
func TestReadCSVFile_ContextCancelled(t *testing.T) {
	filePath := "../testdata/valid.csv"
	urlChan := make(chan string, 50)
	metrics := &Metrics{}
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(10 * time.Millisecond) // Simulate early cancel
		cancel()
	}()

	go readCSVFile(filePath, urlChan, metrics, ctx)
	time.Sleep(50 * time.Millisecond) // Give some time for cancellation

	select {
	case <-urlChan:
		t.Errorf("Channel should be empty or partially filled due to cancellation")
	default:
	}
}

// Test reading a non-existent file
func TestReadCSVFile_NonExistentFile(t *testing.T) {
	urlChan := make(chan string, 10)
	metrics := &Metrics{}
	ctx := context.Background()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Function should panic when file does not exist")
		}
	}()

	readCSVFile("non_existent_file.csv", urlChan, metrics, ctx)
}
