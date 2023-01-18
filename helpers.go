package main

import (
	"io/ioutil"
	"log"
	"os"
)

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
