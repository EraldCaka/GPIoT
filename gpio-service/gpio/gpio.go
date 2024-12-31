package gpio

type PinMode int

const (
	Input  PinMode = 0
	Output PinMode = 1
	Both   PinMode = 2
)

type PinType string

const (
	Digital PinType = "digital"
	Analog  PinType = "analog"
)

type Pin interface {
	SetMode(mode PinMode) error
	Read() (float64, error)
	Write(value float64) error
	Start() error
	Stop() error
	GetType() PinType
}
