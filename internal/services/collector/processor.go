package collector

import (
	"go.uber.org/zap"
	"strings"
	"time"
	"web-crawler/internal/models"
)

func GenerateLinkSlug(projectId, link string) string {
	return projectId + "#" + link
}

func GenerateNodeSlug(projectId, link string) string {
	return projectId + "#id#" + link
}

func (s *Server) WasParsed() bool {
	time.Sleep(1 * time.Second)
	slug := GenerateLinkSlug(s.Message.ProjectId, s.Message.Link)
	status, err := s.DataBase.CheckSlug(slug)
	if err != nil {
		zap.S().Fatalf("failed to check link status, CheckSlug returned err: %s", err)
	}
	return status
}

func (s *Server) AssignLink() error {
	slug := GenerateLinkSlug(s.Message.ProjectId, s.Message.Link)
	err := s.DataBase.UpdateSlug(slug, true)
	if err != nil {
		return err
	}

	err = s.DataBase.CheckCollectorCounter(s.Message.ProjectId)
	if err != nil {
		return err
	}
	s.ProjectTemporaryData = &models.ProjectTemporaryData{}
	tags := strings.Split(s.Message.Link, "/")
	s.Domain = tags[0] + "//" + tags[2]
	s.MaxDepth, err = s.DataBase.GetProjectMaxDepth(s.Message.ProjectId)
	return err
}

func (s *Server) WriteData() {
	actualData, err := s.DataBase.GetProjectTemporaryData(s.Message.ProjectId)
	if err != nil {
		zap.S().Errorf("failed to get actual project temporary data: %s", err)
	}

	actualData.Titles += s.ProjectTemporaryData.Titles
	actualData.Text += s.ProjectTemporaryData.Text
	actualData.Links += s.ProjectTemporaryData.Links
	actualData.TotalCollectorCounter--
	actualData.Nodes += s.ProjectTemporaryData.Nodes
	actualData.CollectorCounterQueue += s.NewCollectors - 1
	actualData.DeadListQueueSites = append(actualData.DeadListQueueSites, s.DeadListSites...)

	if actualData.TotalCollectorCounter == 0 || actualData.CollectorCounterQueue == 0 {
		actualData.TotalCollectorCounter = 0
		actualData.CollectorCounterQueue = 0
		actualData.Links = strings.Trim(actualData.Links, ",")
		actualData.Nodes = strings.Trim(actualData.Nodes, ",")
		prj, err := s.DataBase.GetProject(s.Message.ProjectId)
		if err != nil {
			zap.S().Errorf("failed to get project: %s", err)
		} else {
			prj.Processing = false
			prj.WebGraph = `{"node": [` + actualData.Nodes + `], "link": [` + actualData.Links + `]}`
			prj.DlqSites = actualData.DeadListQueueSites
			err = s.DataBase.UpdateProject(prj)
			if err != nil {
				zap.S().Errorf("failed to update project: %s", err)
			}
		}
	}

	err = s.DataBase.SetProjectTemporaryData(s.Message.ProjectId, actualData)
	if err != nil {
		zap.S().Errorf("failed to set actual project temporary data: %s", err)
	}
}

func (s *Server) Clear() {
	s.ProjectTemporaryData = &models.ProjectTemporaryData{
		Text:                  "",
		Titles:                "",
		Nodes:                 "",
		Links:                 "",
		TotalCollectorCounter: 0,
		CollectorCounterQueue: 0,
	}
	s.Domain = ""
	s.Message = nil
	s.NewCollectors = 0
}

func (s *Server) Process() {
	s.AddNode(s.Message.Link, s.Message.Depth)
	node, err := s.GetNode(s.Message.Link)
	if err != nil {
		s.DeadListSites = append(s.DeadListSites, s.Message.Link)
		zap.S().Debugf("failed to get node. Parsing failed: %s", err)
		return
	}
	zap.S().Debug("Got node")
	s.ParseNodes(node)
	zap.S().Debug("Parsed node")
	s.WriteData()
	zap.S().Debug("Data was written")
	s.Clear()
}
