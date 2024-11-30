package models

type SitesKafka interface {
	AddSiteToParse(link string, projectId int) error
	CheckSitesToParse(link string, projectId int) error
}

//type RegistrationKafka interface {
//	LoginUser(user User) error
//	RegisterUser(user User) error
//	IsFinished(user User) error
//
//	GetLoginRequest() error
//	GetRegisterRequest() error
//	FinishedLoginOrRegister(User, error) error
//}
