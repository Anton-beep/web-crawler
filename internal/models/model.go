package models

type Project struct {
	ID       string // uuid
	OwnerID  string // uuid
	WebGraph string
	DlqSites []string
}

func (p Project) Equals(other Project) bool {
	if !(p.ID == other.ID &&
		p.OwnerID == other.OwnerID &&
		p.WebGraph == other.WebGraph) {
		return false
	}
	if len(p.DlqSites) != len(other.DlqSites) {
		return false
	}
	for site1, _ := range p.DlqSites {
		seen := false
		for site2, _ := range other.DlqSites {
			if site1 == site2 {
				seen = true
				break
			}
		}
		if !seen {
			return false
		}
	}
	return true
}

type ProjectTemporaryData struct {
	Text   string
	Titles string
	Nodes  string
	Links  string
}

type DataBase interface {
	GetProject(id string) (*Project, error)
	GetProjectTemporaryData(id string) (*ProjectTemporaryData, error)

	CreateProject(project Project) (string, error)
	CreateProjectTemporaryData(id string, data *ProjectTemporaryData) error

	DeleteProject(id string) error
	DeleteProjectTemporaryData(id string) error

	GetProjectsByOwnerId(ownerId string) ([]*Project, error)

	CheckLink(slag string) (bool, error)
	UpdateLink(slag string, status bool) error
}
