package gpio

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
)

var (
	analogPins = make(map[int]*AnalogPin)
)

type AnalogPin struct {
	pinNumber int
	mode      PinMode
	state     float64
	path      string
	mu        sync.Mutex
}

func NewAnalogPin(pinNumber int, mode PinMode, state float64, path string) *AnalogPin {
	return &AnalogPin{
		pinNumber: pinNumber,
		mode:      mode,
		state:     state,
		path:      path,
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

	if p.mode == Output {
		return 0, errors.New("pin is not in input mode")
	}

	fileContent, err := os.ReadFile(p.path)
	if err != nil {
		return 0, fmt.Errorf("error reading pin state from file %s: %w", p.path, err)
	}

	strState := string(fileContent)
	state, err := strconv.ParseFloat(strState, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing file content: %w", err)
	}

	return state, nil
}

func (p *AnalogPin) Write(value float64) error {
	if value < 0 || value > 5 {
		return errors.New("analog value out of range (0-5V)")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.mode == Input {
		return errors.New("pin is not in Output mode")
	}

	data := fmt.Sprintf("%f", value)
	err := os.WriteFile(p.path, []byte(data), 0644)
	if err != nil {
		return fmt.Errorf("failed to write analog pin value to file: %w", err)
	}

	p.state = value
	return nil
}

func GetAnalogPin(pinNumber int) (*AnalogPin, error) {
	pinsMu.Lock()
	defer pinsMu.Unlock()
	if pin, exists := analogPins[pinNumber]; exists {
		return pin, nil
	}
	return nil, errors.New("pin not found")
}
func stringifyFloat64(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}
func RegisterAnalogPin(pinNumber int, pin *AnalogPin) {
	pinsMu.Lock()
	defer pinsMu.Unlock()

	data := []byte(stringifyFloat64(pin.state))
	err := os.WriteFile(pin.path, data, 0644)
	if err != nil {
		log.Println("error: couldnt write pin ", pinNumber, " inside the corresponding pin .txt file.", err)
		return
	}

	analogPins[pinNumber] = pin

}

func GetAllAnalogPins() map[int]*AnalogPin {
	pinsMu.Lock()
	defer pinsMu.Unlock()
	return analogPins
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
