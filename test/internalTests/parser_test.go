package internalTests

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
	"web-crawler/internal/broker"
	"web-crawler/internal/models"
	"web-crawler/internal/services/collector"
)

func TestAddTitle(t *testing.T) {
	server := collector.Server{
		ProjectTemporaryData: &models.ProjectTemporaryData{
			Text:                  []string{},
			Titles:                []string{},
			Nodes:                 "",
			Links:                 "",
			TotalCollectorCounter: 0,
		},
		RandomGenerator: rand.New(rand.NewSource(time.Now().UnixNano())),
		Message:         &broker.Message{},
	}

	server.AddTitle("test title")

	assert.Equal(t, server.ProjectTemporaryData.Titles, []string{"test title\n"}, "Expected to add title to titles")
}

func TestAddText(t *testing.T) {
	server := collector.Server{
		ProjectTemporaryData: &models.ProjectTemporaryData{
			Text:                  []string{},
			Titles:                []string{},
			Nodes:                 "",
			Links:                 "",
			TotalCollectorCounter: 0,
		},
		RandomGenerator: rand.New(rand.NewSource(time.Now().UnixNano())),
		Message:         &broker.Message{},
	}

	server.AddText("test text")

	assert.Equal(t, server.ProjectTemporaryData.Text, []string{"test text\n"}, "Expected to add title to titles")
}
