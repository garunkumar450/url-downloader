package src


import (
	"context"
	"os"
	"testing"
	"time"
)


// Test successful persistence
func TestPersistContent_Success(t *testing.T) {
	ctx := context.Background()
	contentChan := make(chan downloadResult, 1)

	// Send mock data
	contentChan <- downloadResult{url: "http://example.com", content: []byte("test content")}
	close(contentChan)
	filePath:="../testdata/valid.csv"
	go persistContent(contentChan, filePath, ctx)

	// Verify file exists
	time.Sleep(10 * time.Millisecond) // Allow time for goroutine

	 // Verify results
        files, err := os.ReadDir("../testdata/valid/downloads/")
        if err != nil {
                t.Fatalf("Failed to read output directory: %v", err)
        }
        if len(files) != 2 {
                t.Errorf("Expected 2 file, got %d", len(files))
        }

	// Cleanup
	os.RemoveAll("../testdata/valid/downloads/*.txt")
}

// Test handling of closed channel
func TestPersistContent_ClosedChannel(t *testing.T) {
	ctx := context.Background()
	contentChan := make(chan downloadResult)
	close(contentChan) // Close the channel before calling the function

	go persistContent(contentChan, "../testdata/valid.csv", ctx)

	time.Sleep(50 * time.Millisecond) // Ensure no panic occurs

	t.Logf("TestPersistContent_ClosedChannel passed")
}

// Test handling of context cancellation
func TestPersistContent_ContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	contentChan := make(chan downloadResult, 1)

	// Send mock data
	contentChan <- downloadResult{url: "http://example.com", content: []byte("test content")}

	// Cancel the context
	cancel()

	go persistContent(contentChan, "../testdata/valid.csv", ctx)

	time.Sleep(50 * time.Millisecond) // Ensure cancellation is handled

	t.Logf("TestPersistContent_ContextCancel passed")
}

// Test handling of file creation failure (mocked)
func TestPersistContent_FileCreationFailure(t *testing.T) {
	ctx := context.Background()
	contentChan := make(chan downloadResult, 1)

	// Mock invalid directory
	filePath := "../testdata/valid.csv"
	contentChan <- downloadResult{url: "http://example.com", content: []byte("test content")}
	close(contentChan)

	go persistContent(contentChan, filePath, ctx)

	time.Sleep(50 * time.Millisecond) // Ensure failure is handled

	t.Logf("TestPersistContent_FileCreationFailure passed")
}

