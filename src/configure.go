package src

import (
	"flag"
	"fmt"
)

const (
	VERSION     = "1.0.0"
	MODULE_NAME = "url-downloader"
)

var usageStr = `
Usage: url-downloader [options]

Command line options: (Mandatory)
        -f, --file <file> absolute path of csv file.
Other Options:
	-h, --help	Show this message
	-v, --version	Show version
`

var (
	showVersion bool
	showHelp    bool
	csvFilePath string
)

// ConfigureOptions accepts a flag set and augments it with URL Downloaded
// specific flags. On success, an options structure is returned configured
// based on the selected flags and/or configuration file.
// The command line options take precedence to the ones in the configuration file.
func Configure(args []string) error {
	// Create a FlagSet and sets the usage
	fs := flag.NewFlagSet(MODULE_NAME, flag.ExitOnError)
	fs.Usage = PrintHelpAndExit
	fs.BoolVar(&showHelp, "h", false, "Show this message")
	fs.BoolVar(&showHelp, "help", false, "Show this message")
	fs.BoolVar(&showVersion, "v", false, "Show version")
	fs.BoolVar(&showVersion, "version", false, "Show version")
	fs.StringVar(&csvFilePath, "f", "", "absolute path of csv file")
	fs.StringVar(&csvFilePath, "file", "", "absolute path of csv file")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if showVersion {
		PrintVersionAndExit(VERSION)
	}

	if showHelp {
		fs.Usage()
	}

	if csvFilePath=="" && len(args)>0{
		csvFilePath=args[0]
	}

	if err := postValidator(); err != nil {
		PrintAndDie(err.Error())
	}

	return nil
}

func postValidator() error {
	if csvFilePath == "" {
		return fmt.Errorf("csv filepath is mandatory")
	}
	if !fileExists(csvFilePath) {
		return fmt.Errorf("csv filepath is not found :%s", csvFilePath)
	}
	if GetFileExtension(csvFilePath) != "csv" {
		return fmt.Errorf("invalid extension")
	}
	return nil
}
