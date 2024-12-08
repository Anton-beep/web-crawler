package models

type SitesKafka interface {
	AddSiteToParse(link string, projectId int) error
	CheckSitesToParse(link string, projectId int) error
}
