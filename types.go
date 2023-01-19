package main

type Key struct {
	Id   string
	Name string
}
type Urls struct {
	Get    string
	Cancel string
}
type DiffusionResp struct {
	Completed_at      string
	Created_at        string
	Error             string
	Id                string
	Input             Input
	Logs              []string
	Metrics           string
	Output            []string
	Started_at        string
	Status            string
	Urls              Urls
	Version           string
	Webhook_completed string
}
type DiffusionGet struct {
	Completed_at      string
	Created_at        string
	Error             string
	Id                string
	Input             Input
	Logs              string
	Metrics           string
	Output            []string
	Started_at        string
	Status            string
	Urls              Urls
	Version           string
	Webhook_completed string
}
type Input struct {
	Prompt string `json:"prompt"`
}
type Payload struct {
	Version string `json:"version"`
	Input   Input  `json:"input"`
}
