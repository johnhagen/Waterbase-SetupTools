package WBU

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func DeleteService(name string, auth string) bool {
	S.Acquire(C, 1)

	url := serverIP + REMOVE_URL + "?type=service"

	jsonData := make(map[string]interface{})

	jsonData["auth"] = auth
	jsonData["servicename"] = name

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(jsonData)
	if err != nil {
		fmt.Println(err.Error())
		S.Release(1)
		return false
	}

	req, err := http.NewRequest(http.MethodDelete, url, b)
	if err != nil {
		fmt.Println(err.Error())
		S.Release(1)
		return false
	}

	req.Header.Add("Authorization", "Basic "+creds)

	res, err := Rclient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		S.Release(1)
		return false
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusAccepted {
		fmt.Println("Failed to delete service")
		S.Release(1)
		return false
	}

	M.Lock()
	delete(services, name)
	M.Unlock()
	S.Release(1)
	return true
}

func (s *Service) DeleteCollection(name string) bool {
	S.Acquire(C, 1)

	data := make(map[string]interface{})

	data["auth"] = s.Authkey
	data["servicename"] = s.Name
	data["collectionname"] = name

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(data)

	req, err := http.NewRequest(http.MethodDelete, serverIP+REMOVE_URL+"?type=collection", b)
	if err != nil {
		fmt.Println(err.Error())
		S.Release(1)
		return false
	}

	req.Header.Add("Authorization", "Basic "+creds)

	res, err := Rclient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		S.Release(1)
		return false
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusAccepted {
		fmt.Println("Failed to delete collection")
		S.Release(1)
		return false
	}

	s.Mutex.Lock()
	delete(s.Collections, name)
	s.Mutex.Unlock()
	S.Release(1)
	return true
}

func (c *Collection) DeleteDocument(name string) bool {

	S.Acquire(C, 1)

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
		S.Release(1)
		return false
	}

	req, err := http.NewRequest(http.MethodDelete, url, b)
	if err != nil {
		fmt.Println(err.Error())
		S.Release(1)
		return false
	}

	req.Header.Add("Authorization", "Basic "+creds)

	res, err := Rclient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		S.Release(1)
		return false
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusAccepted {
		fmt.Println("Failed to delete document")
		S.Release(1)
		return false
	}

	c.Mutex.Lock()
	delete(c.Documents, name)
	c.Mutex.Unlock()
	fmt.Println("Deleted document: " + name)
	S.Release(1)
	return true
}
