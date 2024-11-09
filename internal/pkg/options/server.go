package options

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type options struct {
	readTimeout     time.Duration
	writeTimeout    time.Duration
	idleTimeout     time.Duration
	shutdownTimeout int
}

type ServerOption func(options *options) error

const (
	defaultReadTimeout     = 10
	defaultWriteTimeout    = 10
	defaultIdleTimeout     = 10
	defaultShutdownTimeout = 10
	readTimeoutEnv         = "SERVER_READ_TIMEOUT"
	writeTimeoutEnv        = "SERVER_WRITE_TIMEOUT"
	idleTimeoutEnv         = "SERVER_IDLE_TIMEOUT"
)

var ErrNegativeTimeout = errors.New("Negative timeout")

func WithReadTimeout(timeout int) ServerOption {
	return func(options *options) error {
		if timeout < 0 {
			return fmt.Errorf("server read timeout: %w", ErrNegativeTimeout)
		}

		timeOutStr := os.Getenv(readTimeoutEnv)

		readTimeout, err := strconv.Atoi(timeOutStr)

		if err != nil {
			readTimeout = defaultReadTimeout
		}

		options.readTimeout = time.Duration(readTimeout * int(time.Second))

		return nil
	}
}

func WithWriteTimeout(timeout int) ServerOption {
	return func(options *options) error {
		if timeout < 0 {
			return fmt.Errorf("server write timeout: %w", ErrNegativeTimeout)
		}

		timeOutStr := os.Getenv(writeTimeoutEnv)

		writeTimeout, err := strconv.Atoi(timeOutStr)

		if err != nil {
			writeTimeout = defaultWriteTimeout
		}

		options.writeTimeout = time.Duration(writeTimeout * int(time.Second))

		return nil
	}
}

func WithIdleTimeout(timeout int) ServerOption {
	return func(options *options) error {
		if timeout < 0 {
			return fmt.Errorf("server idle timeout: %w", ErrNegativeTimeout)
		}

		timeOutStr := os.Getenv(idleTimeoutEnv)

		idleTimeout, err := strconv.Atoi(timeOutStr)

		if err != nil {
			idleTimeout = defaultIdleTimeout
		}

		options.idleTimeout = time.Duration(idleTimeout * int(time.Second))

		return nil
	}
}
