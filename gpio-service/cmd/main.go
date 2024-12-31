package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/EraldCaka/GPIoT/gpio-service/config"
	events "github.com/EraldCaka/GPIoT/gpio-service/mqtt"
	services "github.com/EraldCaka/GPIoT/gpio-service/services/gpio"
)

func main() {
	mqttConfig, err := config.LoadMQTTConfig("configs.yaml")
	if err != nil {
		fmt.Printf("Error setting up configurations: %v\n", err)
		return
	}

	messageHandler := events.NewControlMessageHandler(mqttConfig)
	client := events.NewMQTTClient(mqttConfig)
	err = client.Connect()
	if err != nil {
		fmt.Printf("Error connecting to MQTT broker: %v\n", err)
		return
	}
	fmt.Println("Connected to MQTT broker")

	err = client.Subscribe("gpio/control/#", messageHandler.HandleControlMessage())
	if err != nil {
		fmt.Printf("Error subscribing to control topics: %v\n", err)
		return
	}
	gpioHandler := services.NewGPIOHandler(mqttConfig, client)
	gpioHandler.InitPins()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			gpioHandler.InitPins()
		}
	}()

	defer client.Disconnect()
	defer fmt.Println("Disconnected from MQTT broker")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
}
