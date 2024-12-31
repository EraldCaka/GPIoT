package gpio

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
)

var (
	pins   = make(map[int]*DigitalPin)
	pinsMu sync.Mutex
)

type DigitalPin struct {
	pinNumber int
	mode      PinMode
	state     int
	path      string
	mu        sync.Mutex
}

func NewDigitalPin(pinNumber int, mode PinMode, state int, path string) *DigitalPin {
	return &DigitalPin{
		pinNumber: pinNumber,
		mode:      mode,
		state:     state,
		path:      path,
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

func (p *DigitalPin) Read() (int, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.mode != Output {
		return 0, errors.New("pin is not in output mode")
	}

	data, err := os.ReadFile(p.path)
	if err != nil {
		return 0, fmt.Errorf("error reading file: %w", err)
	}

	var state int
	_, err = fmt.Sscanf(string(data), "%d", &state)
	if err != nil {
		return 0, fmt.Errorf("error parsing file content: %w", err)
	}
	if state < 0 || state > 1 {
		return 0, fmt.Errorf("error state is not low/high (0 or 1) value!")
	}
	return state, nil
}

func (p *DigitalPin) Write(value int) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.mode != Input {
		return errors.New("pin is not in input mode")
	}

	data := []byte(fmt.Sprintf("%d", value))
	err := os.WriteFile(p.path, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	p.state = value
	pins[p.pinNumber] = p
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

	data := []byte(fmt.Sprintf("%d", pin.state))
	err := os.WriteFile(pin.path, data, 0644)
	if err != nil {
		log.Println("error: couldnt write pin ", pinNumber, " inside the corresponding pin .txt file.", err)
		return
	}

	pins[pinNumber] = pin

}

func GetAllDigitalPins() map[int]*DigitalPin {
	pinsMu.Lock()
	defer pinsMu.Unlock()
	return pins
}
