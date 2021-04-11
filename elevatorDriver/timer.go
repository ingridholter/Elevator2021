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

/* Prøv å gjøre dette for timerDoor()
go func() {
	//send newOrder message for 2 seconds then stop.
	for timeout := time.After(1 * time.Second); ; {
		select {
		case <-timeout:
			return
		default:
		}
		NewOrderMsgTx <- l
		time.Sleep(100 * time.Millisecond)
	}
}()

var Timer = time.NewTimer(3 * time.Second)

func TimerDoor() { //funker som en sleep i 3 sekunder

	fmt.Println("TIMER")

	<-Timer.C

	fmt.Println("timer expired")
}
*/
