package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/joho/godotenv"
)

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
func writingBytesToFile(dest string, b []byte) {
	file, err := os.Create(dest)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.Write(b)
}

func readingBytesFromFile(dest string) []byte {
	body, err := ioutil.ReadFile(dest)
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}
	return body
}
