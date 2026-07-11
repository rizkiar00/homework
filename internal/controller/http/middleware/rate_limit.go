package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rizkiar00/homework/pkg/constant"
	httpresponse "github.com/rizkiar00/homework/pkg/response"
	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type ipRateLimiter struct {
	mutex       sync.Mutex
	visitors    map[string]*visitor
	limit       rate.Limit
	burst       int
	lastCleanup time.Time
}

func newIPRateLimiter(requestsPerMinute int, burst int) *ipRateLimiter {
	if requestsPerMinute <= 0 {
		requestsPerMinute = 60
	}
	if burst <= 0 {
		burst = requestsPerMinute
	}

	return &ipRateLimiter{
		visitors:    make(map[string]*visitor),
		limit:       rate.Every(time.Minute / time.Duration(requestsPerMinute)),
		burst:       burst,
		lastCleanup: time.Now(),
	}
}

func (l *ipRateLimiter) allow(ip string) bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	now := time.Now()
	if now.Sub(l.lastCleanup) > time.Minute {
		l.cleanup(now)
	}

	current, ok := l.visitors[ip]
	if !ok {
		current = &visitor{
			limiter: rate.NewLimiter(l.limit, l.burst),
		}
		l.visitors[ip] = current
	}

	current.lastSeen = now
	return current.limiter.Allow()
}

func (l *ipRateLimiter) cleanup(now time.Time) {
	for ip, current := range l.visitors {
		if now.Sub(current.lastSeen) > 3*time.Minute {
			delete(l.visitors, ip)
		}
	}
	l.lastCleanup = now
}

func RateLimit(requestsPerMinute int, burst int) echo.MiddlewareFunc {
	limiter := newIPRateLimiter(requestsPerMinute, burst)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if !limiter.allow(ctx.RealIP()) {
				return httpresponse.Error(ctx, http.StatusTooManyRequests, constant.CodeTooManyRequests, constant.MessageTooManyRequests)
			}

			return next(ctx)
		}
	}
}
