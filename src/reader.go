package src

import (
	"bufio"
	"context"
	"encoding/csv"
	"io"
	"log"
	"os"
)

// readCSVFile reads URLs from a CSV file and sends them to a channel for processing.
// It supports graceful shutdown using a context and updates metrics for processed URLs.

// readCSVFile reads URLs from a CSV file and sends them to a channel for processing.
//
// Input:
// - filePath: Path to the CSV file containing URLs (one per line).
// - urlChannel: A channel to send valid URLs for further processing.
// - metrics: A pointer to the Metrics struct to track total URLs processed.
// - ctx: Context for graceful shutdown.
//
// Expected CSV Format:
// - First row is treated as a header and skipped.
// - Each subsequent row contains a single URL.
//
// Output:
// - Sends valid URLs to the urlChannel.
// - Updates the metrics.TotalURLs count.
// - Stops processing when the context is canceled.
//
// Notes:
// - Logs errors for invalid rows but continues processing.
// - Uses a buffered reader for efficient file reading.
func readCSVFile(filePath string, urlChannel chan<- string, metrics *Metrics, ctx context.Context) {
	// Open the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err) // Fatal log will stop execution if the file can't be opened
	}
	defer file.Close() // Ensure the file is closed when function exits

	// Create a CSV reader with buffered input
	reader := csv.NewReader(bufio.NewReader(file))
	reader.FieldsPerRecord = 1 // Enforce only one column per line

	// Skip the header row
	_, err = reader.Read()
	if err != nil {
		zlog.Error().Msgf("Failed to read CSV Header: %v", err) // Log error if header read fails
		return
	}

	// Process each row in the CSV file
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break // Stop reading when reaching end of file
		}
		if err != nil {
			zlog.Error().Msgf("Skipping invalid row: %v", err) // Log and skip malformed rows
			continue
		}
		if len(record) != 1 {
			zlog.Error().Msgf("invalid row colunt : %v", len(record)) // Log and skip malformed rows
			continue                                                  // Skip rows that don't have exactly one field
		}

		metrics.TotalURLs.Add(1) // Update the metrics count

		// Send URL to channel or exit if context is canceled
		select {
		case urlChan <- record[0]: // Send URL to channel
		case <-ctx.Done(): // Handle shutdown scenario
			zlog.Error().Msgf("Stage 1: Context canceled./Shutdown initiated. Stopping file read")
			return
		}
	}
}
