package repository

import (
	"errors"
	"experiments/internal/models"
	"github.com/google/uuid"
)

type DataBase struct {
	proj map[string]*models.Project
	dash map[string]*models.Dashboard
}

func (d DataBase) GetProject(id string) (*models.Project, error) {
	//TODO implement me
	//panic("implement me")

	res, ok := d.proj[id]
	if !ok {
		return nil, errors.New("Project with id " + id + " doesn't exist")
	}
	return res, nil
}

func (d DataBase) GetDashboard(id string) (*models.Dashboard, error) {
	//TODO implement me
	//panic("implement me")

	res, ok := d.dash[id]
	if !ok {
		return nil, errors.New("Dashboard with id " + id + " doesn't exist")
	}
	return res, nil
}

func (d DataBase) CreateProject(project models.Project) (string, error) {
	//TODO implement me
	//panic("implement me")

	project.ID = uuid.New().String()
	d.proj[project.ID] = &project
	return project.ID, nil
}

func (d DataBase) CreateDashboard(dashboard models.Dashboard) (string, error) {
	//TODO implement me
	//panic("implement me")

	dashboard.ID = uuid.New().String()
	d.dash[dashboard.ID] = &dashboard
	return dashboard.ID, nil
}

func (d DataBase) DeleteProject(id string) error {
	//TODO implement me
	//panic("implement me")

	_, ok := d.proj[id]
	if !ok {
		return errors.New("Project with id " + id + " doesn't exist")
	}
	delete(d.proj, id)
	return nil
}

func (d DataBase) DeleteDashBoard(id string) error {
	//TODO implement me
	//panic("implement me")

	_, ok := d.dash[id]
	if !ok {
		return errors.New("Dashboard with id " + id + " doesn't exist")
	}
	delete(d.dash, id)
	return nil
}

func (d DataBase) GetProjectsByOwnerId(ownerId string) ([]*models.Project, error) {
	//TODO implement me
	//panic("implement me")

	var res []*models.Project
	for _, v := range d.proj {
		if v.OwnerID == ownerId {
			res = append(res, v)
		}
	}
	return res, nil
}

func NewDB() models.LongTermDataBase {
	return DataBase{
		proj: make(map[string]*models.Project),
		dash: make(map[string]*models.Dashboard),
	}
}
