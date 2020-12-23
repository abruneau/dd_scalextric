package race

import (
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
)

type pin struct {
	number     string
	high       bool
	connection gpio.DigitalWriter
}

// On sets the led to a high state.
func (p *pin) On() (err error) {
	if err = p.connection.DigitalWrite(p.number, 1); err != nil {
		return
	}
	p.high = true
	return
}

// On sets the led to a high state.
func (p *pin) Off() (err error) {
	if err = p.connection.DigitalWrite(p.number, 0); err != nil {
		return
	}
	p.high = false
	return
}

// ShiftRegisterDriver represents a digital Shift register
type ShiftRegisterDriver struct {
	name       string
	connection gpio.DigitalWriter
	dataPin    pin // ser
	clockPin   pin // srclk
	latchPin   pin // rclk
	clearPin   pin // srclr
	lineSize   int
	gobot.Commander
}

// NewShiftRegister return a new ShiftRegisterDriver given a DigitalWriter, line size, and pins.
func NewShiftRegister(a gpio.DigitalWriter, lineSize int, ser, srclk, rclk, srclr string) *ShiftRegisterDriver {
	register := &ShiftRegisterDriver{
		name:       gobot.DefaultName("Shift Register"),
		connection: a,
		lineSize:   lineSize,
		dataPin:    pin{number: ser, connection: a},
		clockPin:   pin{number: srclk, connection: a},
		latchPin:   pin{number: rclk, connection: a},
		clearPin:   pin{number: srclr, connection: a},
		Commander:  gobot.NewCommander(),
	}

	register.dataPin.Off()
	register.clockPin.Off()
	register.latchPin.Off()
	register.Clear()

	return register
}

// Name returns the ShiftRegisterDriver name
func (s *ShiftRegisterDriver) Name() string { return s.name }

// SetName sets the ShiftRegisterDriver name
func (s *ShiftRegisterDriver) SetName(n string) { s.name = n }

// Start implements the Driver interface
func (s *ShiftRegisterDriver) Start() (err error) { return }

// Halt implements the Driver interface
func (s *ShiftRegisterDriver) Halt() (err error) { return }

// Connection returns the Connection associated with the Driver
func (s *ShiftRegisterDriver) Connection() gobot.Connection { return s.connection.(gobot.Connection) }

// Write writes to the register based on a boolean slice
func (s *ShiftRegisterDriver) Write(data []bool) {
	missingLen := s.lineSize - len(data)
	if missingLen < 0 {
		missingLen = 0
	}

	actualData := append(data, make([]bool, missingLen, missingLen)...)
	actualData = actualData[:s.lineSize]

	for index := range actualData {
		pin := actualData[len(actualData)-index-1]
		if pin {
			s.dataPin.On()
		}
		s.clockPin.On()
		s.clockPin.Off()
		s.dataPin.Off()
	}

	s.latchPin.On()
	s.latchPin.Off()
}

// Clear clears the register
func (s *ShiftRegisterDriver) Clear() {
	empty := make([]bool, 0, 0)
	s.Write(empty)
}
