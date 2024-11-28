package repository

import (
	"errors"
	"experiments/internal/models"
)

type Cache struct {
	mp  map[string]*models.ProjectTemporaryData
	url map[string]bool
}

func (c Cache) GetProjectTemporaryData(id string) (*models.ProjectTemporaryData, error) {
	//TODO implement me
	//panic("implement me")

	res, ok := c.mp[id]
	if !ok {
		return nil, errors.New("ProjectTemporaryData with id " + id + " doesn't exist")
	}
	return res, nil
}

func (c Cache) CreateProjectTemporaryData(id string, data *models.ProjectTemporaryData) error {
	//TODO implement me
	//panic("implement me")
	c.mp[id] = data
	return nil
}

func (c Cache) DeleteProjectTemporaryData(id string) error {
	//TODO implement me
	//panic("implement me")

	_, ok := c.mp[id]
	if !ok {
		return errors.New("ProjectTemporaryData with id " + id + " doesn't exist")
	}
	delete(c.mp, id)
	return nil
}

func (c Cache) CheckLink(slag string) (bool, error) {
	//TODO implement me
	//panic("implement me")

	v, ok := c.url[slag]
	if !ok {
		c.url[slag] = false
		return false, nil
	}
	return v, nil
}

func (c Cache) UpdateLink(slag string, status bool) error {
	//TODO implement me
	//panic("implement me")

	c.url[slag] = status
	return nil
}

func NewCache() models.ShortTermDataBase {
	return Cache{
		mp:  make(map[string]*models.ProjectTemporaryData),
		url: make(map[string]bool),
	}
}
