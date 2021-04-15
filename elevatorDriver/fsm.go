package elevatorDriver

import (
	"fmt"
	. "main/config"
	. "main/elevio"

	//. "main/lights"
	"time"
)

func OnInitBetweenFloors(elevator ElevState, chanElevator chan ElevState) {
	fmt.Println("drive down")
	elevator.Dir = MD_Down
	SetMotorDirection(elevator.Dir)
	elevator.Behaviour = EBmoving
	chanElevator <- elevator
}

func SetLights(elev ElevState) {
	for floor := 0; floor < NumFloors; floor++ {
		for btn := 0; btn < NumButtons-1; btn++ {

			SetButtonLamp(ButtonType(btn), floor, elev.Requests[floor][btn]) //requests == 0/1 is this false/true?
		}
	}
}

func OnRequestButtonPress(elevator ElevState, btnFloor int, btnType ButtonType, timer *time.Timer, chanElevator chan ElevState, lightsNoNetwork chan ElevState) {

	fmt.Println("adding to my requests")
	switch elevator.Behaviour {
	case EBdoorOpen:
		if elevator.Floor == btnFloor {
			fmt.Println("RESET O")
			timer.Reset(2 * time.Second)
			//så slette ordren?
			elevator = RequestClearAtCurrentFloor(elevator)

		} else {
			elevator.Requests[btnFloor][btnType] = true
		}
	case EBmoving:
		elevator.Requests[btnFloor][btnType] = true
	case EBidle:
		if elevator.Floor == btnFloor {
			SetDoorOpenLamp(true)
			fmt.Println("RESET I")
			timer.Reset(2 * time.Second) //TimerDoor() //timer start
			elevator.Behaviour = EBdoorOpen

		} else {
			elevator.Requests[btnFloor][btnType] = true
			elevator.Dir = RequestChooseDirection(elevator)
			SetMotorDirection(elevator.Dir)
			elevator.Behaviour = EBmoving
		}
	}
	lightsNoNetwork <- elevator
	chanElevator <- elevator
}

func OnFloorArrival(elevator ElevState, newFloor int, id string, timer *time.Timer, chanElevator chan ElevState, lightsNoNetwork chan ElevState) {
	//can print the new floor and the state of elevator
	//Id, _ := strconv.Atoi(id)

	elevator.Floor = newFloor
	SetFloorIndicator(elevator.Floor)
	fmt.Println("state:", elevator.Behaviour)

	switch elevator.Behaviour {
	case EBmoving:
		if RequestShouldStop(elevator) {
			SetMotorDirection(MD_Stop)
			fmt.Println("set open door lamp")
			SetDoorOpenLamp(true)

			elevator = RequestClearAtCurrentFloor(elevator)

			fmt.Println("RESET M")
			timer.Reset(2 * time.Second) //start timer for door

			elevator.Behaviour = EBdoorOpen
			lightsNoNetwork <- elevator
		}
	}

	chanElevator <- elevator

}

func OnDoorTimeOut(elevator ElevState, chanElevator chan ElevState) {
	//can print elevator state and function
	fmt.Println("IN ON FLOOR TIME OUT")
	switch elevator.Behaviour {
	case EBdoorOpen:

		elevator.Dir = RequestChooseDirection(elevator)
		fmt.Println("Close door")
		SetDoorOpenLamp(false)
		SetMotorDirection(elevator.Dir)
		if elevator.Dir == MD_Stop {

			elevator.Behaviour = EBidle
		} else {
			elevator.Behaviour = EBmoving
		}
	}
	chanElevator <- elevator
	fmt.Println("going out of on door time out")
}
