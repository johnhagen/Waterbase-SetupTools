package WBU

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
)

//"http://localhost:8080"

var M sync.Mutex
var services map[string]*Service

var serverIP string

// const serverIP = "http://localhost:8080" //"http://192.168.50.121:9420"
const REGISTER_URL = "/waterbase/register"
const RETRIEVE_URL = "/waterbase/retrieve"
const REMOVE_URL = "/waterbase/remove"
const TRANSMITT_URL = "/waterbase/transmitt"

func Init(ServerIP string) {
	serverIP = ServerIP
	services = make(map[string]*Service)
}

func DBCheck() *map[string]*Service {
	return &services
}

// Creates a new service and returns the auth key
func CreateService(name string, owner string, adminkey string) *Service {
	defer http.DefaultClient.CloseIdleConnections()

	req := make(map[string]interface{})
	service := Service{}

	req["name"] = name
	req["owner"] = owner
	req["adminkey"] = adminkey

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(req)

	reqData, err := http.NewRequest(http.MethodPost, serverIP+REGISTER_URL+"?type=service", b)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	res, err := WebClient().Do(reqData)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	var data map[string]interface{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil
	}

	service.Name = name
	service.Owner = owner
	service.Authkey = data["auth"].(string)
	service.Collections = make(map[string]*Collection)
	M.Lock()
	services[service.Name] = &service
	M.Unlock()

	return services[service.Name]
}

func GetService(name string, auth string) *Service {
	defer http.DefaultClient.CloseIdleConnections()

	url := serverIP + RETRIEVE_URL + "?type=service"

	jsonData := make(map[string]interface{})

	jsonData["auth"] = auth
	jsonData["servicename"] = name

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(jsonData)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	req, err := http.NewRequest(http.MethodGet, url, b)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	res, err := WebClient().Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	service := Service{}

	err = json.Unmarshal(body, &service)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	service.Name = name
	service.Authkey = auth
	services[service.Name] = &service

	return services[service.Name]
}

func DeleteService(name string, auth string) bool {
	defer http.DefaultClient.CloseIdleConnections()

	//http://localhost:8080/waterbase/remove?type=service

	url := serverIP + REMOVE_URL + "?type=service"

	jsonData := make(map[string]interface{})

	jsonData["auth"] = auth
	jsonData["servicename"] = name

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(jsonData)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	req, err := http.NewRequest(http.MethodDelete, url, b)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	res, err := WebClient().Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusAccepted {
		fmt.Println("Failed to delete service")
		return false
	}

	M.Lock()
	delete(services, name)
	M.Unlock()
	return true
}

