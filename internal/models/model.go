package models

type Project struct {
	ID      string // uuid
	OwnerID string // uuid
}

type Dashboard struct {
	ID        string // uuid
	ProjectID string // uuid
	WebGraph  string
	DlqSites  []string
}

type ProjectTemporaryData struct {
	text   string
	titles string
	nodes  string
	links  string
}

type LongTermDataBase interface {
	GetProject(id string) (*Project, error)
	GetDashboard(id string) (*Dashboard, error)
	GetProjectTemporaryData(id string) (*ProjectTemporaryData, error)

	CreateProject(project Project) (string, error)
	CreateDashboard(dashboard Dashboard) (string, error)
	CreateProjectTemporaryData(id string, data *ProjectTemporaryData) error

	DeleteProject(id string) error
	DeleteDashBoard(id string) error
	DeleteProjectTemporaryData(id string) error

	GetProjectsByOwnerId(ownerId string) ([]*Project, error)

	CheckLink(slag string) (bool, error)
	UpdateLink(slag string, status bool) error
}
