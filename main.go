package main

import (
	"fmt"
	"net/http"
	"os"
	WBU "sandbox/Waterbase-Utils"
)

func main() {

	authKey := "F0850A2DDFA9CA715AE135F8CBAD234"
	serviceName := "FileServer"

	//WBU.Init("https://waterbase.hagen.fun", "john", "Hagnought99")
	//WBU.Init("http://192.168.50.121:8080")
	WBU.Init("http://localhost:8080", "", "", http.DefaultClient)

	FileService := WBU.CreateService(serviceName, "John", "Keks")

	if FileService == nil {
		FileService = WBU.GetService(serviceName, authKey)
		if FileService == nil {
			fmt.Println("Provide the service authkey")
			return
		}
	}

	fmt.Println(FileService.Authkey)

	FileService.DeleteCollection("JohnFiles")

	JohnCol := FileService.CreateCollection("JohnFiles")

	if JohnCol == nil {
		JohnCol = FileService.GetCollection("JohnFiles")
		if JohnCol == nil {
			fmt.Println("Failed to grab collection")
			return
		}
	}

	files, err := os.ReadDir("./Files")
	if err != nil {
		fmt.Println(err)
		return
	}

	//kek := make(map[string]interface{})
	//kek["Hei"] = "This is data"
	//FileService.DeleteCollection("Temp")
	//tempCol := FileService.CreateCollection("Temp")

	//for i := 0; i < 10000; i++ {
	//	tempCol.CreateDocument(fmt.Sprintf("%d", rand.Intn(1000000)), kek)
	//}

	content := make(map[string]interface{})

	for _, h := range files {

		file, err := os.ReadFile("./Files/" + h.Name())
		if err != nil {
			fmt.Println(err)
			return
		}

		content["Data"] = string(file)

		//t1 := time.Now()

		go JohnCol.CreateDocument(h.Name(), content)

		//fileInfo, _ := h.Info()

		//fileSizeMB := float64(fileInfo.Size()) / (1024 * 1024)

		//timeSince := time.Since(t1).Seconds()

		//fmt.Printf("Size: %.5fMB Time: %.3fsec Average: %.3fMB/s\n", fileSizeMB, time.Since(t1).Seconds(), fileSizeMB/timeSince)
	}

	//WBU.DeleteService(serviceName, authKey)

	/*
		err = os.WriteFile("./Files/newFile.txt", byte(temp["Data"].(string)), 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
	*/

}
