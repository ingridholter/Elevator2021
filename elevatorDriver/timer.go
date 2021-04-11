package elevatorDriver

import (
	"time"
)

var TimerFlag int
var T time.Time

func Timer(chan_timer chan bool) {
	for {
		if TimerFlag == 1 {
			T = time.Now()
			TimerFlag = 0
		}
		if int(time.Now().Second())-int(T.Second()) > 3 {
			chan_timer <- true
		}
	}
}

func Timer_start() {
	TimerFlag = 1
}

