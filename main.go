package main

import (
	"fmt"
	"github.com/garunkumar450/url-downloader/src"
	"log"
	"os"
)

func main() {
	// Configure the options from the flags
	err := src.Configure(os.Args[1:])
	if err != nil {
		src.PrintAndDie(fmt.Sprintf("%s: %s", src.GetExeName(), err))
	}
	// This is a blocking call
	err = src.Start()
	if err != nil {
		src.PrintAndDie(fmt.Sprintf("%s: %s", src.GetExeName(), err))
	}
	log.Printf("Application closed")

}
