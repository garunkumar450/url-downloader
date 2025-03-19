#URLS Downloader Management Based on CSV File



## Overview
	This Project implements a URL Downloader management for downloading the content of URLS based on command line arguments

## Project Setup

### Prerequisites

Before setting up the project, ensure you have the following installed:

- Go (version 1.23 or newer)


### Getting Started

1. **Clone the repository:**
    ```
    git clone https://github.com/garunkumar450/url-downloader.git

    cd url-downloader
    ```

2. **Initialize Go Modules:**
    ```
    If there are no go.mod and go.sum files, then initialize as per below statment

    go mod init github.com/garunkumar450/url-downloader
    ```

3. **Install Dependencies:**
    ```
    Run below command to download and install all go dependencies

    go mod tidy
    ```
4. ** Help and Version  Command:**
	```
	go run main.go --help
		Usage: url-downloader [options]
		Command line options: (Mandatory)
		        -f, --file <file> absolute path of csv file.
		Other Options:
		        -h, --help      Show this message
		        -v, --version   Show version
	```
	go run main.go --version

5. **Run the Application:**

    To run the application, use the following command:

    go run main.go -f <absolute_path of csv file>
 



5. **Running Unit Tests::**
    ```
    To run unit tests, use the following command:

    go test -v ./...
    ```



### Folder Structure
        - `main.go`: Entry point of the application.
        - `src/configure.go`: commandline arguments parsing ang basic validations
        - `src/app.go`:pipeline starts from here
        - `src/reader.go`:Logic for reading URLs from a CSV file
        - `src/downloader.go`:Main logic for orchestrating the download process.
        - `src/persister.go`:Logic for writing downloaded content to files
        - `src/metrics.go`: Logic for tracking and logging metrics
        - `src/constants.go`:constants
        - `src/utils.go`:Utility functions




