package integrationTests

import "testing"

type Message struct {
	topic     string
	partition string
	message   map[string]string
}

type SitesKafka interface {
	AddSiteToParse(link string, projectId int) error
	CheckSitesToParse(link string, projectId int) error
}

type RegistrationKafka interface {
	LoginUser(user User) error
	RegisterUser(user User) error
	IsFinished(user User) error
	GetLoginRequest() error
	GetRegisterRequest() error
	FinishedLoginOrRegister(User, error) error
}

func TestAddSiteToParse(t *testing.T) {
	err := addSiteToParse("link", 1)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestCheckSitesToParse(t *testing.T) {
	//messages := []Message{}
	//Можно по идее и не передавать offset, а читать весь топик и вычленять незавершённые задачи, но это антипаттерн
	err := CheckSitesToParse()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestLoginUser(t *testing.T) {
	err := LoginUser(user)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestRegisterUser(t *testing.T) {
	err := RegisterUser(user)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestIsFinished(t *testing.T) {
	err := IsFinished(user)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestGetLoginRequest(t *testing.T) {
	err := GetLoginRequest()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestGetRegisterRequest(t *testing.T) {
	err := GetRegisterRequest()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestFinishedLoginOrRegister(t *testing.T) {
	err := FinishedLoginOrRegister(user, err)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}
