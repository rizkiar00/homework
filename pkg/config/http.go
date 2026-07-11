package config

import (
	"strings"
	"time"
)

type HTTPConfig struct {
	CORSAllowedOrigins             string `env:"HTTP_CORS_ALLOWED_ORIGINS" envDefault:"http://localhost:3000,http://127.0.0.1:3000"`
	TimeoutSeconds                 int    `env:"HTTP_TIMEOUT_SECONDS" envDefault:"30"`
	BodyLimit                      string `env:"HTTP_BODY_LIMIT" envDefault:"1M"`
	RateLimitRequestsPerMinute     int    `env:"HTTP_RATE_LIMIT_REQUESTS_PER_MINUTE" envDefault:"60"`
	RateLimitBurst                 int    `env:"HTTP_RATE_LIMIT_BURST" envDefault:"60"`
	AuthRateLimitRequestsPerMinute int    `env:"HTTP_AUTH_RATE_LIMIT_REQUESTS_PER_MINUTE" envDefault:"5"`
	AuthRateLimitBurst             int    `env:"HTTP_AUTH_RATE_LIMIT_BURST" envDefault:"5"`
}

func (c HTTPConfig) AllowedOrigins() []string {
	origins := strings.Split(c.CORSAllowedOrigins, ",")
	result := make([]string, 0, len(origins))

	for _, origin := range origins {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			result = append(result, origin)
		}
	}

	return result
}

func (c HTTPConfig) Timeout() time.Duration {
	if c.TimeoutSeconds <= 0 {
		return 30 * time.Second
	}

	return time.Duration(c.TimeoutSeconds) * time.Second
}

func (c HTTPConfig) BodyLimitValue() string {
	if c.BodyLimit == "" {
		return "1M"
	}

	return c.BodyLimit
}
