package utils

import (
	"github.com/brianvoe/gofakeit/v7"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
)

// AddRandomHeaders adds random headers to the request
func AddRandomHeaders(req *http.Request, generator *rand.Rand) {
	err := gofakeit.Seed(generator.Int63())
	if err != nil {
		zap.S().Errorf("Failed to seed gofakeit %s", err)
		return
	}

	req.Header.Set("User-Agent", gofakeit.UserAgent())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "keep-alive")
}
