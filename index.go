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
	"time"

	shell "github.com/ipfs/go-ipfs-api"
)

func main() {
	http.HandleFunc("/retrieveImage", imageHandler)
	port := ":3000"
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Running API on port " + port)
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
	pushToIpfs(imageUrl, "test1")
	//// publishIpns("/ipfs/QmbkDDY6u15Tdc16KKrHXgNu89aAGxbxcG2oCchGVEtR7A")
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
func pushToIpfs(url string, key string) {
	sh := shell.NewShell("localhost:5001")
	resp, err := http.Get(url)
	if err != nil {
		handleError(err)
	}
	defer resp.Body.Close()

	imageBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		handleError(err)
	}
	r := bytes.NewReader(imageBytes)
	s, err := sh.Add(r)
	if err != nil {
		handleError(err)
	}
	err = sh.Pin(s)
	if err != nil {
		fmt.Println("Error with Pinning ipfs")
	}
	publishIpns(s, key)
}

// func generateKey() {
// 	ctx := context.Background()
// 	// Generate Key
// 	_, err := KeyGen(ctx, "test1")
// 	if err != nil {
// 		fmt.Println("Key Gen Failed!")
// 		log.Fatal(err)
// 	}
// }

func publishIpns(ipfsPath string, keyName string) {
	sh := shell.NewShell("localhost:5001")
	resp, err := sh.PublishWithDetails(ipfsPath, keyName, time.Second*100000, time.Second, false)
	if err != nil {
		fmt.Println("Error3")
		log.Fatal(err)
	}
	fmt.Println(resp)
}

// KeyGen Create a new keypair
func KeyGen(ctx context.Context, name string) (*Key, error) {
	sh := shell.NewShell("localhost:5001")
	rb := sh.Request("key/gen", name)

	var out Key
	if err := rb.Exec(ctx, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// func testCrypt() {
// 	privateKey, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
// 	if err != nil {
// 		panic(err)
// 	}
// 	keyBytes := privateKey.Raw
// 	data, err := keyBytes()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	writingBytesToFile("/test/data.txt", data)
// }
