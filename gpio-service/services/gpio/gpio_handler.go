package services

import (
	"fmt"
	"log"
	"time"

	"github.com/EraldCaka/GPIoT/gpio-service/config"
	"github.com/EraldCaka/GPIoT/gpio-service/gpio"
	events "github.com/EraldCaka/GPIoT/gpio-service/mqtt"
)

type GpioHandler struct {
	Config *config.MQTTConfig
	Client *events.MQTTClient
}

func NewGPIOHandler(config *config.MQTTConfig, client *events.MQTTClient) *GpioHandler {
	return &GpioHandler{
		Config: config,
		Client: client,
	}
}

func (g *GpioHandler) InitPins() {
	for _, digital := range g.Config.GPIO.DigitalPins {
		pin := gpio.NewDigitalPin(digital.Pin, digital.Mode, digital.State)
		gpio.RegisterDigitalPin(digital.Pin, pin)
		go g.EventHandler(digital.Pin, digital.Topic)
	}
	for i, digital := range gpio.GetAllDigitalPins() {

		log.Printf("digital pin nr: %v has state", i)
		log.Println("body", digital)

	}
}

func (g *GpioHandler) EventHandler(pinNumber int, topic string) {
	for {
		pin, err := gpio.GetDigitalPin(pinNumber)
		if err != nil {
			log.Printf("Error getting digital pin %d: %v", pinNumber, err)
			time.Sleep(5 * time.Second)
			continue
		}

		state, err := pin.Read()
		if err != nil {
			log.Printf("Error reading pin %d: %v", pinNumber, err)
			time.Sleep(5 * time.Second)
			continue
		}

		message := fmt.Sprintf("Pin %v has state: %v", pinNumber, state)
		g.Client.Publish(topic, message)
		time.Sleep(g.Config.MonitorTime * time.Second)

		select {
		case <-g.Client.Done():
			return
		}
	}
}

func (g *GpioHandler) HealthCheck() {
	err := g.Client.Publish("gpio/control/#", "ok")
	if err != nil {
		fmt.Printf("Error subscribing to control topics: %v\n", err)
		return
	}
}
