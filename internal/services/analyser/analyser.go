package analyser

import (
	"go.uber.org/zap"
	"web-crawler/internal/broker"
	"web-crawler/internal/config"
	"web-crawler/internal/models"
	"web-crawler/internal/repository"
)

type Server struct {
	DataBase models.DataBase
	Broker   *broker.Kafka
	Analyser config.AnalyserConfig
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		DataBase: repository.NewDB(cfg),
		Broker:   broker.New(cfg, true, true, "analyser"),
		Analyser: cfg.Analyser,
	}
}

func (s *Server) Start() {
	for {
		zap.S().Debug("waiting for new message...")
		message, err := s.Broker.Consume()
		zap.S().Debug("received message")
		if err != nil {
			zap.S().Fatalf("failed to consume data for analyser returned err: %s", err)
			continue
		}
		data, err := s.DataBase.GetProjectTemporaryData(message.ProjectId)
		if err != nil {
			zap.S().Errorf("failed to get project temporary data returned err: %s", err)
			continue
		}
		if message.AnalyseType == "KeyWords" {
			resp, err := s.keyWordsPrompt(data.Text)
			if err != nil {
				zap.S().Debug(err)
				continue
			}
			project, err := s.DataBase.GetProject(message.ProjectId)
			if err != nil {
				zap.S().Errorf("failed to get project: %s", err)
				continue
			}
			project.KeyWords = resp.Response
			err = s.DataBase.UpdateProject(project)
			if err != nil {
				zap.S().Errorf("failed to update project: %s", err)
				continue
			}
			project = nil
		} else if message.AnalyseType == "MainIdeas" {
			resp, err := s.mainIdeasPrompt(data.Text)
			if err != nil {
				zap.S().Debug(err)
				continue
			}
			project, err := s.DataBase.GetProject(message.ProjectId)
			if err != nil {
				zap.S().Errorf("failed to get project: %s", err)
				continue
			}
			project.MainIdeas = resp.Response
			err = s.DataBase.UpdateProject(project)
			if err != nil {
				zap.S().Errorf("failed to update project: %s", err)
				continue
			}
			project = nil
		}
	}
}
