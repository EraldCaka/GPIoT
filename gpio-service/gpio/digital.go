package gpio

import (
	"errors"
	"sync"
)

var (
	pins   = make(map[int]*DigitalPin)
	pinsMu sync.Mutex
)

type DigitalPin struct {
	pinNumber int
	mode      PinMode
	state     bool
	mu        sync.Mutex
}

func NewDigitalPin(pinNumber int, mode PinMode, state bool) *DigitalPin {
	return &DigitalPin{
		pinNumber: pinNumber,
		mode:      mode,
		state:     state,
	}
}

func (p *DigitalPin) SetMode(mode PinMode) error {
	if mode != Input && mode != Output {
		return errors.New("invalid mode for digital pin")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.mode = mode
	return nil
}

func (p *DigitalPin) Read() (bool, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.mode != Input {
		return false, errors.New("pin is not in out mode")
	}
	return p.state, nil
}

func (p *DigitalPin) Write(value bool) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.mode != Output {
		return errors.New("pin is not in input mode")
	}
	p.state = value
	return nil
}

func (p *DigitalPin) Start() error {
	return nil
}

func (p *DigitalPin) Stop() error {
	return nil
}

func (p *DigitalPin) GetType() PinType {
	return Digital
}

func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

func GetDigitalPin(pinNumber int) (*DigitalPin, error) {
	pinsMu.Lock()
	defer pinsMu.Unlock()
	if pin, exists := pins[pinNumber]; exists {
		return pin, nil
	}
	return nil, errors.New("pin not found")
}

func RegisterDigitalPin(pinNumber int, pin *DigitalPin) {
	pinsMu.Lock()
	defer pinsMu.Unlock()
	pins[pinNumber] = pin
}

func GetAllDigitalPins() map[int]*DigitalPin {
	pinsMu.Lock()
	defer pinsMu.Unlock()
	return pins
}
