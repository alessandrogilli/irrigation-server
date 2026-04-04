package mqtt

import (
	"encoding/json"
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var Client mqtt.Client

func Init(broker string) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID("go-scheduler")

	Client = mqtt.NewClient(opts)

	token := Client.Connect()
	token.Wait()

	if token.Error() != nil {
		log.Fatal(token.Error())
	}

	fmt.Println("Connected to MQTT broker")
}

type Message struct {
	Line int    `json:"line"`
	Cmd  string `json:"cmd"`
}

func Publish(lineID int, cmd string) {
	msg := Message{
		Line: lineID,
		Cmd:  cmd,
	}

	payload, _ := json.Marshal(msg)

	topic := "irrigation-server/lines"

	token := Client.Publish(topic, 0, false, payload)
	token.Wait()

	fmt.Printf("MQTT sent: %s\n", payload)
}
