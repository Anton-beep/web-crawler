package integrationTests

import (
	"experiments/internal/config"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestNginx(t *testing.T) {
	cfg := config.NewConfig()
	if !cfg.IsRanInDocker {
		return
	}

	resp, err := http.Get("http://nginx:80")
	assert.Equal(t, err, nil)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, err, nil)

	assert.Equal(t, string(body), "Hello World!")
}
