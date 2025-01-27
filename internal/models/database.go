package models

type Project struct {
	ID               string   `db:"id" json:"id"`
	OwnerID          string   `db:"owner_id" json:"owner_id"`
	Name             string   `db:"name" json:"name"`
	StartUrl         string   `db:"start_url" json:"start_url"`
	Processing       bool     `db:"processing" json:"processing"`
	WebGraph         string   `db:"web_graph" json:"web_graph"`
	DlqSites         []string `db:"dlq_sites" json:"dlq_sites"`
	MaxDepth         int      `db:"max_depth" json:"max_depth"`
	MaxNumberOfLinks int      `db:"max_number_of_links" json:"max_number_of_links"`
	KeyWords         string   `db:"key_words" json:"key_words"`
	MainIdeas        string   `db:"main_ideas" json:"main_ideas"`
}

type ShortProject struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AnalyserTask struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type ProjectTemporaryData struct {
	Text                  []string `json:"text"`
	Titles                []string `json:"titles"`
	Nodes                 string   `json:"nodes"`
	Links                 string   `json:"links"`
	TotalCollectorCounter int      `json:"collector_counter"`
	CollectorCounterQueue int      `json:"collector_counter_queue"`
	DeadListQueueSites    []string `json:"dlq_sites"`
}

type User struct {
	ID       string `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
}

// DataBase is an interface that contains methods for working with the database
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

	Push2Queue(key string, value interface{}) error
	PopFromQueue(key string) (string, error)

	AddAnalyserTask(projectId, typeOfAnalysis string) error
	GetAnalyserTask() (AnalyserTask, error)

	GetUserByUsername(username string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	AddUser(user *User) (string, error)
	UpdateUser(user *User) error
}
