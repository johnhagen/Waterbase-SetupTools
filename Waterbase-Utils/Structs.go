package WBU

type Service struct {
	Authkey     string
	Name        string `json:"name"`
	Owner       string `json:"owner"`
	Collections map[string]*Collection
}

type Collection struct {
	Name        string `json:"name"`
	Owner       string `json:"owner"`
	Servicename string `json:"servicename"`
	Authkey     string `json:"auth"`
	LastUpdated string
	Documents   map[string]*Document
}

type Document struct {
	UpdatedBy    string
	Name         string `json:"name"`
	Owner        string `json:"owner"`
	CreationDate string
	LastUpdated  string
	Content      interface{} `json:"content"`
}
