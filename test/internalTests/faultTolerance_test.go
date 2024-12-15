package internalTests

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"web-crawler/internal/utils"
)

func TestRetryCount(t *testing.T) {
	counter := 0
	err := utils.RetryCount(3, 0, nil, func() error {
		counter++
		if counter < 3 {
			return errors.New("error")
		}
		return nil
	})
	assert.NoError(t, err)
}

func TestRetryTimeout(t *testing.T) {
	startTime := time.Now()
	err := utils.RetryTimeout(time.Second*3, time.Second, nil, func() error {
		if time.Now().Sub(startTime) < time.Second {
			return errors.New("error")
		}
		return nil
	})
	assert.NoError(t, err)
}
