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
	JobTitle                string `json:"jobTitle"`
	HireDate                string `json:"hireDate"`
	OriginalHireDate        string `json:"originalHireDate"`
	EmployeeNumber          string `json:"employeeNumber"`
	Location                string `json:"location"`
	PreferredName           string `json:"preferredName"`
	DisplayName             string `json:"displayName"`
	ReportsTo               string `json:"reportsTo"`

	// CustomFields holds any non-standard field returned by the report
	// (i.e. fields requested via the --custom-fields config option).
	CustomFields map[string]string `json:"-"`
}

// StandardFields are the BambooHR field names the connector always requests and
// maps to first-class profile attributes. They are excluded from CustomFields.
var StandardFields = []string{
	"firstName",
	"lastName",
	"supervisor",
	"supervisorEId",
	"supervisorId",
	"supervisorEmail",
	"workEmail",
	"status",
	"department",
	"division",
	"employmentHistoryStatus",
	"terminationDate",
	"jobTitle",
	"hireDate",
	"originalHireDate",
	"employeeNumber",
	"location",
	"preferredName",
	"displayName",
	"reportsTo",
}

// standardFieldSet is StandardFields plus "id" for O(1) lookups during unmarshal.
var standardFieldSet = func() map[string]bool {
	m := map[string]bool{"id": true}
	for _, f := range StandardFields {
		m[f] = true
	}
	return m
}()

// UnmarshalJSON decodes the standard fields normally and captures every other
// returned field into CustomFields as a string.
func (u *User) UnmarshalJSON(data []byte) error {
	type alias User
	aux := alias{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	*u = User(aux)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	u.CustomFields = make(map[string]string)
	for key, val := range raw {
		if standardFieldSet[key] {
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
