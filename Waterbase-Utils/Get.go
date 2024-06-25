package WBU

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetService(name string, auth string) *Service {
	//defer http.DefaultClient.CloseIdleConnections()

	url := serverIP + RETRIEVE_URL + "?type=service"

	jsonData := make(map[string]interface{})

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

	req.Header.Add("Servicename", name)
	req.Header.Add("Auth", auth)
	req.Header.Add("Authorization", "Basic "+creds)

	res, err := Rclient.Do(req)
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

func (s *Service) GetAllCollections() []string {
	//defer http.DefaultClient.CloseIdleConnections()

	url := serverIP + TRANSMITT_URL + "?type=collections"

	reqData := make(map[string]interface{})

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(reqData)

	req, err := http.NewRequest(http.MethodGet, url, b)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	req.Header.Add("Servicename", s.Name)
	req.Header.Add("Auth", s.Authkey)
	req.Header.Add("Authorization", "Basic "+creds)

	res, err := Rclient.Do(req)
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
	//defer http.DefaultClient.CloseIdleConnections()

	url := serverIP + RETRIEVE_URL + "?type=collection"

	body := make(map[string]interface{})

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

	req.Header.Add("Auth", s.Authkey)
	req.Header.Add("Servicename", s.Name)
	req.Header.Add("Collectionname", name)
	req.Header.Add("Authorization", "Basic "+creds)

	res, err := Rclient.Do(req)
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

func (c *Collection) GetAllDocuments() []string {
	//defer http.DefaultClient.CloseIdleConnections()

	url := serverIP + TRANSMITT_URL + "?type=documents"

	reqData := make(map[string]interface{})

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(reqData)

	req, err := http.NewRequest(http.MethodGet, url, b)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	req.Header.Add("Auth", c.Authkey)
	req.Header.Add("Servicename", c.Servicename)
	req.Header.Add("Collectionname", c.Name)
	req.Header.Add("Authorization", "Basic "+creds)

	res, err := Rclient.Do(req)
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
	//defer http.DefaultClient.CloseIdleConnections()

	url := serverIP + RETRIEVE_URL + "?type=document"

	reqData := make(map[string]interface{})

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(reqData)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	req.Header.Add("Servicename", c.Servicename)
	req.Header.Add("Auth", c.Authkey)
	req.Header.Add("Collectionname", c.Name)
	req.Header.Add("Documentname", name)
	req.Header.Add("Authorization", "Basic "+creds)

	res, err := Rclient.Do(req)
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
