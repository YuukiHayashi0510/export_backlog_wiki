package main

import (
	"encoding/json"
	"log"

	"github.com/kenzo0107/backlog"
)

const (
	spaceName = ""
	baseUrl   = "https://" + spaceName + ".backlog.com"
	projectID = ""
	apiKey    = ""
)

// wikiを取得する
func main() {
	c := backlog.New(apiKey, baseUrl)

	wiki, err := c.GetWiki(3080262)
	if err != nil {
		log.Fatal(err)
	}

	data, err := json.Marshal(wiki)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(data))
}
