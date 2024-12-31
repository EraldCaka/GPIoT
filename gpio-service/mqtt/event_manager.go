package events

import (
	"fmt"
	"log"

	"github.com/EraldCaka/GPIoT/gpio-service/config"
	"github.com/EraldCaka/GPIoT/gpio-service/gpio"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type ControlMessageHandler struct {
	config *config.MQTTConfig
}

func NewControlMessageHandler(config *config.MQTTConfig) *ControlMessageHandler {
	return &ControlMessageHandler{
		config: config,
	}
}

func (c *ControlMessageHandler) HandleControlMessage() func(mqtt.Client, mqtt.Message) {
	return func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		payload := string(msg.Payload())
		fmt.Printf("Received message from topic %s: %s\n", topic, payload)
		for _, digital := range c.config.GPIO.DigitalPins {
			if topic == digital.Topic {
				state := false
				if payload == "1" {
					state = true
				} else if payload == "0" {
					state = false
				}
				log.Println("here<<>>><<>>", digital)
				err := handleDigitalPin(client, digital.Pin, digital.Mode, state)
				if err != nil {
					fmt.Printf("Error processing pin %d: %v\n", digital.Pin, err)
				}
				return
			}
		}

		fmt.Printf("Unhandled control message: %s - %s\n", topic, payload)
	}
}

func handleDigitalPin(client mqtt.Client, pinNumber int, mode gpio.PinMode, state bool) error {

	digitalPin := gpio.NewDigitalPin(pinNumber, mode, state)
	err := digitalPin.Write(state)
	log.Println(">>>>>>>>>>--------------done", digitalPin)
	if err != nil {
		return fmt.Errorf("failed to write state %d to digital pin %d: %v", 1, pinNumber, err)
	}

	currentState, _ := digitalPin.Read()
	log.Println("current state updated to ", currentState)
	client.Publish(fmt.Sprintf("gpio/control/response/digital/%d", pinNumber), 0, false, fmt.Sprintf("Digital Pin %d turned %s", pinNumber, map[int]string{1: "ON", 0: "OFF"}[1]))
	log.Println(">>>>>>> published:", digitalPin)
	return nil
}
