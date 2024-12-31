package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/EraldCaka/GPIoT/gpio-service/gpio"

	"gopkg.in/yaml.v3"
)

type DigitalEventHandler func(pin int, state bool, topic string, monitorTime time.Duration)

type MQTTConfig struct {
	Broker      string        `yaml:"broker"`
	ClientID    string        `yaml:"client_id"`
	QoS         int           `yaml:"qos"`
	Retain      bool          `yaml:"retain"`
	GPIO        GPIOConfig    `yaml:"gpio"`
	MonitorTime time.Duration `yaml:"monitor_time"`
}

type GPIOConfig struct {
	DigitalPins []DigitalPinConfig `yaml:"digital-pins"`
	AnalogPins  []AnalogPinConfig  `yaml:"analog-pins"`
}

type DigitalPinConfig struct {
	Name  string       `yaml:"gpio-name"`
	Pin   int          `yaml:"pin"`
	Mode  gpio.PinMode `yaml:"mode"`
	State int          `yaml:"state"`
	Topic string       `yaml:"topic"`
	Path  string       `yaml:"path"`
}

type AnalogPinConfig struct {
	Name  string       `yaml:"gpio-name"`
	Pin   int          `yaml:"pin"`
	Mode  gpio.PinMode `yaml:"mode"`
	State float64      `yaml:"state"`
	Topic string       `yaml:"topic"`
	Path  string       `yaml:"path"`
}

func NewMQTTConfig() *MQTTConfig {
	return &MQTTConfig{
		Broker:   "tcp://localhost:1883",
		ClientID: "gpio-monitor",
		QoS:      1,
		Retain:   false,
		GPIO: GPIOConfig{
			DigitalPins: []DigitalPinConfig{
				{Name: "led", Pin: 1, Mode: gpio.Output, State: 0},
			},
			AnalogPins: []AnalogPinConfig{
				{Name: "temperature-sensor", Pin: 0, Mode: gpio.Input, State: 0.0},
				{Name: "light-sensor", Pin: 2, Mode: gpio.Both, State: 0.5},
			},
		},
	}
}

type State interface {
}

type MQTTMessage struct {
	Pin   int
	State State
	Type  gpio.PinType
	Error string
}

func (m MQTTMessage) String() string {
	payloadBytes, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		log.Printf("Error marshalling MQTT message: %v", err)
		return ""
	}
	return string(payloadBytes)
}

func LoadMQTTConfig(filePath string) (*MQTTConfig, error) {
	config := NewMQTTConfig()
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open MQTT config file: %v", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return nil, fmt.Errorf("failed to decode MQTT config: %v", err)
	}

	return config, nil
}
