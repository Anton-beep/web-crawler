package models

type NotFound struct {
}

func (NotFound) Error() string {
	return "requested element not found"
}

var DataBaseNotFound NotFound
