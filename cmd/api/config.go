package api

import (
	"github.com/tomknobel/ip2country/pkg/rate"
	"github.com/tomknobel/ip2country/pkg/utils"
	"strconv"
	"time"
)

type Config struct {
	rateLimiter rate.Config
	port        string
}

func initConfig() Config {
	window, err := time.ParseDuration(utils.GetEnv("WINDOW", "2s"))
	if err != nil {
		panic(err)
	}
	ratePerRequest, err := strconv.ParseInt(utils.GetEnv("RATE", "2"), 10, 64)
	if err != nil {
		panic(err)
	}
	rateLimiter := rate.Config{
		Rate:   ratePerRequest,
		Window: window,
	}
	return Config{
		rateLimiter: rateLimiter,
		port:        utils.GetEnv("PORT", "3001"),
	}
}
