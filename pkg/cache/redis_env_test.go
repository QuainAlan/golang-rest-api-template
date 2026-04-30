package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRedisOptionsFromEnvDefaults(t *testing.T) {
	t.Setenv("REDIS_ADDR", "")
	t.Setenv("REDIS_HOST", "")
	t.Setenv("REDIS_PORT", "")
	t.Setenv("REDIS_DB", "")
	t.Setenv("REDIS_PASSWORD", "")
	t.Setenv("REDIS_DIAL_TIMEOUT", "")
	t.Setenv("REDIS_READ_TIMEOUT", "")
	t.Setenv("REDIS_WRITE_TIMEOUT", "")
	t.Setenv("REDIS_TLS", "")

	o, err := redisOptionsFromEnv()
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "127.0.0.1:6379", o.Addr)
	assert.Equal(t, 0, o.DB)
	assert.Equal(t, 5*time.Second, o.DialTimeout)
	assert.Equal(t, 3*time.Second, o.ReadTimeout)
	assert.Equal(t, 3*time.Second, o.WriteTimeout)
	assert.Nil(t, o.TLSConfig)
}

func TestRedisOptionsFromEnvAddrOverridesHostPort(t *testing.T) {
	t.Setenv("REDIS_ADDR", "redis.internal:6380")
	t.Setenv("REDIS_HOST", "should-not-matter")
	t.Setenv("REDIS_PORT", "9999")

	o, err := redisOptionsFromEnv()
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "redis.internal:6380", o.Addr)
}

func TestRedisOptionsFromEnvHostAndPort(t *testing.T) {
	t.Setenv("REDIS_ADDR", "")
	t.Setenv("REDIS_HOST", "10.0.0.5")
	t.Setenv("REDIS_PORT", "16379")

	o, err := redisOptionsFromEnv()
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "10.0.0.5:16379", o.Addr)
}

func TestRedisOptionsFromEnvPasswordAndDB(t *testing.T) {
	t.Setenv("REDIS_ADDR", "127.0.0.1:6379")
	t.Setenv("REDIS_PASSWORD", "secret")
	t.Setenv("REDIS_DB", "2")

	o, err := redisOptionsFromEnv()
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "secret", o.Password)
	assert.Equal(t, 2, o.DB)
}

func TestRedisOptionsFromEnvInvalidDB(t *testing.T) {
	t.Setenv("REDIS_ADDR", "127.0.0.1:6379")
	t.Setenv("REDIS_DB", "x")

	_, err := redisOptionsFromEnv()
	assert.Error(t, err)
}

func TestRedisOptionsFromEnvNegativeDB(t *testing.T) {
	t.Setenv("REDIS_ADDR", "127.0.0.1:6379")
	t.Setenv("REDIS_DB", "-1")

	_, err := redisOptionsFromEnv()
	assert.Error(t, err)
}

func TestRedisOptionsFromEnvTLSEnabled(t *testing.T) {
	t.Setenv("REDIS_ADDR", "127.0.0.1:6379")
	t.Setenv("REDIS_TLS", "true")
	t.Setenv("REDIS_TLS_INSECURE", "")

	o, err := redisOptionsFromEnv()
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NotNil(t, o.TLSConfig) {
		return
	}
	assert.False(t, o.TLSConfig.InsecureSkipVerify)
}

func TestRedisOptionsFromEnvTLSInsecure(t *testing.T) {
	t.Setenv("REDIS_ADDR", "127.0.0.1:6379")
	t.Setenv("REDIS_TLS", "1")
	t.Setenv("REDIS_TLS_INSECURE", "true")

	o, err := redisOptionsFromEnv()
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NotNil(t, o.TLSConfig) {
		return
	}
	assert.True(t, o.TLSConfig.InsecureSkipVerify)
}

func TestRedisOptionsFromEnvTruthyTLSOn(t *testing.T) {
	t.Setenv("REDIS_ADDR", "127.0.0.1:6379")
	t.Setenv("REDIS_TLS", "on")
	t.Setenv("REDIS_TLS_INSECURE", "")

	o, err := redisOptionsFromEnv()
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NotNil(t, o.TLSConfig) {
		return
	}
	assert.False(t, o.TLSConfig.InsecureSkipVerify)
}

func TestRedisOptionsFromEnvInvalidDialTimeout(t *testing.T) {
	t.Setenv("REDIS_ADDR", "127.0.0.1:6379")
	t.Setenv("REDIS_DIAL_TIMEOUT", "nope")

	_, err := redisOptionsFromEnv()
	assert.Error(t, err)
}

func TestRedisOptionsFromEnvNonPositiveDialTimeout(t *testing.T) {
	t.Setenv("REDIS_ADDR", "127.0.0.1:6379")
	t.Setenv("REDIS_DIAL_TIMEOUT", "0s")

	_, err := redisOptionsFromEnv()
	assert.Error(t, err)
}
