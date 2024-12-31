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
	gpioHandler := services.NewGPIOHandler(mqttConfig, client)
	gpioHandler.InitPins()

	err = client.Subscribe("gpio/control/#", messageHandler.HandleControlMessage())
	if err != nil {
		fmt.Printf("Error subscribing to control topics: %v\n", err)
		return
	}
	timer := mqttConfig.MonitorTime
	ticker := time.NewTicker(timer)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			gpioHandler.EventsTriggerer()
			// for i, pin := range gpio.GetAllDigitalPins() {
			// 	state, _ := pin.Read()
			// 	log.Println("pin nr", i, " state :", state)
			// }
			// for i, pin := range gpio.GetAllAnalogPins() {
			// 	state, _ := pin.Read()
			// 	log.Println("pin nr", i, " state :", state)
			// }
		}
	}()

	defer client.Disconnect()
	defer fmt.Println("Disconnected from MQTT broker")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
}
