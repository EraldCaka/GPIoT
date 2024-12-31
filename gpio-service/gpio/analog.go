package gpio

import (
	"errors"
	"math"
	"sync"
)

type AnalogPin struct {
	pinNumber int
	mode      PinMode
	value     float64 // voltage value any number
	mu        sync.Mutex
}

func NewAnalogPin(pinNumber int) *AnalogPin {
	return &AnalogPin{
		pinNumber: pinNumber,
		mode:      Input,
		value:     0.0,
	}
}

func (p *AnalogPin) SetMode(mode PinMode) error {
	if mode != Input && mode != Output {
		return errors.New("invalid mode for analog pin")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.mode = mode
	return nil
}

func (p *AnalogPin) Read() (float64, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.mode == Input {
		return 0, errors.New("pin is not in output mode")
	}
	return p.value, nil
}

func (p *AnalogPin) Write(value float64) error {
	// if value < 0 || value > 5 {
	// 	return errors.New("analog value out of range (0-5V)")
	// }
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.mode == Output {
		return errors.New("pin is not in input mode")
	}
	p.value = math.Max(0, math.Min(5, value))
	return nil
}

func (p *AnalogPin) Start() error {
	return nil
}

func (p *AnalogPin) Stop() error {
	return nil
}

func (p *AnalogPin) GetType() PinType {
	return Analog
}
