package models

type NotFound struct {
}

func (NotFound) Error() string {
	return "requested element not found"
}

var DataBaseNotFound NotFound

type WrongUUID struct {
}

func (WrongUUID) Error() string {
	return "invalid uuid"
}

var DataBaseWrongID WrongUUID
