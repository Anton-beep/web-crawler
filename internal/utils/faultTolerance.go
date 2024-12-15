package utils

import (
	"go.uber.org/zap"
	"slices"
	"time"
)

func RetryCount(retries int, pause time.Duration, allowedErrors []error, f func() error) error {
	var err error
	if retries < 1 {
		retries = 1
	}
	for i := 0; i < retries; i++ {
		err = f()
		if err == nil {
			return nil
		}
		if allowedErrors != nil && slices.Contains(allowedErrors, err) {
			return err
		}
		zap.S().Errorf("Error: %v. Retrying...", err)
		time.Sleep(pause)
	}
	return err
}

func RetryTimeout(timeout time.Duration, pause time.Duration, allowedErrors []error, f func() error) error {
	var err error
	if timeout == 0 {
		timeout = time.Second
	}
	start := time.Now()
	for time.Since(start) < timeout {
		err = f()
		if err == nil {
			return nil
		}
		if slices.Contains(allowedErrors, err) {
			return err
		}
		zap.S().Errorf("Error: %v. Retrying...", err)
		time.Sleep(pause)
	}
	return err
}
