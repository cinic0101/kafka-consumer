package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"./config"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	log "github.com/sirupsen/logrus"
)

type Request struct {
	Index   string 			`json:"index"`
	Document interface{}	`json:"doc"`
}

const configPath = "consumer.yml"
const logPath = "log/consumer.log"

func main() {
	conf := config.Unmarshal(configPath)

	_ = os.Mkdir("log", 0700)
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
	 	log.SetOutput(io.MultiWriter(os.Stdout, file))
	} else {
	 	log.Info("Failed to log to file, using default stderr")
	}

	c, err := kafka.NewConsumer(&kafka.ConfigMap {
		"bootstrap.servers": conf.Kafka.BootstrapServers,
		"group.id":          conf.Kafka.GroupID,
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	log.Info("kafka-consumer was started")
	c.SubscribeTopics(strings.Split(conf.Kafka.Topics, ","), nil)
	log.Infof("Subscribed Topics: %s", conf.Kafka.Topics)

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			req := Request{}
			if e := json.Unmarshal(msg.Value, &req); e != nil {
				log.WithFields(log.Fields{ "kafka-msg": string(msg.Value) }).Warning("Received wrong format message")
				continue
			}

			doc, err := json.Marshal(req.Document)
			if err != nil {
				log.WithFields(log.Fields{ "kafka-msg": string(msg.Value) }).Warning("Received wrong format message")
				continue
			}

			resp, err := http.Post(
				fmt.Sprintf("http://%s/%s/%s", conf.Destination.Params["server"], req.Index, conf.Destination.Params["default-type"]),
				"application/json;charset=utf-8", bytes.NewReader(doc))

			if err != nil {
				log.WithFields(log.Fields{ "kafka-msg": string(msg.Value) }).Errorf("Something wrong with elasticsearch server: %s", err)
				continue
			}

			if resp.StatusCode > http.StatusIMUsed /* 226 // RFC 3229, 10.4.1 */ {
				b, _ := ioutil.ReadAll(resp.Body)
				log.WithFields(log.Fields{ "kafka-msg": string(msg.Value) }).Errorf("Received not ok response:[%v] %s", resp.StatusCode, string(b))
				continue
			}

			resp.Body.Close()

		} else {
			log.WithFields(log.Fields{ "kafka-msg": string(msg.Value) }).Errorf("Read kafka message fail: %v", err)
			break
		}
	}

	c.Close()
}
