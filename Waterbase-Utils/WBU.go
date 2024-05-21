package WBU

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

//"http://localhost:8080"

const serverIP = "http://localhost:8080" //"http://192.168.50.121:9420"
const REGISTER_URL = "/waterbase/register"
const RETRIEVE_URL = "/waterbase/retrieve"
const DELETE_URL = "/waterbase/remove"

// Creates a new service and returns the auth key
func CreateService(name string, owner string, adminkey string) (Service, bool) {

	req := make(map[string]interface{})
	service := Service{}

	req["name"] = name
	req["owner"] = owner
	req["adminkey"] = adminkey

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(req)

	res, err := http.Post(serverIP+REGISTER_URL+"?type=service", "application-json", b)
	if err != nil {
		fmt.Println(err.Error())
		return Service{}, false
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		return Service{}, false
	}

	var data map[string]interface{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return Service{}, false
	}

	service.Name = name
	service.Owner = owner
	service.Authkey = data["auth"].(string)

	return service, false
}

func GetService(name string, auth string) (Service, bool) {

	url := serverIP + RETRIEVE_URL + "?service=" + name

	jsonData := make(map[string]interface{})

	jsonData["auth"] = auth
	jsonData["servicename"] = name

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(jsonData)
	if err != nil {
		fmt.Println(err.Error())
		return Service{}, false
	}

	req, err := http.NewRequest(http.MethodGet, url, b)
	if err != nil {
		fmt.Println(err.Error())
		return Service{}, false
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return Service{}, false
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		return Service{}, false
	}

	service := Service{}

	//data := make(map[string]interface{})

	err = json.Unmarshal(body, &service)
	if err != nil {
		fmt.Println(err.Error())
		return Service{}, false
	}

	service.Name = name
	service.Authkey = auth

	return service, true
}

func DeleteService(name string, auth string) bool {

	//http://localhost:8080/waterbase/remove?name=Sandbox&type=service

	url := serverIP + DELETE_URL + "?type=service&name=" + name

	fmt.Println(url)

	jsonData := make(map[string]interface{})

	jsonData["auth"] = auth

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

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	if res.StatusCode != http.StatusAccepted {
		fmt.Println("Failed to delete service")
		return false
	}

	return true
}

func (s *Service) CreateCollection(name string) *Collection {

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

	req["name"] = name
	req["owner"] = s.Owner
	req["auth"] = s.Authkey
	req["servicename"] = s.Name

	/*
		fmt.Println("Collection creation with JSON:")
		fmt.Println(name)
		fmt.Println(s.Owner)
		fmt.Println(s.Authkey)
		fmt.Println(s.Name)
	*/

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(collection)

	res, err := http.DefaultClient.Post(serverIP+REGISTER_URL+"?type=collection", "application/json", b)
	if err != nil {
		fmt.Println("Create collection: " + err.Error())
		return nil
	}
	defer res.Body.Close()
	return collection
}

func (s *Service) GetCollection(name string) *Collection {

	url := serverIP + RETRIEVE_URL + "?service=" + s.Name + "&collection=" + name

	body := make(map[string]interface{})

	body["auth"] = s.Authkey

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

	res, err := http.DefaultClient.Do(req)
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

	s.Collections[name] = collection

	return collection
}

func (c *Collection) CreateDocument(name string, content map[string]interface{}) error {

	req := make(map[string]interface{})

	req["name"] = name
	req["owner"] = c.Owner
	req["auth"] = c.Authkey
	req["collectionname"] = c.Name
	req["servicename"] = c.Servicename
	req["content"] = content

	//collection := Collection{}

	fmt.Println("Creating document with data:")
	fmt.Println(req)

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(req)

	res, err := http.Post(serverIP+REGISTER_URL+"?type=document", "application/json", b)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer res.Body.Close()
	return nil
}

func StressTest(rounds int) {
	for i := 0; i < rounds; i++ {
		CreateService(fmt.Sprintf("%v", rand.Intn(1000000)), "John", "Keks")
	}
}
