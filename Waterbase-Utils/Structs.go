package WBU

import (
	"net/http"
	"sync"
)

type DB struct {
	Services map[string]*Service
	Mutex    sync.Mutex
}

type Service struct {
	Authkey     string
	Name        string `json:"name"`
	Owner       string `json:"owner"`
	Collections map[string]*Collection
	Mutex       sync.Mutex
}

type Collection struct {
	Name        string `json:"name"`
	Owner       string `json:"owner"`
	Servicename string `json:"servicename"`
	Authkey     string `json:"auth"`
	LastUpdated string
	Documents   map[string]*Document
	Mutex       sync.Mutex
}

type Document struct {
	UpdatedBy    string      `json:"updatedBy"`
	Name         string      `json:"name"`
	Owner        string      `json:"owner"`
	CreationDate string      `json:"creationDate"`
	LastUpdated  string      `json:"lastUpdated"`
	Content      interface{} `json:"content"`
}

type SetupConfig struct {
	ServerIP string
	Username string
	Password string
	Router   *http.Client
	Threads  int64
}
