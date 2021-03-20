package main

import (
	"fmt"
	"time"
)

func timerStart(sec int) {

	timer := timer.NewTimer(sec * time.Second)
	<-timer1.C
	fmt.Println("Timer fired")

}

func timerStop() {
	stopTimer := timer.Stop()
}
