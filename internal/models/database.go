package models

type Project struct {
	ID               string   `db:"id"`
	OwnerID          string   `db:"owner_id"`
	Name             string   `db:"name"`
	StartUrl         string   `db:"start_url"`
	Processing       bool     `db:"processing"`
	WebGraph         string   `db:"web_graph"`
	DlqSites         []string `db:"dlq_sites"`
	MaxDepth         int      `db:"max_depth"`
	MaxNumberOfLinks int      `db:"max_number_of_links"`
}

type ShortProject struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ProjectTemporaryData struct {
	Text                  string   `json:"text"`
	Titles                string   `json:"titles"`
	Nodes                 string   `json:"nodes"`
	Links                 string   `json:"links"`
	TotalCollectorCounter int      `json:"collector_counter"`
	CollectorCounterQueue int      `json:"collector_counter_queue"`
	DeadListQueueSites    []string `json:"dlq_sites"`
}

type DataBase interface {
	GetProject(id string) (*Project, error)
	GetProjectTemporaryData(id string) (*ProjectTemporaryData, error)

	CreateProject(project *Project) (string, error)
	SetProjectTemporaryData(id string, data *ProjectTemporaryData) error

	UpdateProject(project *Project) error

	DeleteProject(id string) error
	DeleteProjectTemporaryData(id string) error

	GetProjectsByOwnerId(ownerId string) ([]*ShortProject, error)

	CheckSlug(slag string) (bool, error)
	UpdateSlug(slag string, status bool) error

	CheckCollectorCounter(id string) error

	GetProjectMaxDepth(id string) (int, error)
}
