package server

import (
	"os"
	"strconv"
	"time"
)

type httpConfig struct {
	port                 int
	httpWriteTimeout     time.Duration
	httpReadTimeout      time.Duration
	httpIdleTimeout      time.Duration
	httpGracefulShutdown time.Duration
}

var config httpConfig

func init() {
	config = httpConfig{}

	portEnv := os.Getenv("PORT")
	if portEnv == "" {
		panic("PORT env is required")
	} else {
		port, err := strconv.Atoi(portEnv)
		if err != nil {
			panic("PORT env is not a valid number")
		}
		if port < 0 || port > 65535 {
			panic("PORT env is out of range 0-65535")
		}
		config.port = port
	}

	httpWriteTimeoutEnv := os.Getenv("HTTP_WRITE_TIMEOUT")
	if httpWriteTimeoutEnv == "" {
		config.httpWriteTimeout = 5 * time.Second
	} else {
		httpWriteTimeout, err := time.ParseDuration(httpWriteTimeoutEnv)
		if err != nil {
			panic("HTTP_WRITE_TIMEOUT env is not a valid duration")
		}
		config.httpWriteTimeout = httpWriteTimeout
	}

	httpReadTimeoutEnv := os.Getenv("HTTP_READ_TIMEOUT")
	if httpReadTimeoutEnv == "" {
		config.httpReadTimeout = 5 * time.Second
	} else {
		httpReadTimeout, err := time.ParseDuration(httpReadTimeoutEnv)
		if err != nil {
			panic("HTTP_READ_TIMEOUT env is not a valid duration")
		}
		config.httpReadTimeout = httpReadTimeout
	}

	httpIdleTimeoutEnv := os.Getenv("HTTP_IDLE_TIMEOUT")
	if httpIdleTimeoutEnv == "" {
		config.httpIdleTimeout = 5 * time.Second
	} else {
		httpIdleTimeout, err := time.ParseDuration(httpIdleTimeoutEnv)
		if err != nil {
			panic("HTTP_IDLE_TIMEOUT env is not a valid duration")
		}
		config.httpIdleTimeout = httpIdleTimeout
	}

	httpGracefulShutdownEnv := os.Getenv("HTTP_GRACEFUL_SHUTDOWN")
	if httpGracefulShutdownEnv == "" {
		config.httpGracefulShutdown = 5 * time.Second
	} else {
		httpGracefulShutdown, err := time.ParseDuration(httpGracefulShutdownEnv)
		if err != nil {
			panic("HTTP_GRACEFUL_SHUTDOWN env is not a valid duration")
		}
		config.httpGracefulShutdown = httpGracefulShutdown
	}
}
