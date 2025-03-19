package src

import (
	"time"
)

const (
	MAX_WORKERS        = 50
	SHUTDOWN_DEAD_LINE = 5 * time.Second // Graceful shutdown deadline

)