func (s *Service) CreateCollection(name string) *Collection {
	defer http.DefaultClient.CloseIdleConnections()

	if s.Name == "" {
		fmt.Println("Service doesnt exist")
		return nil
	}

	req := make(map[string]interface{})

	collection := new(Collection)

	collection.Name = name
	collection.Owner = s.Owner
	collection.Servicename = s.Name
	collection.Authkey = s.Authkey
	collection.Documents = make(map[string]*Document)

	req["name"] = name
	req["owner"] = s.Owner
	req["auth"] = s.Authkey
	req["servicename"] = s.Name

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(collection)

	res, err := WebClient().Post(serverIP+REGISTER_URL+"?type=collection", "application/json", b)
	if err != nil {
		fmt.Println("Create collection: " + err.Error())
		return nil
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusAlreadyReported {
		fmt.Println("Collection already exists")
		return nil
	} else if res.StatusCode != http.StatusAccepted {
		fmt.Println("Failed to create collection")
		return nil
	}

	s.Collections[name] = collection
	return s.Collections[name]
}

func (s *Service) GetAllCollections() []string {
	defer http.DefaultClient.CloseIdleConnections()

	url := serverIP + TRANSMITT_URL + "?type=collections"

	reqData := make(map[string]interface{})
	reqData["auth"] = s.Authkey
	reqData["servicename"] = s.Name

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(reqData)

	req, err := http.NewRequest(http.MethodGet, url, b)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	res, err := WebClient().Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer res.Body.Close()

	var colNames []string

	data, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	err = json.Unmarshal(data, &colNames)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return colNames
}

func (s *Service) GetCollection(name string) *Collection {
	defer http.DefaultClient.CloseIdleConnections()

	url := serverIP + RETRIEVE_URL + "?type=collection"

	body := make(map[string]interface{})

	body["auth"] = s.Authkey
	body["servicename"] = s.Name
	body["collectionname"] = name

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(body)
	if err != nil {
		fmt.Println("Fuck me")
		fmt.Println(err.Error())
		return nil
	}

	req, err := http.NewRequest(http.MethodGet, url, b)
	if err != nil {
		fmt.Println("Fuck me")
		fmt.Println(err.Error())
		return nil
	}

	res, err := WebClient().Do(req)
	if err != nil {
		fmt.Println("Fuck me")
		fmt.Println(err.Error())
		return nil
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Println("Could not find the collection")
		return nil
	}

	collection := new(Collection)

	data, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	err = json.Unmarshal(data, &collection)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	collection.Authkey = s.Authkey
	s.Collections[name] = collection

	return s.Collections[name]
}

func (s *Service) DeleteCollection(name string) bool {
	defer http.DefaultClient.CloseIdleConnections()

	data := make(map[string]interface{})

	data["auth"] = s.Authkey
	data["servicename"] = s.Name
	data["collectionname"] = name

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(data)

	req, err := http.NewRequest(http.MethodDelete, serverIP+REMOVE_URL+"?type=collection", b)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	res, err := WebClient().Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusAccepted {
		fmt.Println("Failed to delete collection")
		return false
	}

	delete(s.Collections, name)
	return true
}

func (c *Collection) CreateDocument(name string, content map[string]interface{}) *Document {
	defer http.DefaultClient.CloseIdleConnections()

	req := make(map[string]interface{})

	req["name"] = name
	req["owner"] = c.Owner
	req["auth"] = c.Authkey
	req["collectionname"] = c.Name
	req["servicename"] = c.Servicename
	req["content"] = content

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(req)

	reqData, err := http.NewRequest(http.MethodPost, serverIP+REGISTER_URL+"?type=document", b)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	res, err := WebClient().Do(reqData)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusAccepted {
		fmt.Println("Failed to create document")
		return nil
	}

	doc := new(Document)
	doc.Content = content
	doc.CreationDate = "temp"
	doc.LastUpdated = "temp"
	doc.Name = name
	doc.Owner = c.Owner
	doc.UpdatedBy = "temp"

	fmt.Println("Created document: " + doc.Name)
	c.Documents[name] = doc
	return c.Documents[name]
}

func (c *Collection) GetAllDocuments() []string {
	defer http.DefaultClient.CloseIdleConnections()

	url := serverIP + TRANSMITT_URL + "?type=documents"

	reqData := make(map[string]interface{})
	reqData["auth"] = c.Authkey
	reqData["servicename"] = c.Servicename
	reqData["collectionname"] = c.Name

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(reqData)

	req, err := http.NewRequest(http.MethodGet, url, b)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	res, err := WebClient().Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer res.Body.Close()

	var docNames []string

	data, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	err = json.Unmarshal(data, &docNames)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return docNames
}

func (c *Collection) GetDocument(name string) *Document {
	defer http.DefaultClient.CloseIdleConnections()

	url := serverIP + RETRIEVE_URL + "?type=document"

	reqData := make(map[string]interface{})
	reqData["auth"] = c.Authkey
	reqData["servicename"] = c.Servicename
	reqData["collectionname"] = c.Name
	reqData["documentname"] = name

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(reqData)

	req, err := http.NewRequest(http.MethodGet, url, b)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	res, err := WebClient().Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Println("Failed to get document")
		return nil
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	document := new(Document)

	err = json.Unmarshal(data, &document)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	c.Documents[name] = document
	return c.Documents[name]
}

func (c *Collection) DeleteDocument(name string) bool {
	defer http.DefaultClient.CloseIdleConnections()

	url := serverIP + REMOVE_URL + "?type=document"

	reqData := make(map[string]interface{})
	reqData["servicename"] = c.Servicename
	reqData["collectionname"] = c.Name
	reqData["documentname"] = name
	reqData["auth"] = c.Authkey

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(reqData)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	req, err := http.NewRequest(http.MethodDelete, url, b)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	res, err := WebClient().Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusAccepted {
		fmt.Println("Failed to delete document")
		return false
	}

	fmt.Println("Deleted document: " + name)
	delete(c.Documents, name)

	return true
}

func StressTest(rounds int) {
	for i := 0; i < rounds; i++ {
		CreateService(fmt.Sprintf("%v", rand.Intn(1000000)), "John", "Keks")
	}
}

func HyperStressTest(rounds int) {

	content := make(map[string]interface{})

	content["tempData"] = "HyperTesting for the win"

	for i := 0; i < rounds; i++ {
		fmt.Println("Round: " + fmt.Sprintf("%v", i))
		service := CreateService(fmt.Sprintf("%v", rand.Intn(1000000)), "John", "Keks")
		if service == nil {
			fmt.Println("Failed to create service")
			continue
		}

		collection := service.CreateCollection(fmt.Sprintf("%v", rand.Intn(1000000)))

		collection.CreateDocument(fmt.Sprintf("%v", rand.Intn(1000000)), content)
	}
}
