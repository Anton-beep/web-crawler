package broker

type Message struct {
	Link      string `json:"link"`
	ProjectId int    `json:"projectId"`
}

type SitesKafka struct {
	Consumer interface{}
	Producer interface{}
	Topic    string
}

func (s *SitesKafka) AddSiteToParse(link string, projectId int) error {
	//TODO implement me
	panic("func is not implement")
}

func (s *SitesKafka) CheckSitesToParse(link string, projectId int) error {
	//TODO implement me
	panic("func is not implement")
}
