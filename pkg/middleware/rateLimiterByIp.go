package middleware

import (
	"encoding/json"
	"github.com/tomknobel/ip2country/pkg/rate"
	"go.uber.org/zap"
	"net"
	"net/http"
	"sync"
	"time"
)

type rateLimitMiddleware struct {
	cfg    rate.Config
	logger *zap.SugaredLogger
}

type message struct {
	status string
	body   string
}

func NewRateLimitMiddleware(rateLimit int64, window time.Duration, logger *zap.SugaredLogger) *rateLimitMiddleware {
	return &rateLimitMiddleware{
		cfg: rate.Config{
			Rate:   rateLimit,
			Window: window,
		},
		logger: logger,
	}

}
func (rm *rateLimitMiddleware) RateLimiterByIp(next http.Handler) http.Handler {
	rm.logger.Info("innn0")

	type client struct {
		limiter  rate.ILimiter
		lastSeen time.Time
	}
	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)
	go func() {
		for {
			time.Sleep(time.Minute)
			// Lock the mutex to protect this section from race conditions.
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the IP address from the request.
		rm.logger.Info("innn")
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		rm.logger.Info(ip)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Lock the mutex to protect this section from race conditions.
		mu.Lock()
		if _, found := clients[ip]; !found {
			clients[ip] = &client{limiter: rate.NewLimiter(rm.cfg.Window, rm.cfg.Rate)}
		}
		clients[ip].lastSeen = time.Now()
		if !clients[ip].limiter.Allow() {
			mu.Unlock()

			message := message{
				status: "Request Failed",
				body:   "The API is at capacity, try again later.",
			}

			w.WriteHeader(http.StatusTooManyRequests)
			_ = json.NewEncoder(w).Encode(&message)
			return
		}
		mu.Unlock()
		next.ServeHTTP(w, r)
	})
}
