package src

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// persistContent receives downloaded content from a channel and writes it to files.
//
// Input:
// - contentChan: A channel that provides downloadResult objects containing URL and content.
// - filePath: The base file path used to determine the output directory.
// - ctx: Context for graceful shutdown.
//
// Output:
// - Saves downloaded content as files in the directory `<filePath_without_extension>/downloads/`.
// - Logs errors if file creation or writing fails.
// - Stops processing when the context is canceled.
//
// Notes:
// - Creates an output directory if it doesn’t exist.
// - Uses a random filename for each saved file.
// - Ensures graceful shutdown if the context is canceled.
func persistContent(contentChan <-chan downloadResult, filePath string, ctx context.Context) {
	// Determine the output directory based on the file path
	outputDir := strings.Split(filePath, ".")[0] + "/downloads/"

	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		log.Fatalf("Error creating output directory: %v", err) // Fatal log stops execution on failure
	}

	// Continuously listen for download results
	for {
		select {
		case result, ok := <-contentChan:
			if !ok {
				return // Exit if the channel is closed
			}

			// Generate a random file name and construct the full path
			fileName := filepath.Join(outputDir, generateRandomFileName())

			// Create the output file
			file, err := os.Create(fileName)
			if err != nil {
				zlog.Error().Msgf("Error creating file: %v for URL: %s", err, result.url)
				continue
			}
			defer file.Close() // Ensure file is closed after writing

			// Write content to the file
			_, err = file.WriteString(string(result.content))
			if err != nil {
				zlog.Error().Msgf("Error writing to file: %v for URL: %s", err, result.url)
				continue
			}

			// Log success
			zlog.Info().Msgf("Saved content to %s for URL: %s", fileName, result.url)

		case <-ctx.Done(): // Handle shutdown scenario
			zlog.Info().Msgf("Stage 3: Context canceled. Stopping file write.")
			return
		}
	}
}

// generateRandomFilename creates a unique filename using a random string and a timestamp
func generateRandomFileName() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%d-%v.txt", rand.Intn(1000000), time.Now().UnixNano())
}
