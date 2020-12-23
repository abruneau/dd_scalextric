package race

import (
	"time"

	"gobot.io/x/gobot/drivers/gpio"
)

func buzzerCountdown() {
	for i := 0; i < 3; i++ {
		buzzer.Tone(gpio.A4, 0.1)
		time.Sleep(900 * time.Millisecond)
	}
	buzzer.Tone(gpio.A5, 1)
}

func buzzerFinish() {
	buzzer.Tone(gpio.A5, 2)
}
