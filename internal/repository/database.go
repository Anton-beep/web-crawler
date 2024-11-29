package repository

import (
	"web-crauler/internal/models"
)

type DataBase struct{}

func (d DataBase) GetProject(id string) (*models.Project, error) {
	//TODO implement me
	panic("implement me")
}

func (d DataBase) CreateProject(project models.Project) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (d DataBase) DeleteProject(id string) error {
	//TODO implement me
	panic("implement me")
}

func (d DataBase) GetProjectsByOwnerId(ownerId string) ([]*models.Project, error) {
	//TODO implement me
	panic("implement me")
}

func (d DataBase) GetProjectTemporaryData(id string) (*models.ProjectTemporaryData, error) {
	//TODO implement me
	panic("implement me")
}

func (d DataBase) CreateProjectTemporaryData(id string, data *models.ProjectTemporaryData) error {
	//TODO implement me
	panic("implement me")
}

func (d DataBase) DeleteProjectTemporaryData(id string) error {
	//TODO implement me
	panic("implement me")
}

func (d DataBase) CheckLink(slag string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (d DataBase) UpdateLink(slag string, status bool) error {
	//TODO implement me
	panic("implement me")
}

func NewDB() models.DataBase {
	return DataBase{}
}
