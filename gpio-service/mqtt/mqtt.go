package events

import (
	"fmt"
	"log"

	"github.com/EraldCaka/GPIoT/gpio-service/config"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClient struct {
	client mqtt.Client
	config *config.MQTTConfig
	done   chan bool
}

func NewMQTTClient(cfg *config.MQTTConfig) *MQTTClient {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(cfg.Broker).
		SetClientID(cfg.ClientID).
		SetDefaultPublishHandler(messageHandler).
		SetCleanSession(true)
	opts.SetClientID(cfg.ClientID)
	return &MQTTClient{
		config: cfg,
		done:   make(chan bool),
	}
}

func messageHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

func (c *MQTTClient) Connect() error {
	c.client = mqtt.NewClient(mqtt.NewClientOptions().AddBroker(c.config.Broker).SetClientID(c.config.ClientID))
	token := c.client.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (c *MQTTClient) Disconnect() {
	c.client.Disconnect(250)
	close(c.done)
	fmt.Println("MQTT client disconnected")
}

func (c *MQTTClient) Subscribe(topic string, callback mqtt.MessageHandler) error {
	token := c.client.Subscribe(topic, byte(c.config.QoS), callback)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	log.Printf("Successfully subscribed to topic: %s", topic)
	return nil
}

func (c *MQTTClient) Publish(topic string, message string) mqtt.Token {
	return c.client.Publish(topic, byte(c.config.QoS), c.config.Retain, message)
}

func (c *MQTTClient) Done() <-chan bool {
	return c.done
}
