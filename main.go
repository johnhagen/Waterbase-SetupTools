package main

import (
	"fmt"
	"math/rand"
	WBU "sandbox/Waterbase-Utils"
)

func main() {

	//WBU.Init("http://192.168.50.121:8080")
	WBU.Init("http://localhost:8080")

	//var bf5 *WBU.Service
	servicename := "Test"
	key := "BB53476B03CC0FF0043E657FC268C91"

	bf5 := WBU.GetService(servicename, key)
	if bf5 == nil {
		bf5 = WBU.CreateService("Test", "John", "Keks")
	}

	if bf5 == nil {
		fmt.Println("Fill in service key")
		return
	}

	Cheaters := bf5.CreateCollection("Cheaters")

	data := make(map[string]interface{})
	data["info"] = "test"

	for i := 0; i < 10000; i++ {
		fmt.Println(i)
		Cheaters.CreateDocument(fmt.Sprintf("%v", rand.Intn(100000)), data)
	}

	for _, h := range Cheaters.GetAllDocuments() {
		Cheaters.DeleteDocument(h)
	}

	bf5.DeleteCollection("Cheaters")

	WBU.DeleteService(bf5.Name, bf5.Authkey)

	fmt.Println(*WBU.DBCheck())
}
