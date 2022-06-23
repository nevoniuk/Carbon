package clickhouse

import (
	"math/rand"
	"strings"
	"time"
)

// TotalRetryDuration returns the total duration to retry before giving up.
var TotalRetryDuration = time.Minute

// MaxRetryDuration returns the maximum duration to wait before retrying.
var MaxRetryDuration = 10 * time.Second

// WithRetries executes a function with retries.
func WithRetries(f func() error) error {
	err := f()
	if !isRetriable(err) {
		return err
	}
	tries := 0
	then := time.Now()
	for {
		tries++
		time.Sleep(retryPause(tries))
		err = f()
		if !isRetriable(err) || time.Since(then) > TotalRetryDuration {
			return err
		}
	}
}

// retryPause returns the duration to wait before retrying.
func retryPause(attempt int) time.Duration {
	max := int64(time.Duration(attempt*2) * time.Second)
	// Select a duration between 1ns and the current max. It might seem
	// counterintuitive to have so much jitter, but
	// https://www.awsarchitectureblog.com/2015/03/backoff.html argues that
	// that is the best strategy.
	pause := time.Duration(1 + rand.Int63n(max))
	if pause > MaxRetryDuration {
		return MaxRetryDuration
	}
	return pause
}

// isRetriable returns true if the error is retriable.
func isRetriable(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "broken pipe") ||
		strings.Contains(msg, "connection reset by peer") ||
		strings.Contains(msg, "EOF")
}