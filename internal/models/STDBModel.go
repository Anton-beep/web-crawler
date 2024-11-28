package models

type ProjectTemporaryData struct {
	text   string
	titles string
	nodes  string
	links  string
}

type ShortTermDataBase interface {
	GetProjectTemporaryData(id string) (*ProjectTemporaryData, error)
	CreateProjectTemporaryData(id string, data *ProjectTemporaryData) error
	DeleteProjectTemporaryData(id string) error
	// CheckLink checks if the link has been processed
	CheckLink(slag string) (bool, error)
	UpdateLink(slag string, status bool) error
}
