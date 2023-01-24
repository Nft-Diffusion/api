package main

type Key struct {
	Id   string
	Name string
}

type Urls struct {
	Get    string
	Cancel string
}

type CreateImageReq struct {
	TokenId string `json:"tokenId"`
	Param   string `json:"param"`
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
type Metadata struct {
	Description  string `json:"description"`
	External_url string `json:"external_url"`
	Image        string `json:"image"`
	Name         string `json:"name"`
}

// {
//   "description": "Friendly OpenSea Creature that enjoys long swims in the ocean.",
//   "external_url": "https://openseacreatures.io/3",
//   "image": "https://storage.googleapis.com/opensea-prod.appspot.com/puffs/3.png",
//   "name": "Dave Starbelly",
//   "attributes": [ ... ]
// }
