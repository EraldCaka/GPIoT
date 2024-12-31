package main

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	broker := "tcp://localhost:1883"
	clientID := "gpio-monitor"
	topic := "gpio/control/digital/1"
	message := "1"

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker).SetClientID(clientID)

	client := mqtt.NewClient(opts)

	// Connect to MQTT broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error connecting to MQTT broker: %v", token.Error())
	}

	token := client.Publish(topic, 1, false, message)
	if token.Wait() && token.Error() != nil {
		log.Printf("Error publishing message: %v", token.Error())
	} else {
		fmt.Printf("Message sent to topic %s: %s\n", topic, message)
	}

	client.Disconnect(250)
}
