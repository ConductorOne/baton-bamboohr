package client

type User struct {
	Id              string `json:"id"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Supervisor      string `json:"supervisor"`
	SupervisorEId   string `json:"supervisorEId"`
	SupervisorId    string `json:"supervisorId"`
	SupervisorEmail string `json:"supervisorEmail"`
	Email           string `json:"workEmail"`
	Status          string `json:"status"`
}

type Fields struct {
	Id   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

type ReqFields struct {
	Title  string   `json:"title"`
	Fields []string `json:"fields"`
}

type ReportUserResults struct {
	Title  string   `json:"title"`
	Fields []Fields `json:"fields"`
	Users  []*User  `json:"employees"`
}
