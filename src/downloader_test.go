package src


import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

// Mock HTTP server to simulate downloads
func mockHTTPServer(response string, status int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Write([]byte(response))
	}))
}

// Test successful download
func TestDownloadURL_Success(t *testing.T) {
	server := mockHTTPServer("test content", http.StatusOK)
	defer server.Close()
	ctx := context.Background()
	data, err := downloadURL(ctx, server.URL)
	if err != nil {
		t.Fatalf("Expected success but got error: %v", err)
	}
	
	expected := "test content"
	if string(data) != expected {
		t.Errorf("Expected %q, got %q", expected, string(data))
	} else {
		t.Logf("TestDownloadURL_Success passed")
	}
}

// Test invalid URL format
func TestDownloadURL_InvalidURL(t *testing.T) {
	ctx := context.Background()
	_, err := downloadURL(ctx, "invalid-url")
	if err == nil {
		t.Errorf("Expected error for invalid URL, but got nil")
	} else {
		t.Logf("TestDownloadURL_InvalidURL passed")
	}
}

// Test HTTP error response
func TestDownloadURL_HTTPError(t *testing.T) {
	server := mockHTTPServer("Not Found", http.StatusNotFound)
	defer server.Close()

	ctx := context.Background()
	_, err := downloadURL(ctx, server.URL)

	if err == nil {
		t.Errorf("Expected HTTP error, but got nil")
	} else {
		t.Logf("TestDownloadURL_HTTPError passed")
	}
}

// Test with context cancellation
func TestDownloadURL_ContextCancel(t *testing.T) {
	server := mockHTTPServer("Delayed Response", http.StatusOK)
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := downloadURL(ctx, server.URL)

	if err == nil {
		t.Errorf("Expected context deadline exceeded error, but got nil")
	} else {
		t.Logf("TestDownloadURL_ContextCancel passed")
	}
}

// Test downloadURLs function (concurrent downloads)
func TestDownloadURLs(t *testing.T) {
	server := mockHTTPServer("mock data", http.StatusOK)
	defer server.Close()

	urlChan := make(chan string, 2)
	contentChan := make(chan downloadResult, 2)
	metrics := &Metrics{}
	ctx := context.Background()
	wg := &sync.WaitGroup{}

	// Send URLs to the channel
	urlChan <- server.URL
	urlChan <- server.URL
	close(urlChan)
	
	// Start downloading
	go downloadURLs(urlChan, contentChan, metrics, ctx, wg)

	wg.Wait()
	close(contentChan)

	count := 0
	for range contentChan {
		count++
	}

	if count != 2 {
		t.Errorf("Expected 2 downloads, got %d", count)
	} else {
		t.Logf("TestDownloadURLs passed")
	}
}



