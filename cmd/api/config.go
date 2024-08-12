package api

import (
	"github.com/tomknobel/ip2country/internal/db"
	"github.com/tomknobel/ip2country/pkg/utils"
	"strconv"
	"time"
)

type limiterCfg struct {
	rateLimit  int64
	windowSize time.Duration
}
type config struct {
	limiter limiterCfg
	port    string
	dbType  string
	dbCfg   db.DbConfig
}

func initConfig() config {
	windowSize, err := time.ParseDuration(utils.GetEnv("WINDOW_SIZE", "60s"))
	if err != nil {
		panic(err)
	}
	rateLimit, err := strconv.ParseInt(utils.GetEnv("RATE_LIMIT", "10"), 10, 64)
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}
	return config{
		limiter: limiterCfg{
			rateLimit:  rateLimit,
			windowSize: windowSize,
		},

		dbCfg: db.DbConfig{
			ConnectionString: utils.GetEnv("DB_CONNECTION_STRING", ""),
		},
		dbType: utils.GetEnv("DB_TYPE", "csv"),
		port:   utils.GetEnv("PORT", "3001"),
	}
}
