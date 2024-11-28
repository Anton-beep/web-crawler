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

type LongTermDataBase interface {
	// GetProject gets a project by id
	GetProject(id string) (*Project, error)
	// GetDashboard gets a dashboard by id
	GetDashboard(id string) (*Dashboard, error)

	CreateProject(project Project) (string, error)
	CreateDashboard(dashboard Dashboard) (string, error)

	DeleteProject(id string) error
	DeleteDashBoard(id string) error

	GetProjectsByOwnerId(ownerId string) ([]*Project, error)
}
