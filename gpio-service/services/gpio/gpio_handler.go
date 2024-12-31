package services

import (
	"log"
	"sync"
	"time"

	"github.com/EraldCaka/GPIoT/gpio-service/config"
	"github.com/EraldCaka/GPIoT/gpio-service/gpio"
	events "github.com/EraldCaka/GPIoT/gpio-service/mqtt"
)

type GpioHandler struct {
	Config *config.MQTTConfig
	Client *events.MQTTClient
	mu     sync.Mutex
}

func NewGPIOHandler(config *config.MQTTConfig, client *events.MQTTClient) *GpioHandler {
	return &GpioHandler{
		Config: config,
		Client: client,
	}
}

func (g *GpioHandler) InitPins() {
	for _, digital := range g.Config.GPIO.DigitalPins {
		pin := gpio.NewDigitalPin(digital.Pin, digital.Mode, digital.State, digital.Path)
		gpio.RegisterDigitalPin(digital.Pin, pin)
		//	go g.EventHandler(digital.Pin, digital.Topic)
	}
	for _, analog := range g.Config.GPIO.AnalogPins {
		pin := gpio.NewAnalogPin(analog.Pin, analog.Mode, analog.State, analog.Path)
		gpio.RegisterAnalogPin(analog.Pin, pin)
		//	go g.EventHandler(digital.Pin, digital.Topic)
	}
}

func (g *GpioHandler) EventsTriggerer() {
	for _, digital := range g.Config.GPIO.DigitalPins {
		go g.EventHandlerDigital(digital.Pin, digital.Topic)
	}

	for _, analog := range g.Config.GPIO.AnalogPins {
		go g.EventHandlerAnalog(analog.Pin, analog.Topic)
	}
}
func (g *GpioHandler) EventHandlerDigital(pinNumber int, topic string) {
	for {
		g.mu.Lock()
		pin, err := gpio.GetDigitalPin(pinNumber)
		if err != nil {
			g.mu.Unlock()
			log.Printf("Error getting pin %d: %v", pinNumber, err)
			time.Sleep(5 * time.Second)
			return
		}
		state, err := pin.Read()
		g.mu.Unlock()

		if err != nil {
			log.Printf("Error reading pin %d: %v", pinNumber, err)
			time.Sleep(5 * time.Second)
			return
		}

		message := config.MQTTMessage{
			Pin:   pinNumber,
			State: state,
			Type:  pin.GetType(),
		}
		token := g.Client.Publish(topic, message.String())
		if token.Wait() && token.Error() != nil {
			log.Printf("Error publishing message: %v", token.Error())
		} else {
			//	fmt.Printf("Message sent to topic %s\n", topic)
		}
		time.Sleep(g.Config.MonitorTime * time.Second)

		select {
		case <-g.Client.Done():
			return
		default:
		}
	}
}

func (g *GpioHandler) EventHandlerAnalog(pinNumber int, topic string) {
	for {
		g.mu.Lock()
		pin, err := gpio.GetAnalogPin(pinNumber)
		if err != nil {
			g.mu.Unlock()
			log.Printf("Error getting pin %d: %v", pinNumber, err)
			time.Sleep(5 * time.Second)
			return
		}
		state, err := pin.Read()
		g.mu.Unlock()

		if err != nil {
			log.Printf("Error reading pin %d: %v", pinNumber, err)
			time.Sleep(5 * time.Second)
			return
		}

		message := config.MQTTMessage{
			Pin:   pinNumber,
			State: state,
			Type:  pin.GetType(),
		}
		token := g.Client.Publish(topic, message.String())
		if token.Wait() && token.Error() != nil {
			log.Printf("Error publishing message: %v", token.Error())
		} else {
			//		fmt.Printf("Message sent to topic %s\n", topic)
		}
		time.Sleep(g.Config.MonitorTime * time.Second)

		select {
		case <-g.Client.Done():
			return
		default:
		}
	}
}
