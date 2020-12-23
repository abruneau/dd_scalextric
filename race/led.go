package race

import "time"

func ledCountdown() {
	sequence := [][]bool{
		{true, false, false, false, false, false, false, false},
		{true, true, false, false, false, false, false, false},
		{true, true, true, false, false, false, false, false},
		{false, false, false, true, false, false, false, false},
	}

	for _, v := range sequence {
		register.Write(v)
		time.Sleep(1 * time.Second)
	}

	register.Clear()
}

func ledFinish() {
	sequence := []bool{true, true, true, true, false, false, false, false}

	for i := 0; i < 3; i++ {
		register.Write(sequence)
		time.Sleep(200 * time.Millisecond)
		register.Clear()
		time.Sleep(200 * time.Millisecond)
	}
}
