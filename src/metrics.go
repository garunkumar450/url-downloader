package src

import (
	"log"
	"sync/atomic"
	"time"
)

// Metrics tracks the progress of URL processing
type Metrics struct {
	TotalURLs     atomic.Uint64 // Total number of URLs processed
	SuccessCount  atomic.Uint64 // Number of successful downloads
	FailureCount  atomic.Uint64 // Number of failed downloads
	TotalDuration atomic.Uint64 // Total duration of all successful downloads (in nanoseconds)
	PrcStartTime  time.Time
	PrcEndTime    time.Time
}

func (m *Metrics) AddSuccess(duration time.Duration) {
	m.SuccessCount.Add(1)
	m.TotalDuration.Add(uint64(duration.Nanoseconds()))
}

func (m *Metrics) AddFailure() {
	m.FailureCount.Add(1)
}

func (m *Metrics) LogSummary() {
	totalURLs := m.TotalURLs.Load()
	successCount := m.SuccessCount.Load()
	failureCount := m.FailureCount.Load()
	totalDuration := time.Duration(m.TotalDuration.Load())
	avgDuration := time.Duration(0)
	if successCount > 0 {
		avgDuration = totalDuration / time.Duration(successCount)
	}
	log.Printf("Summary: Total URLs=%d, Success=%d, Failures=%d, Avg Download Duration=%v", totalURLs, successCount, failureCount, avgDuration)
	zlog.Info().Uint64("Total URLs", totalURLs).Uint64("Success", successCount).Uint64("Failures", failureCount).Str("Avg Download Duration", avgDuration.String()).Str("Latency", m.PrcEndTime.Sub(m.PrcStartTime).String()).Msg("Summary")
}
