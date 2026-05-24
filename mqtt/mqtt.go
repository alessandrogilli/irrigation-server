package mqtt

import (
	"encoding/json"
	"errors"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var Client mqtt.Client

func Init(broker string) error {
	if broker == "" {
		fmt.Println("No MQTT broker configured, skipping MQTT setup")
		return nil
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID("go-scheduler")

	Client = mqtt.NewClient(opts)

	token := Client.Connect()
	token.Wait()

	if token.Error() != nil {
		return errors.New(token.Error().Error())
	}

	fmt.Println("Connected to MQTT broker")
	return nil
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
