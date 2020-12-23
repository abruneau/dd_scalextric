package race

import (
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
)

// PPRMotionDriver represents a digital Proximity Photo Resistor (PPR) motion detecter
type PPRMotionDriver struct {
	Active     bool
	pin        string
	name       string
	halt       chan bool
	interval   time.Duration
	connection gpio.DigitalReader
	gobot.Eventer
}

// NewPPRMotionDriver returns a new PPRMotionDriver with a polling interval of
// 10 Milliseconds given a DigitalReader and pin.
//
// Optionally accepts:
//  time.Duration: Interval at which the PPRMotionDriver is polled for new information
func NewPPRMotionDriver(a gpio.DigitalReader, pin string, v ...time.Duration) *PPRMotionDriver {
	b := &PPRMotionDriver{
		name:       gobot.DefaultName("PPRMotion"),
		connection: a,
		pin:        pin,
		Active:     false,
		Eventer:    gobot.NewEventer(),
		interval:   10 * time.Millisecond,
		halt:       make(chan bool),
	}

	if len(v) > 0 {
		b.interval = v[0]
	}

	b.AddEvent(gpio.MotionDetected)
	b.AddEvent(gpio.MotionStopped)
	b.AddEvent(gpio.Error)

	return b
}

// Start starts the PPRMotionDriver and polls the state of the sensor at the given interval.
//
// Emits the Events:
// 	MotionDetected - On motion detected
//	MotionStopped int - On motion stopped
//	Error error - On button error
//
// The PPRMotionDriver will send the MotionStopped event over and over,
// just as long as motion is still not being detected.
// It will only send the MotionDetected event once, however, until
// motion stop being detected again
func (p *PPRMotionDriver) Start() (err error) {
	go func() {
		for {
			newValue, err := p.connection.DigitalRead(p.Pin())
			if err != nil {
				p.Publish(gpio.Error, err)
			}
			switch newValue {
			case 1:
				if !p.Active {
					p.Active = true
					p.Publish(gpio.MotionStopped, newValue)
				}
			case 0:
				if p.Active {
					p.Active = false
					p.Publish(gpio.MotionDetected, newValue)
				}
			}

			select {
			case <-time.After(p.interval):
			case <-p.halt:
				return
			}
		}
	}()
	return
}

// Halt stops polling the button for new information
func (p *PPRMotionDriver) Halt() (err error) {
	p.halt <- true
	return
}

// Name returns the PPRMotionDriver name
func (p *PPRMotionDriver) Name() string { return p.name }

// SetName sets the PPRMotionDriver name
func (p *PPRMotionDriver) SetName(n string) { p.name = n }

// Pin returns the PPRMotionDriver pin
func (p *PPRMotionDriver) Pin() string { return p.pin }

// Connection returns the PPRMotionDriver Connection
func (p *PPRMotionDriver) Connection() gobot.Connection { return p.connection.(gobot.Connection) }
