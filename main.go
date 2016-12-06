package main

import (
	"encoding/json"
	"log"

	hc "cirello.io/HumorChecker"
	"github.com/kelseyhightower/envconfig"
	"github.com/timberslide/gotimberslide"
)

type specification struct {
	TsToken  string `envconfig:"ts_token"`
	TsHost   string `envconfig:"ts_host" default:"gw.timberslide.com:443"`
	InTopic  string `envconfig:"ts_in_topic"`
	OutTopic string `envconfig:"ts_out_topic"`
}

// Event contains information from an interesting tweet
type Event struct {
	Timestamp string   `json:"timestamp"` // Timestamp of tweet
	TweetID   int64    `json:"tweetID"`   // ID of tweet
	Raw       string   // Raw text of the tweet
	Keywords  []string // Any keywords that were found
}

func main() {
	var s specification
	err := envconfig.Process("APP", &s)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(s)

	client, err := ts.NewClient(s.TsHost, s.TsToken)
	err = client.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	ch, err := client.CreateChannel(s.OutTopic)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Starting")
	for event := range client.Iter(s.InTopic, ts.PositionOldest) {
		e := &Event{}
		err := json.Unmarshal([]byte(event.Message), e)
		if err != nil {
			log.Println(err)
			continue
		}
		msg, err := json.Marshal(hc.Analyze(event.Message))
		if err != nil {
			log.Println(err)
			continue
		}
		ch.Send(string(msg))
	}
}
