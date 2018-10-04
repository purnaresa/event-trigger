package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/tokopedia/feeds/src/utils"
	"gopkg.in/robfig/cron.v2"
)

type (
	Event struct {
		Name     string   `json:"name"`
		URL      string   `json:"url"`
		Headers  []string `json:"headers"`
		Schedule string   `json:"schedule"`
	}
)

func main() {
	file := flag.String("file", "event.json", "Events list")
	flag.Parse()

	log.Printf("initiating app using %s", *file)
	c := cron.New()
	events := getEvent(*file)

	for _, v := range events {
		c.AddFunc(v.Schedule, func() { trigger(v) })
		log.Printf("run : *%s* at %s\n", v.Name, v.Schedule)
	}

	c.Start()

	select {}
}

func getEvent(file string) (events []Event) {
	featureFile, err := ioutil.ReadFile(file)
	if err != nil {
		utils.LogError.Println(err)
		return
	}
	err = json.Unmarshal(featureFile, &events)
	if err != nil {
		utils.LogError.Println(err)
		return
	}

	return
}

func trigger(event Event) {

	req, _ := http.NewRequest("GET", event.URL, nil)
	for _, v := range event.Headers {
		texts := strings.Split(v, "|")
		req.Header.Add(texts[0], texts[1])
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	log.Printf(string(body))
}
