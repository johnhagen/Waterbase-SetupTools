package main

import (
	"fmt"
	WBU "sandbox/Waterbase-Utils"
)

func main() {

	WBU.Init()

	//content := make(map[string]interface{})

	service := WBU.GetService("testservice", "A75635B24B50CA0F358F6469F0E0AF3")
	if service == nil {
		service = WBU.CreateService("testservice", "john", "Keks")
	}

	//collection := service.CreateCollection("Kekekeke")

	fmt.Println(service.GetAllCollections())

	//WBU.HyperStressTest(2000)

	//collection.DeleteDocument("Testing")

	//WBU.DeleteService("107594", "67D8A9F78068320D780D973C85B04FA")

	//service.DeleteCollection("TestCollection")

	//WBU.DeleteService(service.Name, service.Authkey)

	/*
		var service WBU.Service
		serviceName := "theBibble"
		authKey := "C0DC633F345890732107BFC5BED4AE7C9739E9CCE8FF478191E6AE09125E89EEFA7C14445C44D7D28CBA1521157EAC5335BD235672B9E0C0EE9E1BB05207A4A18776A05C30297796BFB6178E79681101D9C1A07456EC1781C41179138A3ECF43A02BCE35E228164B10578DE1F0ADF680E3A96F7B9AA8CF2E807AF93198886BB98164DF6AA77492963D724AB44964591D078C8BD3115D04B475B7552772AB0E37275D828203B8CC85AF72536F93885D221D44227F5F5A04BDDFE52C9C2D8803FF51A1CBB86516528CA314E7329943ACCB5906D34FF8E10355B0DA7909252B7EC59CD33D680922636F572EFCECF31ED37B9AFCF42985EB0D30DC841BC377BB93105870A1E3680F012B618A377E6F67FD223980926E5BD0E8418F55DCE5A297B74B5EB28C3E3EA93D4FA74B2D09F482A2B710C5286F3B996E98BB515E140959EC33C8A7F7FBC214F95205CF5FC6946B9BAE4570CB9AB565DE4F608D85424823CB03C7B945DD6C07C0C3CC4BBECAA30239533AB90C870EE90225C45583123B912897081A17D585B9A00A04D37B51E971F5A442F987361D64B13FECE76A5BEE93218A6788724C17CDCE3BAB82CD4FB4FC4828F321E46110710BB14BD9F6F7F11DC30BCBDFC30968F595308BE0C905FD50966C4B3D5532566253B2DC6A7D7B23F239415DBBC01E1F3B5B7499C7C7ABB99F7D5377F53F8"

		Foundbf5, found := WBU.GetService(serviceName, authKey)

		if found {
			fmt.Println("Found old service")
			service = Foundbf5
		} else {
			fmt.Println("New Service created")
			Newbf5, _ := WBU.CreateService(serviceName, "John", "Keks")
			service = Newbf5
		}

		theBibble := service.CreateCollection("BibleVerses")

		for i := 0; i < 1; i++ {

			res, err := http.Get("https://bible-api.com/?random=verse")
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			result := make(map[string]interface{})

			err = json.Unmarshal(body, &result)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			err = theBibble.CreateDocument(result["reference"].(string), result)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	*/
}
