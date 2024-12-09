package collector

import (
	"go.uber.org/zap"
	"strings"
	"web-crawler/internal/models"
)

func GenerateLinkSlug(projectId, link string) string {
	return projectId + "#" + link
}

func GenerateNodeSlug(projectId, link string) string {
	return projectId + "#id#" + link
}

// GetDomain returns domain from url
func GetDomain(url string) string {
	if strings.Count(url, "http") > 0 {
		tags := strings.Split(url, "/")
		return tags[0] + "//" + tags[2]
	}
	return "http://" + url
}

// WasParsed checks if link was parsed
func (s *Server) WasParsed() bool {
	slug := GenerateLinkSlug(s.Message.ProjectId, s.Message.Link)
	status, err := s.DataBase.CheckSlug(slug)
	if err != nil {
		zap.S().Fatalf("failed to check link status, CheckSlug returned err: %s", err)
	}
	return status
}

// AssignLink assigns link to collector
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
	s.Domain = GetDomain(s.Message.Link)

	if strings.Count(s.Message.Link, "http") == 0 {
		s.Message.Link = "http://" + s.Message.Link
	}

	s.MaxDepth, err = s.DataBase.GetProjectMaxDepth(s.Message.ProjectId)
	return err
}

// WriteData writes data to cache database
func (s *Server) WriteData() {
	actualData, err := s.DataBase.GetProjectTemporaryData(s.Message.ProjectId)
	if err != nil {
		zap.S().Errorf("failed to get actual project temporary data: %s", err)
	}

	actualData.Titles = append(actualData.Titles, s.ProjectTemporaryData.Titles...)
	actualData.Text = append(actualData.Text, s.ProjectTemporaryData.Text...)
	actualData.Links += s.ProjectTemporaryData.Links
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
			prj.WebGraph = `{"nodes": [` + actualData.Nodes + `], "links": [` + actualData.Links + `]}`
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

// Clear clears server data
func (s *Server) Clear() {
	s.ProjectTemporaryData = &models.ProjectTemporaryData{
		Text:                  []string{},
		Titles:                []string{},
		Nodes:                 "",
		Links:                 "",
		TotalCollectorCounter: 0,
		CollectorCounterQueue: 0,
	}
	s.DeadListSites = []string{}
	s.Domain = ""
	s.Message = nil
	s.NewCollectors = 0
}

// Process makes request to site and parses it
func (s *Server) Process() {
	s.AddNode(s.Message.Link, s.Message.Depth)
	node, err := s.GetNode(s.Message.Link)
	if err != nil {
		zap.S().Debugf("failed to get node. Parsing failed: %s", err)
		s.DeadListSites = append(s.DeadListSites, s.Message.Link)

		s.ProjectTemporaryData = &models.ProjectTemporaryData{
			Text:                  []string{},
			Titles:                []string{},
			Nodes:                 "",
			Links:                 "",
			TotalCollectorCounter: 0,
			CollectorCounterQueue: 0,
		}
		s.NewCollectors = 0

		s.WriteData()
		s.Clear()
		return
	}
	zap.S().Debug("Got node")
	s.ParseNodes(node)
	zap.S().Debug("Parsed node")
	s.WriteData()
	zap.S().Debug("Data was written")
	s.Clear()
}
