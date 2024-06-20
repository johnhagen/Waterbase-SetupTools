package WBU

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

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
