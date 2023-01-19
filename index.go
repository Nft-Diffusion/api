package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	http.HandleFunc("/retrieveImage", imageHandler)
	http.HandleFunc("/genkeys", genKeys)
	port := ":3000"
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Running API on port " + port)
}

func genKeys(w http.ResponseWriter, r *http.Request) {
	//todo
	// Error handlers
	ctx := context.Background()
	tokenId := 0
	for tokenId <= 100 {
		wg.Add(1)
		fmt.Println("running!")
		go KeyGen(ctx, "token"+strconv.Itoa(tokenId))
		tokenId++
	}
	wg.Wait()
	seedMetadata()
	io.WriteString(w, "Keys have been successfully generated!")
}
func seedMetadata() string {
	tokenId := 0
	for tokenId <= 100 {
		wg.Add(1)
		go publishIpnsAndWrite(strconv.Itoa(tokenId), "QmVMWN461ChmwQt1q4HuQ8zdenbMttAhFTFLndsDxru226")
		tokenId++
	}
	wg.Wait()
	ipfs, err := addFolderToIpfs("/metadata")
	if err != nil {
		return ""
		// handle err
	}
	return ipfs
}
func imageHandler(w http.ResponseWriter, r *http.Request) {
	version := getEnv("version")
	p := Payload{
		Version: version,
		Input: Input{
			Prompt: "A tiger on a boat",
		},
	}
	resp, err := postData("https://api.replicate.com/v1/predictions", p)
	if err != nil {
		panic(err)
	}
	var d DiffusionResp
	checkValid := json.Valid(resp)

	if checkValid {
		fmt.Println("JSON was valid")
		json.Unmarshal(resp, &d)
	} else {
		fmt.Println("JSON is not valid!")
	}
	fmt.Printf("%#v\n", d.Urls.Get)
	imageUrl := retrieveImage(d.Urls.Get)
	pushUrlToIpfs(imageUrl, "test1")
	io.WriteString(w, "This is my website!")
}
func retrieveImage(url string) string {
	// Loop
	attempts := 0
	for attempts <= 1000 {
		time.Sleep(1 * time.Second)
		status, getUrl := getDiffusionState(url)
		fmt.Println(status)
		if status == "succeeded" {
			return getUrl
		}
	}
	return ""
}
func postData(url string, payload Payload) ([]byte, error) {
	jsonValue, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}
	token := getEnv("token")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}
func getDiffusionState(url string) (string, string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "error", ""
	}
	req.Header.Set("Authorization", "Token 55161acfcb4b8c4cb767169002c204c0b09f9dae")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "error", ""
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "error", ""
	}
	var d DiffusionGet
	checkValid := json.Valid(body)

	if checkValid {
		fmt.Println("JSON was valid")
		json.Unmarshal(body, &d)
		fmt.Printf("%#v\n", d)
		if len(d.Output) > 0 {
			return d.Status, d.Output[0]
		} else {
			return d.Status, ""
		}

	} else {
		fmt.Println("JSON is not valid!")
		return "error", ""
	}

}
func handleError(err error) {
	fmt.Println("Error Occured!")
	fmt.Println(err)
}
