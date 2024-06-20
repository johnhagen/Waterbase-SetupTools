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
	//defer http.DefaultClient.CloseIdleConnections()

	reqData := make(map[string]interface{})
	service := Service{}

	reqData["name"] = name
	reqData["owner"] = owner
	reqData["adminkey"] = adminkey

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(reqData)

	req, err := http.NewRequest(http.MethodPost, serverIP+REGISTER_URL+"?type=service", b)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	req.Header.Add("Authorization", "Basic "+creds)

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

func (s *Service) CreateCollection(name string) *Collection {
	//defer http.DefaultClient.CloseIdleConnections()

	if s.Name == "" {
		fmt.Println("Service doesnt exist")
		return nil
	}

	reqData := make(map[string]interface{})

	collection := new(Collection)

	collection.Name = name
	collection.Owner = s.Owner
	collection.Servicename = s.Name
	collection.Authkey = s.Authkey
	collection.Documents = make(map[string]*Document)

	reqData["name"] = name
	reqData["owner"] = s.Owner
	reqData["auth"] = s.Authkey
	reqData["servicename"] = s.Name

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(collection)

	req, err := http.NewRequest(http.MethodPost, serverIP+REGISTER_URL+"?type=collection", b)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	req.Header.Add("Authorization", "Basic "+creds)

	res, err := WebClient().Do(req)
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

func (c *Collection) CreateDocument(name string, content map[string]interface{}) *Document {
	//defer http.DefaultClient.CloseIdleConnections()

	reqData := make(map[string]interface{})

	reqData["name"] = name
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
		return nil
	}

	req.Header.Add("Authorization", "Basic "+creds)

	res, err := WebClient().Do(req)
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
