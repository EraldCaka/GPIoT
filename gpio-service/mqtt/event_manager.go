package events

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

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
		match, err := regexp.MatchString(`^[0-9]*\.?[0-9]+$`, payload)
		if !match {
			fmt.Printf("error: payload does not contain valid numeric data: %s\n", err)
			return
		}
		//	fmt.Printf("Received message from topic %s: %s\n", topic, payload)
		for _, digital := range c.config.GPIO.DigitalPins {
			if topic == digital.Topic && digital.Mode == gpio.Input {
				state := 0
				if payload == "1" {
					state = 1
				} else if payload == "0" {
					state = 0
				}

				err := handleDigitalPin(client, digital, state)
				if err != nil {
					fmt.Printf("Error processing pin %d: %v\n", digital.Pin, err)
				}
				return
			}
			/*
				else if digital.Mode != gpio.Input {
					message := config.MQTTMessage{
						Pin:   digital.Pin,
						State: digital.State,
						Error: fmt.Sprintf("error: couldnt update pin: %v because it is in output mode", digital.Pin),
					}
					client.Publish(digital.Path, 0, false, message.String())

				}
			*/
		}

		for _, analog := range c.config.GPIO.AnalogPins {
			if topic == analog.Topic && analog.Mode != gpio.Output {

				value, err := strconv.ParseFloat(payload, 64)
				if err != nil {
					fmt.Printf("error: couldn't parse string to float for payload: %s, error: %v\n", payload, err)
					return
				}
				err = handleAnalogPin(client, analog, value)
				if err != nil {
					fmt.Printf("Error processing pin %d: %v\n", analog.Pin, err)
				}
				return
			}
			/*
				else if analog.Mode == gpio.Output {
					message := config.MQTTMessage{
						Pin:   analog.Pin,
						State: analog.State,
						Error: fmt.Sprintf("error: couldnt update pin: %v because it is in output mode", analog.Pin),
					}
					client.Publish(analog.Path, 0, false, message.String())

				}
			*/
		}
	}
}

func handleDigitalPin(client mqtt.Client, pin config.DigitalPinConfig, state int) error {
	digitalPin := gpio.NewDigitalPin(pin.Pin, pin.Mode, pin.State, pin.Path)
	err := digitalPin.Write(state)
	if err != nil {
		return fmt.Errorf("failed to Write state %d to digital pin %d: %v", 1, pin.Pin, err)
	}

	message := config.MQTTMessage{
		Pin:   pin.Pin,
		State: state,
		Type:  digitalPin.GetType(),
	}
	token := client.Publish(pin.Path, 0, false, message.String())
	if token.Wait() && token.Error() != nil {
		log.Printf("Error publishing message: %v", token.Error())
	} else {
		fmt.Printf("Message sent to topic %s: \n", pin.Topic)
	}

	return nil
}

func handleAnalogPin(client mqtt.Client, pin config.AnalogPinConfig, state float64) error {
	analogPin := gpio.NewAnalogPin(pin.Pin, pin.Mode, pin.State, pin.Path)
	err := analogPin.Write(state)
	if err != nil {
		return fmt.Errorf("failed to Write state %d to digital pin %d: %v", 1, pin.Pin, err)
	}

	message := config.MQTTMessage{
		Pin:   pin.Pin,
		State: state,
		Type:  analogPin.GetType(),
	}
	if pin.Pin == 3 {
		log.Println("state for 3", state)
	}
	token := client.Publish(pin.Topic, 0, false, message.String())
	if token.Wait() && token.Error() != nil {
		log.Printf("Error publishing message: %v", token.Error())
	} else {
		fmt.Printf("Message sent to topic %s: \n", pin.Topic)
	}

	return nil
}
