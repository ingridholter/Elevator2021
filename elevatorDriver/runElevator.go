package elevatorDriver

import (
	"fmt"
	. "main/config"
	. "main/elevio"
	"time"
)

func DrvElevator(id string, chanNewOrder <-chan ButtonEvent, chanFloors <-chan int, chanObstr <-chan bool,
	chanStop chan bool, chanElevator chan ElevState, ElevStateMsgTx chan<- ElevStateMsg, lightsNoNetwork chan ElevState) {

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
	DoorTimer := time.NewTimer(3 * time.Second)
	DoorTimer.Stop()
	SendStateTicker := time.NewTicker(500 * time.Millisecond)
	
	
	for {

		select {
		case <-SendStateTicker.C:

			elevator = <-chanElevator

			ElevStateMsgTx <- ElevStateMsg{
				SenderId: id,
				Elevator: elevator,
			}

			chanElevator <- elevator

		case b := <-chanNewOrder: //SyncAllLights(ElevStateArray, id)
			fmt.Println("lagt til i min ordre kanal: ", id)
			elevator = <-chanElevator
			
			OnRequestButtonPress(elevator, b.Floor, b.Button, DoorTimer, chanElevator, lightsNoNetwork)

		case f := <-chanFloors: //SyncAllLights(ElevStateArray, id)
			fmt.Println("In case Floor: ", f)

			elevator = <-chanElevator
			OnFloorArrival(elevator, f, id, DoorTimer, chanElevator, lightsNoNetwork)

		case <-DoorTimer.C:
			fmt.Println("in time out")
			elevator = <-chanElevator
			OnDoorTimeOut(elevator, chanElevator)

		case a := <-chanObstr:
			fmt.Printf("Obstuction! %+v\n", a)
			elevator = <-chanElevator
			chanElevator <- elevator
			if a && elevator.Behaviour == EBdoorOpen {
				DoorTimer.Reset(3 * time.Second)
				DoorTimer.Stop()
			}
			if !a {
				DoorTimer.Stop()
				elevator = <-chanElevator
				OnDoorTimeOut(elevator, chanElevator)
			}

		case a := <-chanStop:
			fmt.Printf("stop button: %+v\n", a)
			SetStopLamp(true)
			if a {
				elevator = <-chanElevator
				elevator.Dir = MD_Stop
				SetMotorDirection(elevator.Dir)
				chanElevator <- elevator
			} else {
				SetStopLamp(false)
				elevator = <-chanElevator
				if !(elevator.Behaviour == EBdoorOpen) {
					elevator.Dir = RequestChooseDirection(elevator)
					SetMotorDirection(elevator.Dir)
				}
				chanElevator <- elevator
			}
		}
		time.Sleep(10 * time.Millisecond)
	}

}
