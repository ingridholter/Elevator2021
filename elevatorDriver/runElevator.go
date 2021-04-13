package elevatorDriver

import (
	"fmt"
	. "main/config"
	. "main/elevio"
	"time"
)

func DrvElevator(id string, chanNewOrder <-chan ButtonEvent, chanFloors <-chan int, chanObstr <-chan bool, chanStop <-chan bool, chanElevator chan ElevState, ElevStateMsgTx chan<- ElevStateMsg) {

	elevator := ElevState{
		Floor:     GetFloor(),
		Dir:       MD_Stop,
		Behaviour: EBidle,
		Requests:  [NumFloors][NumButtons]bool{},
	}
	if Between {
		fmt.Println("on init between floors")
		OnInitBetweenFloors(elevator, chanElevator)
	} else {
		chanElevator <- elevator
	}
	Timer := time.NewTimer(3 * time.Second)
	Timer.Stop()
	SendStateTicker := time.NewTicker(500 * time.Millisecond)

	for {

		select {
		case <-SendStateTicker.C:

			elevator = <-chanElevator

			ElevStateMsgTx <- ElevStateMsg{
				SenderId: id,
				Message:  "State Update",
				Elevator: elevator,
			}

			chanElevator <- elevator

		case b := <-chanNewOrder: //SyncAllLights(ElevStateArray, id)
			fmt.Println("lagt til i min ordre kanal")
			elevator = <-chanElevator

			OnRequestButtonPress(elevator, b.Floor, b.Button, Timer, chanElevator)

		case f := <-chanFloors: //SyncAllLights(ElevStateArray, id)
			fmt.Println("In case Floor: ", f)

			elevator = <-chanElevator
			OnFloorArrival(elevator, f, id, Timer, chanElevator)

		case <-Timer.C:
			fmt.Println("in time out")
			elevator = <-chanElevator
			OnDoorTimeOut(elevator, chanElevator)

		case a := <-chanObstr:
			fmt.Printf("Obstuction! %+v\n", a)
			elevator = <-chanElevator
			if a && elevator.Behaviour == EBdoorOpen {
				Timer.Reset(3 * time.Second)
				Timer.Stop()
				chanElevator <- elevator
			}
			if !a {
				Timer.Stop()
				OnDoorTimeOut(elevator, chanElevator)
			}

		case a := <-chanStop:
			fmt.Printf("stop button: %+v\n", a)
			SetStopLamp(true)
			for a {
				for f := 0; f < NumFloors; f++ {
					for b := ButtonType(0); b < 3; b++ {
						SetButtonLamp(b, f, false)
					}
				}
				time.Sleep(500 * time.Millisecond)
				for f := 0; f < NumFloors; f++ {
					for b := ButtonType(0); b < 3; b++ {
						SetButtonLamp(b, f, true)

					}
				}
				time.Sleep(500 * time.Millisecond)
			}
		}
		time.Sleep(10 * time.Millisecond)
	}

}
