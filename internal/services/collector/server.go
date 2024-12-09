package collector

import (
	"errors"
	"go.uber.org/zap"
	"math/rand"
	"strings"
	"time"
	"web-crawler/internal/broker"
	"web-crawler/internal/config"
	"web-crawler/internal/models"
	"web-crawler/internal/repository"
)

type Server struct {
	DataBase             models.DataBase
	Broker               *broker.SitesKafka
	ProjectTemporaryData *models.ProjectTemporaryData
	DeadListSites        []string
	TextTags             map[string]bool
	Domain               string
	RandomGenerator      *rand.Rand
	Message              *broker.Message
	MaxDepth             int
	NewCollectors        int
}

func NewServer(cfg *config.Config) *Server {
	tagMap := make(map[string]bool)
	for _, tag := range strings.Split(cfg.Collector.Tags, ",") {
		tagMap[tag] = true
	}
	return &Server{
		DataBase:             repository.NewDB(cfg),
		Broker:               broker.New(cfg, true, true),
		ProjectTemporaryData: &models.ProjectTemporaryData{},
		TextTags:             tagMap,
		Domain:               "",
		RandomGenerator:      rand.New(rand.NewSource(time.Now().UnixNano())),
		Message:              nil,
		MaxDepth:             0,
		NewCollectors:        0,
	}
}

// Start starts endless loop for parsing sites
func (s *Server) Start() {
	for {
		zap.S().Debug("waiting for new message...")
		message, err := s.Broker.CheckSitesToParse()
		zap.S().Debug("received message")
		if err != nil {
			zap.S().Fatalf("failed to check sites to parse, CheckSitesToParse returned err: %s", err)
		}
		if message.Link[len(message.Link)-1] == '/' {
			message.Link = message.Link[:len(message.Link)-1]
		}
		s.Message = message
		if s.WasParsed() {
			s.WriteData()
			s.Clear()
			continue
		}

		err = s.AssignLink()
		if err != nil {
			if errors.Is(err, models.CollectorCounterIsNegative) {
				continue
			} else {
				zap.S().Fatalf("failed to assign link: %s", err)
			}
		}

		s.Process()
		zap.S().Debug("finished collecting ", message.Link)
	}
}
