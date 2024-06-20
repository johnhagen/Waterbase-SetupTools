package WBU

import (
	"fmt"
	"math/rand"
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

		docName := fmt.Sprintf("%v", rand.Intn(1000000))

		collection.CreateDocument(docName, content)

		collection.GetDocument(docName)
	}
}
