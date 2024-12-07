package internalTests

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net/http"
	"testing"
	"time"
	"web-crawler/internal/utils"
)

func IsCorrectLinkTest(t *testing.T) {
	assert.True(t, utils.IsCorrectLink("https://google.com"), "Must be valid")
	assert.True(t, utils.IsCorrectLink("http://example.com"), "Must be valid")
	assert.True(t, utils.IsCorrectLink("http://example.com/"), "Must be valid")
	assert.True(t, utils.IsCorrectLink("http://localhost:8080"), "Must be valid")
	assert.True(t, utils.IsCorrectLink("http://localhost:8080/"), "Must be valid")
	assert.True(t, utils.IsCorrectLink("http://127.0.0.1:9090"), "Must be valid")
	assert.True(t, utils.IsCorrectLink("http://127.0.0.1:9090/"), "Must be valid")
	assert.True(t, utils.IsCorrectLink("https://sub.domain.com/path?query=1"), "Must be valid")
	assert.True(t, utils.IsCorrectLink("http://example.com:8000/path#fragment"), "Must be valid")
	assert.True(t, utils.IsCorrectLink("http://[::1]:8080"), "Must be valid for IPv6")

	assert.False(t, utils.IsCorrectLink("ftp://example.com"), "Should not be valid")
	assert.False(t, utils.IsCorrectLink("http://invalid_url"), "Should not be valid")
	assert.False(t, utils.IsCorrectLink("http://example..com"), "Should not be valid")
	assert.False(t, utils.IsCorrectLink("http://"), "Should not be valid")
	assert.False(t, utils.IsCorrectLink("localhost:8080"), "Shouldn't be valid without a schema")
	assert.False(t, utils.IsCorrectLink("https:/example.com"), "Should not be valid due to incorrect schema")
	assert.False(t, utils.IsCorrectLink("randomstring"), "Should not be valid for a random string")
}

func TestAddRandomHeadersDeterministic(t *testing.T) {
	generator := rand.New(rand.NewSource(42))

	req, _ := http.NewRequest("GET", "http://example.com", nil)

	utils.AddRandomHeaders(req, generator)

	expectedUserAgent := "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/5322 (KHTML, like Gecko) Chrome/36.0.822.0 Mobile Safari/5322"
	assert.Equal(t, expectedUserAgent, req.Header.Get("User-Agent"), "User-Agent must match a fixed seed")
}

func TestAddRandomHeadersRandomness(t *testing.T) {
	generator := rand.New(rand.NewSource(time.Now().UnixNano()))

	userAgents := make([]string, 5)

	for i := 0; i < 5; i++ {
		req, _ := http.NewRequest("GET", "http://example.com", nil)
		utils.AddRandomHeaders(req, generator)
		userAgents[i] = req.Header.Get("User-Agent")
	}

	allEqual := true
	for i := 1; i < len(userAgents); i++ {
		if userAgents[i] != userAgents[0] {
			allEqual = false
			break
		}
	}

	assert.False(t, allEqual, "User-Agent values must not be the same for 5 consecutive calls")
}
