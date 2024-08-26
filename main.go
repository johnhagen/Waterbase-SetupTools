package main

import (
	"log"
	"net/http"
	"os"
	"runtime"
	WBU "sandbox/Waterbase-Utils"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	var Config WBU.SetupConfig

	Config.ServerIP = os.Getenv("SERVERIP")
	Config.Username = os.Getenv("BASIC_AUTH")
	Config.Password = os.Getenv("BASIC_AUTH_PASS")
	Config.Router = http.DefaultClient
	Config.Threads = int64(runtime.NumCPU())

	WBU.Init(Config)

	WBU.HyperStressTest(300)

}
