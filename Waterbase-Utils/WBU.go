package WBU

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

var C context.Context
var S *semaphore.Weighted
var M sync.Mutex
var services map[string]*Service
var Rclient *http.Client

var serverIP string
var creds string

const REGISTER_URL = "/waterbase/register"
const RETRIEVE_URL = "/waterbase/retrieve"
const REMOVE_URL = "/waterbase/remove"
const TRANSMITT_URL = "/waterbase/transmitt"

func Init(config SetupConfig) {

	C = context.Background()
	S = semaphore.NewWeighted(config.Threads)
	Rclient = config.Router
	serverIP = config.ServerIP
	services = make(map[string]*Service)
	creds = base64.StdEncoding.EncodeToString([]byte(config.Username + ":" + config.Password))
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
		t1 := time.Now()
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

		DeleteService(service.Name, service.Authkey)
		fmt.Printf("Amount per second: %d\n", 1000/(time.Since(t1).Milliseconds()))
	}

}
