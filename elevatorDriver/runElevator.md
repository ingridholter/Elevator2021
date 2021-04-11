package elevatorDriver

import (
	"fmt"
	. "main/config"
	. "main/elevio"
	"time"
)

func Run_Elevator(chanButtons chan ButtonEvent, chanFloors chan int, chanObstr chan bool, chanStop chan bool) {
	d := MD_Up

	chan_timer := make(chan bool, 1)
	go Timer(chan_timer)
	Timer := time.NewTimer(3 * time.Second)
	Timer.Stop()

	for {
		select {
		case b := <-chanButtons:
			msg := NewOrderDistributer(ElevStateArray, b.Button, b.Floor, id, Elevator) //ny mld med hvem som skal ha ordre!
			go func() {
				//send newOrder message for 2 seconds then stop.
				for timeout := time.After(1 * time.Second); ; {
					select {
					case <-timeout:
						return
					default:
					}
					NewOrderMsgTx <- msg
					time.Sleep(100 * time.Millisecond)
				}
			}()

		case f := <-chanFloors:

		case o := <-chanObstr:
			fmt.Printf("%+v\n", o)
			if o {
				SetMotorDirection(MD_Stop)
			} else {
				SetMotorDirection(d)
			}

		case s := <-chanStop:
			fmt.Printf("%+v\n", s)
			for f := 0; f < NumFloors; f++ {
				for b := ButtonType(0); b < 3; b++ {
					SetButtonLamp(b, f, false)
				}
			}

		case <-Timer.C:
			OnFloorTimeOut()

		}
	}
}
