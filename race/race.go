package race

import (
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
)

var (
	startButton, stopButton *gpio.ButtonDriver
	sensor1, sensor2        *PPRMotionDriver
	buzzer                  *gpio.BuzzerDriver
	register                *ShiftRegisterDriver
)

const (
	// RaceStart event
	RaceStart = "race-started"
	// RaceStop event
	RaceStop = "race-stoped"
	// RaceFinished event
	RaceFinished = "race-finished"
	// LapDone event
	LapDone = "lap-done"
)

// Race discribes the race
type Race struct {
	totalLaps int
	halt      chan bool
	Running   bool
	gobot.Eventer
}

// Lap describes a Lap
type Lap struct {
	Car    int
	Number int
	Time   time.Time
}

// NewRace returns a new race with a number off total laps
func NewRace(totalLaps int) *Race {
	return &Race{
		totalLaps: totalLaps,
		Eventer:   gobot.NewEventer(),
		halt:      make(chan bool),
		Running:   false,
	}
}

// Start starts the race
func (r *Race) Start() {
	r.Publish(RaceStart, nil)
	r.Running = true

	counter1 := 0
	counter2 := 0

	sensor1.On(gpio.MotionDetected, func(data interface{}) {
		if counter1 < r.totalLaps {
			l := Lap{
				Car:    1,
				Number: counter1,
				Time:   time.Now(),
			}
			r.Publish(LapDone, l)
			counter1++
		}
	})
	sensor2.On(gpio.MotionDetected, func(data interface{}) {
		if counter2 < r.totalLaps {
			l := Lap{
				Car:    2,
				Number: counter2,
				Time:   time.Now(),
			}
			r.Publish(LapDone, l)
			counter2++
		}
	})

	for {

		if counter1 >= r.totalLaps && counter2 >= r.totalLaps {
			r.Running = false
			r.Publish(RaceFinished, nil)
			return
		}

		select {
		case <-time.After(10 * time.Millisecond):
		case <-r.halt:
			r.Running = false
			r.Publish(RaceStop, nil)
			return
		}
	}

}

// Halt stops the race
func (r *Race) Halt() (err error) {
	r.halt <- true
	return
}
