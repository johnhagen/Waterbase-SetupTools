package WBU

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Creates a new service and returns the auth key
func CreateService(name string, owner string, adminkey string) *Service {

	S.Acquire(C, 1)

	reqData := make(map[string]interface{})
	service := Service{}

	reqData["servicename"] = name
	reqData["owner"] = owner
	reqData["adminkey"] = adminkey

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(reqData)

	req, err := http.NewRequest(http.MethodPost, serverIP+REGISTER_URL+"?type=service", b)
	if err != nil {
		fmt.Println(err.Error())
		S.Release(1)
		return nil
	}

	req.Header.Add("Authorization", "Basic "+creds)

	res, err := Rclient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		S.Release(1)
		return nil
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		S.Release(1)
		return nil
	}

	var data map[string]interface{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		S.Release(1)
		return nil
	}

	service.Name = name
	service.Owner = owner
	service.Authkey = data["auth"].(string)
	service.Collections = make(map[string]*Collection)
	M.Lock()
	services[service.Name] = &service
	M.Unlock()
	S.Release(1)
	return &service
}

func (s *Service) CreateCollection(name string) *Collection {
	S.Acquire(C, 1)

	if s.Name == "" {
		fmt.Println("Service doesnt exist")
		S.Release(1)
		return nil
	}

	reqData := make(map[string]interface{})

	collection := new(Collection)

	collection.Name = name
	collection.Owner = s.Owner
	collection.Servicename = s.Name
	collection.Authkey = s.Authkey
	collection.Documents = make(map[string]*Document)

	reqData["collectionname"] = name
	reqData["owner"] = s.Owner
	reqData["auth"] = s.Authkey
	reqData["servicename"] = s.Name

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(reqData)

	req, err := http.NewRequest(http.MethodPost, serverIP+REGISTER_URL+"?type=collection", b)
	if err != nil {
		fmt.Println(err.Error())
		S.Release(1)
		return nil
	}

	req.Header.Add("Authorization", "Basic "+creds)

	res, err := Rclient.Do(req)
	if err != nil {
		fmt.Println("Create collection: " + err.Error())
		S.Release(1)
		return nil
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusAlreadyReported {
		fmt.Println("Collection already exists")
		S.Release(1)
		return nil
	} else if res.StatusCode != http.StatusAccepted {
		fmt.Println("Failed to create collection")
		S.Release(1)
		return nil
	}

	s.Mutex.Lock()
	s.Collections[name] = collection
	s.Mutex.Unlock()
	S.Release(1)
	return collection
}

func (c *Collection) CreateDocument(name string, content interface{}) *Document {
	S.Acquire(C, 1)

	reqData := make(map[string]interface{})

	reqData["documentname"] = name
	reqData["owner"] = c.Owner
	reqData["auth"] = c.Authkey
	reqData["collectionname"] = c.Name
	reqData["servicename"] = c.Servicename
	reqData["content"] = content

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(reqData)

	req, err := http.NewRequest(http.MethodPost, serverIP+REGISTER_URL+"?type=document", b)
	if err != nil {
		fmt.Println(err.Error())
		S.Release(1)
		return nil
	}

	req.Header.Add("Authorization", "Basic "+creds)

	res, err := Rclient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		S.Release(1)
		return nil
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusAccepted {
		fmt.Println("Failed to create document")
		S.Release(1)
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
	c.Mutex.Lock()
	c.Documents[name] = doc
	c.Mutex.Unlock()
	S.Release(1)
	return doc
}
