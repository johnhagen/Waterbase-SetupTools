package main

import (
	"fmt"
	"math/rand"
	WBU "sandbox/Waterbase-Utils"
)

func main() {

	WBU.Init("https://waterbase.hagen.fun", "john", "Hagnought99")
	//WBU.Init("http://192.168.50.121:8080")
	//WBU.Init("http://localhost:8080")

	for i := 0; i < 10; i++ {
		newService := WBU.CreateService(fmt.Sprintf("%d", rand.Intn(1000)), "John", "Keks")
		if newService != nil {
			for j := 0; j < 10; j++ {
				newCollection := newService.CreateCollection(fmt.Sprintf("%d", rand.Intn(100000)))
				if newCollection != nil {
					testing := make(map[string]interface{})
					for k := 0; k < 10; k++ {
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
