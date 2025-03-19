package src

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type downloadResult struct {
	url     string
	content []byte
}

// downloadURLs concurrently downloads content from URLs received via a channel.
//
// Input:
// - urlChan: A channel that provides URLs for downloading.
// - contentChan: A channel to send the downloaded content for persistence.
// - metrics: A pointer to the Metrics struct for tracking success and failure counts.
// - ctx: Context for graceful shutdown and cancellation handling.
// - wg: WaitGroup to synchronize goroutines.
//
// Output:
// - Downloads content from URLs and sends results to contentChan.
// - Updates metrics for successful and failed downloads.
// - Ensures a maximum of MAX_WORKERS concurrent downloads.
//
// Notes:
// - Uses a semaphore (channel) to limit concurrent downloads.
// - Supports graceful shutdown by listening to ctx.Done().
// - Ensures goroutine cleanup with wg.Done().
func downloadURLs(urlChan <-chan string, contentChan chan<- downloadResult, metrics *Metrics, ctx context.Context, wg *sync.WaitGroup) {
	semaphore := make(chan struct{}, MAX_WORKERS) // Limit to MAX_WORKERS concurrent downloads

	// Process each URL received from the urlChan
	for url := range urlChan {
		select {
		case semaphore <- struct{}{}: // Acquire a semaphore slot
			wg.Add(1)
			go func(u string) {
				defer wg.Done()                // Ensure the goroutine signals completion
				defer func() { <-semaphore }() // Release semaphore slot

				start := time.Now() // Record start time for metrics
				content, err := downloadURL(ctx, ensureScheme(u))
				if err != nil {
					zlog.Error().Msgf("Error downloading %s: %v", u, err)
					metrics.AddFailure() // Track failed downloads
					return
				}

				metrics.AddSuccess(time.Since(start)) // Track successful download duration

				// Send the downloaded content to contentChan or handle shutdown
				select {
				case contentChan <- downloadResult{url: u, content: content}:
				case <-ctx.Done():
					zlog.Info().Msgf("Stage 2: Context canceled / Shutdown initiated. Skipping content persistence.")
					return
				}
			}(url)

		case <-ctx.Done(): // Handle shutdown scenario
			zlog.Error().Msgf("Stage 2: Context canceled / Shutdown initiated. Stopping new downloads.")
			return
		}
	}
}

// downloadURL fetches the content of a given URL using an HTTP GET request.
//
// Input:
// - ctx: Context for handling timeouts or cancellations.
// - url: The URL to download.
//
// Output:
// - Returns the response body as a byte slice ([]byte).
// - Returns an error if the request fails or the response status is not 200 OK.
//
// Notes:
// - Uses http.NewRequestWithContext to support graceful shutdown.
// - Ensures the response body is closed properly to prevent resource leaks.
func downloadURL(ctx context.Context, url string) ([]byte, error) {
	// Create a new HTTP GET request with context for cancellation support
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err // Return error if request creation fails
	}

	// Send the HTTP request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err // Return error if request execution fails
	}
	defer resp.Body.Close() // Ensure the response body is closed

	// Check for non-200 HTTP status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	// Read and return the response body
	return io.ReadAll(resp.Body)
}
