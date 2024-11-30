package integrationTests

import "testing"

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
