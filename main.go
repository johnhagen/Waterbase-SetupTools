package main

import (
	"fmt"
	"math/rand"
	WBU "sandbox/Waterbase-Utils"
)

func main() {

	//WBU.Init("http://192.168.50.121:8081")
	WBU.Init("http://localhost:8080")

	serName := "TestSer1"
	authKey := "5CF04D232C957476CB074506E6B5711"

	service := WBU.CreateService(serName, "John", "Keks")

	if service == nil {
		service = WBU.GetService(serName, authKey)
	}

	if service == nil {
		fmt.Println("Need correct auth key")
		return
	}

	fmt.Println(service.Authkey)

	col := service.CreateCollection("Cheaters")
	if col == nil {
		col = service.GetCollection("Cheaters")
	}

	fmt.Println(col)

	new := make(map[string]interface{})

	new["test"] = "lalalal"

	if col.GetDocument("thisisatest") != nil {
		col.DeleteDocument("thisisatest")
		col.CreateDocument("thisisatest", new)
	}

	service.GetCollection("Cheaters").GetDocument("thisisatest")

	for i := 0; i < 100; i++ {
		service.CreateCollection(fmt.Sprintf("%d", rand.Intn(10000)))
	}

	for i := 0; i < 10; i++ {
		newService := WBU.CreateService(fmt.Sprintf("%d", rand.Intn(1000)), "John", "Keks")
		if newService != nil {
			for j := 0; j < 100; j++ {
				newCollection := newService.CreateCollection(fmt.Sprintf("%d", rand.Intn(100000)))
				if newCollection != nil {
					testing := make(map[string]interface{})
					for k := 0; k < 100; k++ {
						testing["Data1"] = fmt.Sprintf("%d", rand.Intn(1002020))
						testing["Data2"] = fmt.Sprintf("%d", rand.Intn(1002020))
						testing["Data3"] = fmt.Sprintf("%d", rand.Intn(1002020))
						testing["Data4"] = fmt.Sprintf("%d", rand.Intn(1002020))
						testing["Data5"] = fmt.Sprintf("%d", rand.Intn(1002020))
						newCollection.CreateDocument(fmt.Sprintf("%d", rand.Intn(10000)), testing)
					}
				}
			}
		}
	}
}
