package models

type NotFound struct {
}

// Error NotFound is an error that is returned when the requested element is not found
func (NotFound) Error() string {
	return "requested element not found"
}

var DataBaseNotFound NotFound

type WrongUUID struct {
}

// Error WrongUUID is an error that is returned when the uuid is invalid
func (WrongUUID) Error() string {
	return "invalid uuid"
}

var DataBaseWrongID WrongUUID

type QueueIsEmpty struct{}

func (QueueIsEmpty) Error() string {
	return "queue is empty"
}

var DataBaseQueueIsEmpty QueueIsEmpty
