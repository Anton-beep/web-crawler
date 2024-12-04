package internalTests

import (
	"testing"
	"web-crawler/internal/config"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigWithDevConfig(t *testing.T) {
	path := "../../configs/.env"

	cfg := config.NewConfig(path)
	assert.Truef(t, cfg != nil, "Config shouldn't be nil")
}

func TestNewConfigPanics(t *testing.T) {
	invalidPath := "invalid_path.env"

	defer func() {
		if r := recover(); r != nil {
			assert.Contains(t, r, "failed to read config", "The panic should contain a message about the configuration being unable to be read")
		} else {
			t.Errorf("Panic was expected, but it didn’t happen")
		}
	}()

	config.NewConfig(invalidPath)
}
