package src

import (
	"context"
	"github.com/rs/zerolog"
	"sync"
	"time"
)

var (
	zlog        zerolog.Logger
	wg          sync.WaitGroup // main wait group
	urlChan     = make(chan string, MAX_WORKERS)
	contentChan = make(chan downloadResult, MAX_WORKERS)
	metrics     *Metrics

)

func Start() error {
	var err error

	err, zlog = initLogger(csvFilePath)
	if err != nil {
		return err
	}
	// Create a context with a 5-second timeout for graceful shutdowt
	ctx, cancel := context.WithTimeout(context.Background(), SHUTDOWN_DEAD_LINE)
	defer cancel()
	metrics = &Metrics{}
	metrics.PrcStartTime = time.Now()

	// Stage 1: Read file
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(urlChan)
		zlog.Info().Msg("Stage-1 Started Reading Csv file")
		readCSVFile(csvFilePath, urlChan, metrics, ctx)
		zlog.Info().Msg("Stage-1 Completed ")
	}()

	// Stage 2: Download URLs
	wg.Add(1)
	go func() {
		defer wg.Done()
		zlog.Info().Msg("Stage-2 Started  download URLS")
		downloadURLs(urlChan, contentChan, metrics, ctx, &wg)
		zlog.Info().Msg("Stage-2 Completed ")
	}()

	// Close contentChan once all download goroutines are done
	go func() {
		wg.Wait()
		close(contentChan)
	}()




	// Stage 3: Persist Contents (Single Goroutine)
	var persistWg sync.WaitGroup
	persistWg.Add(1)
	go func() {
		zlog.Info().Msg("Stage-3 Started  Persistent")
		defer persistWg.Done()
		persistContent(contentChan, csvFilePath, ctx)
		zlog.Info().Msg("Stage-3 Completed ")
	}()
	persistWg.Wait()

	metrics.PrcEndTime = time.Now()
	metrics.LogSummary()

	// Graceful shutdown
	select {
	case <-ctx.Done():
		zlog.Info().Msg("Shutdown deadline reached. Exiting...")
	default:
		zlog.Info().Msg("All tasks completed. Exiting...")
	}

	return nil
}
