package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/joho/godotenv"
)

// KeyGen Create a new keypair
func KeyGen(ctx context.Context, name string) (*Key, error) {
	defer wg.Done()
	sh := shell.NewShell("localhost:5001")
	rb := sh.Request("key/gen", name)

	var out Key
	if err := rb.Exec(ctx, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
func addFolderToIpfs(folderPath string) (string, error) {
	sh := shell.NewShell("localhost:5001")
	added, err := sh.AddDir(folderPath)
	if err != nil {
		fmt.Println(err)
		return "", errors.New("adding dir to ipfs")
	}
	return added, nil
}
func pushUrlToIpfs(url string) string {
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
	return s
}
func publishIpnsAndWrite(tokenId string, ipfs string) {
	defer wg.Done()
	ipns, err := publishIpns(ipfs, "token"+tokenId)
	if err != nil {
		fmt.Println("Error Occured in PublishIpnsAndWrite")
	}
	data := Metadata{
		Description:  "A description of NFT",
		Name:         "Nft Diffusion",
		External_url: "http://google.com",
		Image:        "ipns://" + ipns,
	}
	file, _ := json.MarshalIndent(data, "", "")
	_ = ioutil.WriteFile("metadata/"+tokenId+".json", file, 0644)
	fmt.Println("Finished!")

}
func publishIpns(ipfsPath string, keyName string) (string, error) {
	sh := shell.NewShell("localhost:5001")
	resp, err := sh.PublishWithDetails(ipfsPath, keyName, time.Second*100000, time.Second, false)
	if err != nil {
		fmt.Println("An Error Occured.")
		return "", err
	}
	fmt.Println(resp)
	return resp.Name, nil
}
func getEnv(env string) string {
	ve := os.Getenv(env)
	_, ok := os.LookupEnv(ve)
	if !ok {
		err := godotenv.Load()
		if err != nil {
			fmt.Println("Error loading .env file")
		}
	}
	version := os.Getenv(env)
	return version
}
