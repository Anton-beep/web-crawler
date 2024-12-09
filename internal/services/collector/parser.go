package collector

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/net/html"
	"net/http"
	"strings"
	"web-crawler/internal/utils"
)

// AddTitle adds title to project temporary data
func (s *Server) AddTitle(title string) {
	if len(s.ProjectTemporaryData.Titles) == 0 {
		s.ProjectTemporaryData.Titles = append(s.ProjectTemporaryData.Titles, "")
	}
	s.ProjectTemporaryData.Titles[0] += title + "\n"
}

// AddText adds text to project temporary data
func (s *Server) AddText(text string) {
	if len(s.ProjectTemporaryData.Text) == 0 {
		s.ProjectTemporaryData.Text = append(s.ProjectTemporaryData.Text, "")
	}
	s.ProjectTemporaryData.Text[0] += text + "\n"
}

// AddNode adds node to project temporary data
// If node already exists, or TotalCollectorCounter
// is less than 0, does nothing
func (s *Server) AddNode(link string, depth int) bool {
	actualData, err := s.DataBase.GetProjectTemporaryData(s.Message.ProjectId)
	if err != nil {
		zap.S().Errorf("failed to get actual project temporary data: %s", err)
		return false
	}

	if actualData.TotalCollectorCounter <= 0 {
		return false
	}

	actualData.TotalCollectorCounter--

	err = s.DataBase.SetProjectTemporaryData(s.Message.ProjectId, actualData)
	if err != nil {
		zap.S().Errorf("failed to set actual project temporary data: %s", err)
	}

	slag := GenerateNodeSlug(s.Message.ProjectId, link)

	status, err := s.DataBase.CheckSlug(slag)
	if err != nil {
		zap.S().Errorf("failed to check status of node slug: %s", err)
		return false
	}
	if status {
		return true
	}

	format := `{"id": "%s", "name": "%s", "val": %d},`
	s.ProjectTemporaryData.Nodes += fmt.Sprintf(format, link, link, depth)

	err = s.DataBase.UpdateSlug(slag, true)
	if err != nil {
		zap.S().Errorf("failed to update node slug: %s", err)
	}
	return true
}

// AddLink adds link to project temporary data
func (s *Server) AddLink(link string) {
	if link[len(link)-1] == '/' {
		link = link[:len(link)-1]
	}
	zap.S().Debug("AddLink ", link)
	status, err := s.DataBase.CheckSlug(GenerateLinkSlug(s.Message.ProjectId, link))
	if err != nil {
		zap.S().Errorf("failed to check link status, CheckSlug returned err: %s", err)
	}
	if status {
		format := `{"source": "%s", "target": "%s"},`
		s.ProjectTemporaryData.Links += fmt.Sprintf(format, s.Message.Link, link)
		return
	}

	nodeStatus := s.AddNode(link, s.Message.Depth+1)
	if !nodeStatus {
		return
	}

	format := `{"source": "%s", "target": "%s"},`
	s.ProjectTemporaryData.Links += fmt.Sprintf(format, s.Message.Link, link)

	if s.Message.Depth < s.MaxDepth {
		s.NewCollectors++
		err = s.Broker.AddSiteToParse(link, s.Message.ProjectId, s.Message.Depth+1)
		if err != nil {
			zap.S().Errorf("failed to add site %s to parse, AddSiteToParse returned %s", link, err)
		}
	}
}

// PrepareLink prepares link for
// adding to project temporary data
func (s *Server) PrepareLink(link string) {
	if utils.IsCorrectLink(link) {
		s.AddLink(link)
		return
	}
	if strings.Count(link, "http") == 0 {
		link = s.Domain + link
		if utils.IsCorrectLink(link) {
			s.AddLink(link)
		}
	}
}

// GetNode makes GET request to link and returns html.Node
func (s *Server) GetNode(link string) (*html.Node, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return nil, err
	}

	utils.AddRandomHeaders(req, s.RandomGenerator)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected status code: " + resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

// ParseNodes parses html.Node and adds
func (s *Server) ParseNodes(n *html.Node) {
	if _, ok := s.TextTags[n.Data]; ok {
		if n.FirstChild != nil {
			s.AddText(n.FirstChild.Data)
		}
	}

	if n.Data == "title" {
		if n.FirstChild != nil {
			s.AddTitle(n.FirstChild.Data)
		}
	}

	if n.Data == "a" {
		if n.FirstChild != nil {
			s.AddText(n.FirstChild.Data)
		}
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				s.PrepareLink(attr.Val)
				break
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		s.ParseNodes(c)
	}
}
