package src

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// FileExists checks if the given file path exists and is not a directory.
func fileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		return !os.IsNotExist(err) // Returns false if file does not exist, true for other errors
	}
	return !info.IsDir()
}

// GetExeName returns the name of the executable file.
// If an error occurs, it logs the error and returns an empty string.
func GetExeName() string {
	executablePath, err := os.Executable()
	if err != nil {
		log.Printf("Failed to get executable path: %v", err)
		return ""
	}
	return filepath.Base(executablePath)
}

// PrintAndDie prints an error message to stderr and exits with code 1.
func PrintAndDie(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

// PrintHelpAndExit prints the usage instructions and exits with code 0.
func PrintHelpAndExit() {
	fmt.Println(usageStr)
	os.Exit(0)
}

// PrintVersionAndExit prints the executable name and version, then exits with code 0.
func PrintVersionAndExit(version string) {
	fmt.Printf("%s: v%s\n", GetExeName(), version)
	os.Exit(0)
}

// getFileName extracts the file name from the given file path.
func getFileName(filePath string) string {
	return strings.Split(filepath.Base(filePath), ".")[0]
}

// GetFileExtension extracts the file extension from the given file path.
func GetFileExtension(filePath string) string {
	ext := filepath.Ext(filePath)
	return strings.TrimPrefix(ext, ".") // Remove the leading dot
}

func ensureScheme(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "https://" + url // Default to HTTPS
	}
	return url
}
