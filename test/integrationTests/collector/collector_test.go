package collector

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"os"
	"testing"
	"time"
	"web-crawler/internal/broker"
	"web-crawler/internal/config"
	"web-crawler/internal/models"
	"web-crawler/internal/services/collector"
)

var (
	cfg *config.Config
)

func TestMain(m *testing.M) {
	time.Sleep(1 * time.Second)
	cfg = config.NewConfig("../../../configs/.env")
	config.InitLogger(true)
	code := m.Run()
	os.Exit(code)
}

func TestAddLink(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	s := collector.NewServer(cfg)
	s.Domain = "https://example.com"

	prj := models.Project{
		OwnerID:          uuid.New().String(),
		StartUrl:         "https://example.com",
		MaxNumberOfLinks: 1,
	}

	_, _ = s.DataBase.CreateProject(&prj)

	ptd := models.ProjectTemporaryData{
		Text:                  []string{},
		Titles:                []string{},
		Nodes:                 "",
		Links:                 "",
		TotalCollectorCounter: prj.MaxNumberOfLinks,
	}

	_ = s.DataBase.SetProjectTemporaryData(prj.ID, &ptd)
	s.ProjectTemporaryData = &ptd

	msg := broker.Message{
		Link:      prj.StartUrl,
		ProjectId: prj.ID,
		Depth:     0,
	}

	s.Message = &msg

	s.PrepareLink("https://example.com/")

	assert.Equal(t, s.ProjectTemporaryData.Nodes, `{"id": "https://example.com", "name": "https://example.com", "val": 1},`)
	assert.Equal(t, s.ProjectTemporaryData.Links, `{"source": "https://example.com", "target": "https://example.com"},`)
}

func TestAddShortLink(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	s := collector.NewServer(cfg)
	s.Domain = "https://example.com"

	prj := models.Project{
		OwnerID:          uuid.New().String(),
		StartUrl:         "https://example.com",
		MaxNumberOfLinks: 1,
	}

	_, _ = s.DataBase.CreateProject(&prj)

	ptd := models.ProjectTemporaryData{
		Text:                  []string{},
		Titles:                []string{},
		Nodes:                 "",
		Links:                 "",
		TotalCollectorCounter: prj.MaxNumberOfLinks,
	}

	_ = s.DataBase.SetProjectTemporaryData(prj.ID, &ptd)
	s.ProjectTemporaryData = &ptd

	msg := broker.Message{
		Link:      prj.StartUrl,
		ProjectId: prj.ID,
		Depth:     0,
	}

	s.Message = &msg

	s.PrepareLink("/")

	assert.Equal(t, s.ProjectTemporaryData.Nodes, `{"id": "https://example.com", "name": "https://example.com", "val": 1},`)
	assert.Equal(t, s.ProjectTemporaryData.Links, `{"source": "https://example.com", "target": "https://example.com"},`)
}

func TestProceedMessage(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	s := collector.NewServer(cfg)
	s.Domain = "https://example.com"

	prj := models.Project{
		OwnerID:          uuid.New().String(),
		StartUrl:         "https://example.com/",
		MaxNumberOfLinks: 2,
	}

	_, _ = s.DataBase.CreateProject(&prj)

	ptd := models.ProjectTemporaryData{
		Text:                  []string{},
		Titles:                []string{},
		Nodes:                 "",
		Links:                 "",
		TotalCollectorCounter: prj.MaxNumberOfLinks,
		CollectorCounterQueue: 1,
	}

	_ = s.DataBase.SetProjectTemporaryData(prj.ID, &ptd)
	s.ProjectTemporaryData = &ptd

	msg := broker.Message{
		Link:      prj.StartUrl,
		ProjectId: prj.ID,
		Depth:     0,
	}

	s.Message = &msg

	assert.False(t, s.WasParsed())
	assert.True(t, s.AssignLink() == nil)
	s.Process()

	result, _ := s.DataBase.GetProjectTemporaryData(prj.ID)

	assert.Equal(t, result.Text[0], `Example Domain
Example Domain
This domain is for use in illustrative examples in documents. You may use this
    domain in literature without prior coordination or asking for permission.
a
More information...
`)
	assert.Equal(t, result.Titles[0], "Example Domain\n")
	assert.Equal(t, result.Nodes, `{"id": "https://example.com/", "name": "https://example.com/", "val": 0},{"id": "https://www.iana.org/domains/example", "name": "https://www.iana.org/domains/example", "val": 1}`)
	assert.Equal(t, result.Links, `{"source": "https://example.com/", "target": "https://www.iana.org/domains/example"}`)
	assert.Equal(t, result.TotalCollectorCounter, 0)
}

func TestProceedNonExistingLink(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	s := collector.NewServer(cfg)
	s.Domain = "https://non_existing_link.com"

	prj := models.Project{
		OwnerID:          uuid.New().String(),
		StartUrl:         "https://non_existing_link.com",
		MaxNumberOfLinks: 1,
	}

	_, _ = s.DataBase.CreateProject(&prj)

	ptd := models.ProjectTemporaryData{
		Text:                  []string{},
		Titles:                []string{},
		Nodes:                 "",
		Links:                 "",
		TotalCollectorCounter: prj.MaxNumberOfLinks,
		CollectorCounterQueue: 1,
	}

	_ = s.DataBase.SetProjectTemporaryData(prj.ID, &ptd)
	s.ProjectTemporaryData = &ptd

	msg := broker.Message{
		Link:      prj.StartUrl,
		ProjectId: prj.ID,
		Depth:     0,
	}

	s.Message = &msg

	assert.False(t, s.WasParsed())
	err := s.AssignLink()
	zap.S().Debug(err)
	assert.True(t, err == nil)
	s.Process()

	result, _ := s.DataBase.GetProjectTemporaryData(prj.ID)

	assert.Equal(t, result.Text, []string{})
	assert.Equal(t, result.Titles, []string{})
	assert.Equal(t, result.Nodes, "")
	assert.Equal(t, result.Links, "")
	assert.Equal(t, result.TotalCollectorCounter, 0)
	assert.Equal(t, result.DeadListQueueSites, []string{prj.StartUrl})
}

func TestNegativeCounter(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	s := collector.NewServer(cfg)
	s.Domain = "example.com"

	prj := models.Project{
		OwnerID:          uuid.New().String(),
		StartUrl:         "example.com",
		MaxNumberOfLinks: 1,
	}

	_, _ = s.DataBase.CreateProject(&prj)

	ptd := models.ProjectTemporaryData{
		Text:                  []string{},
		Titles:                []string{},
		Nodes:                 "",
		Links:                 "",
		TotalCollectorCounter: 0,
		CollectorCounterQueue: 1,
	}

	_ = s.DataBase.SetProjectTemporaryData(prj.ID, &ptd)
	s.ProjectTemporaryData = &ptd

	msg := broker.Message{
		Link:      prj.StartUrl,
		ProjectId: prj.ID,
		Depth:     0,
	}

	s.Message = &msg

	assert.False(t, s.WasParsed())
	err := s.AssignLink()
	assert.Error(t, models.CollectorCounterIsNegative, err)
}
