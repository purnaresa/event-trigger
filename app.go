package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	cron "gopkg.in/robfig/cron.v2"
)

type (
	// Event data structure for event
	Event struct {
		Name     string   `json:"name"`
		Method   string   `json:"method"`
		URL      string   `json:"url"`
		Headers  []string `json:"headers"`
		Schedule string   `json:"schedule"`
	}
)

var (
	logger *log.Logger
)

func main() {
	eventFile := flag.String("event", "event.json", "Events list")
	logFile := flag.String("log", "event.log", "log file")
	flag.Parse()
	f, err := os.OpenFile(*logFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	logger = log.New(f, "", log.LstdFlags)

	log.Printf("initiating app using %s", *eventFile)
	events := getEvent(*eventFile)
	c := cron.New()

	for _, event := range events {
		log.Printf("run : *%s* at %s\n", event.Name, event.Schedule)
		go func(v Event) {
			c.AddFunc(v.Schedule, func() { trigger(v) })
		}(event)
	}

	c.Start()
	select {}
}

func getEvent(file string) (events []Event) {
	featureFile, err := ioutil.ReadFile(file)
	if err != nil {
		logger.Println(err)
		return
	}
	err = json.Unmarshal(featureFile, &events)
	if err != nil {
		logger.Println(err)
		return
	}

	return
}

func trigger(event Event) {

	req, _ := http.NewRequest(event.Method, event.URL, nil)
	for _, v := range event.Headers {
		texts := strings.Split(v, "|")
		req.Header.Add(texts[0], texts[1])
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Println(err)
		return
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	logger.Printf("event : *%s* respond : %s", event.Name, string(body))
}
