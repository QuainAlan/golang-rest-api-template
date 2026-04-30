package cache

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Keys(context.Context, string) *redis.StringSliceCmd
	Del(context.Context, ...string) *redis.IntCmd
}

// NewRedisClient builds a go-redis client from environment variables.
// See redisOptionsFromEnv for supported keys and defaults.
func NewRedisClient() (*redis.Client, error) {
	opts, err := redisOptionsFromEnv()
	if err != nil {
		return nil, err
	}
	return redis.NewClient(opts), nil
}

func redisOptionsFromEnv() (*redis.Options, error) {
	opts := &redis.Options{
		Password: strings.TrimSpace(os.Getenv("REDIS_PASSWORD")),
	}

	if u := strings.TrimSpace(os.Getenv("REDIS_USERNAME")); u != "" {
		opts.Username = u
	}

	addr := strings.TrimSpace(os.Getenv("REDIS_ADDR"))
	if addr != "" {
		opts.Addr = addr
	} else {
		host := strings.TrimSpace(os.Getenv("REDIS_HOST"))
		if host == "" {
			host = "127.0.0.1"
		}
		port := strings.TrimSpace(os.Getenv("REDIS_PORT"))
		if port == "" {
			port = "6379"
		}
		opts.Addr = host + ":" + port
	}

	if s := strings.TrimSpace(os.Getenv("REDIS_DB")); s != "" {
		db, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("REDIS_DB: %w", err)
		}
		if db < 0 {
			return nil, fmt.Errorf("REDIS_DB must be >= 0, got %d", db)
		}
		opts.DB = db
	}

	dial, err := durationFromEnv("REDIS_DIAL_TIMEOUT", 5*time.Second)
	if err != nil {
		return nil, err
	}
	read, err := durationFromEnv("REDIS_READ_TIMEOUT", 3*time.Second)
	if err != nil {
		return nil, err
	}
	write, err := durationFromEnv("REDIS_WRITE_TIMEOUT", 3*time.Second)
	if err != nil {
		return nil, err
	}
	opts.DialTimeout = dial
	opts.ReadTimeout = read
	opts.WriteTimeout = write

	tlsVal := strings.ToLower(strings.TrimSpace(os.Getenv("REDIS_TLS")))
	if tlsVal == "1" || tlsVal == "true" || tlsVal == "yes" {
		cfg := &tls.Config{MinVersion: tls.VersionTLS12}
		if insecure := strings.ToLower(strings.TrimSpace(os.Getenv("REDIS_TLS_INSECURE"))); insecure == "1" || insecure == "true" || insecure == "yes" {
			cfg.InsecureSkipVerify = true
		}
		opts.TLSConfig = cfg
	}

	return opts, nil
}

func durationFromEnv(key string, defaultVal time.Duration) (time.Duration, error) {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return defaultVal, nil
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return 0, fmt.Errorf("%s must be a Go duration (e.g. 5s): %w", key, err)
	}
	if d <= 0 {
		return 0, fmt.Errorf("%s must be positive, got %s", key, v)
	}
	return d, nil
}
