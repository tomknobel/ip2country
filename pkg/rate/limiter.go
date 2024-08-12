package rate

import "time"

type ILimiter interface {
	Allow() bool
}

type Config struct {
	Rate   int64
	Window time.Duration
}
