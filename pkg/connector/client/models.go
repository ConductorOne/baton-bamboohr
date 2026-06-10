package client

import "encoding/json"

type User struct {
	Id                      string `json:"id"`
	FirstName               string `json:"firstName"`
	LastName                string `json:"lastName"`
	Supervisor              string `json:"supervisor"`
	SupervisorEId           string `json:"supervisorEId"`
	SupervisorId            string `json:"supervisorId"`
	SupervisorEmail         string `json:"supervisorEmail"`
	Email                   string `json:"workEmail"`
	Status                  string `json:"status"`
	Division                string `json:"division"`
	Department              string `json:"department"`
	EmploymentHistoryStatus string `json:"employmentHistoryStatus"`
	TerminationDate         string `json:"terminationDate"`
	CustomFields            map[string]string
}

var knownFields = map[string]bool{
	"id": true, "firstName": true, "lastName": true,
	"supervisor": true, "supervisorEId": true, "supervisorId": true,
	"supervisorEmail": true, "workEmail": true, "status": true,
	"division": true, "department": true,
	"employmentHistoryStatus": true, "terminationDate": true,
}

func (u *User) UnmarshalJSON(data []byte) error {
	type Alias User
	aux := &Alias{}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	*u = User(*aux)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	u.CustomFields = make(map[string]string)
	for key, val := range raw {
		if knownFields[key] {
			continue
		}
		var s string
		if json.Unmarshal(val, &s) == nil {
			u.CustomFields[key] = s
		} else {
			u.CustomFields[key] = string(val)
		}
	}

	return nil
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
