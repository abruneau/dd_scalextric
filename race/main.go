package race

import (
	"fmt"
	"sync"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/abruneau/dd_scalextric/utils"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

// InitRace Initiate race
func InitRace(conf *utils.Configuration, client *statsd.Client) {
	r := raspi.NewAdaptor()

	startButton = gpio.NewButtonDriver(r, fmt.Sprint(conf.Button.Start))
	stopButton = gpio.NewButtonDriver(r, fmt.Sprint(conf.Button.Stop))
	sensor1 = NewPPRMotionDriver(r, fmt.Sprint(conf.Sensor.One))
	sensor1.SetName("Car 1")
	sensor2 = NewPPRMotionDriver(r, fmt.Sprint(conf.Sensor.Two))
	sensor1.SetName("Car 2")
	buzzer = gpio.NewBuzzerDriver(r, fmt.Sprint(conf.Countdown.Beep))
	register = NewShiftRegister(r, 8, fmt.Sprint(conf.Countdown.Data), fmt.Sprint(conf.Countdown.Clk), fmt.Sprint(conf.Countdown.Latch), "")

	work := func() {

		race := NewRace(conf.Laps)
		var raceStartTime time.Time
		var car1, car2 Lap

		startButton.On(gpio.ButtonPush, func(data interface{}) {
			if !race.Running {
				countdown()
				raceStartTime = time.Now()
				race.Start()
				finish()
			}

		})

		stopButton.On(gpio.ButtonPush, func(data interface{}) {
			fmt.Println("Stop Button pushed")
			if race.Running {
				race.Halt()
			}

		})

		race.On(RaceStart, func(data interface{}) {
			fmt.Println("Race Started")
			e := statsd.NewEvent("Race Start", "A new race started")
			client.Event(e)
			client.Incr("race.count", nil, 1)
		})
		race.On(RaceStop, func(data interface{}) {
			fmt.Println("Race Stoped")
			e := statsd.NewEvent("Race Stop", "The race was stoped")
			e.AlertType = statsd.Error
			client.Event(e)
		})
		race.On(RaceFinished, func(data interface{}) {
			fmt.Println("Race Finished")
			racetime := time.Now().Sub(raceStartTime).Milliseconds()
			e := statsd.NewEvent("Race Finish", "The race is finished")
			e.AlertType = statsd.Success
			client.Event(e)
			client.Gauge("race.duration", float64(racetime), nil, 1)
		})

		race.On(LapDone, func(lap interface{}) {
			l := lap.(Lap)
			var previousLap *Lap
			if l.Number > 0 {
				if l.Car == 1 {
					previousLap = &car1
				}
				if l.Car == 2 {
					previousLap = &car2
				}
				lapTime(&l, previousLap, client)
			}
			if l.Car == 1 {
				car1 = l
			}
			if l.Car == 2 {
				car2 = l
			}
		})
	}

	robot := gobot.NewRobot("buttonBot",
		[]gobot.Connection{r},
		[]gobot.Device{startButton, stopButton, sensor1, sensor2, buzzer, register},
		work,
	)

	robot.Start()
}

func countdown() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		buzzerCountdown()
		wg.Done()
	}()
	go func() {
		ledCountdown()
		wg.Done()
	}()
	wg.Wait()
}

func finish() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		buzzerFinish()
		wg.Done()
	}()
	go func() {
		ledFinish()
		wg.Done()
	}()
	wg.Wait()
}

func lapTime(current, previus *Lap, client *statsd.Client) {
	lapTime := current.Time.Sub(previus.Time).Milliseconds()
	client.Gauge("race.lap", float64(lapTime), []string{fmt.Sprintf("car:%v", current.Car), fmt.Sprintf("lap:%v", current.Number)}, 1)
	fmt.Printf("Car %v finished lap %v in %v milliseconds at %v\n", current.Car, current.Number, lapTime, current.Time)
}
