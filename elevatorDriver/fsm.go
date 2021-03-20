package elevatorDriver

import (
	. "main/config"
	. "main/elevio"
	"time"
)

var elevator ElevState //elevator state variable

func TimerDoor() { //funker som en sleep i 3 sekunder

	timer := time.NewTimer(3 * time.Second)
	<-timer.C
}

func OnInitBetweenFloors() {
	SetMotorDirection(MD_Down)
	elevator.Dir = MD_Down
	elevator.Behaviour = EBmoving
}

func SetLights(elev ElevState) {
	for floor := 0; floor < NumFloors; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			SetButtonLamp(ButtonType(btn), floor, elev.Requests[floor][btn]) //requests == 0/1 is this false/true?
		}
	}
}

func OnRequestButtonPress(btnFloor int, btnType ButtonType) {
	switch elevator.Behaviour {
	case EBdoorOpen:
		if elevator.Floor == btnFloor {
			TimerDoor() //start timer for door
		} else {
			elevator.Requests[btnFloor][btnType] = true
		}
	case EBmoving:
		elevator.Requests[btnFloor][btnType] = true
	case EBidle:
		if elevator.Floor == btnFloor {
			SetDoorOpenLamp(true)
			TimerDoor() //timer start
			elevator.Behaviour = EBdoorOpen
		} else {
			elevator.Requests[btnFloor][btnType] = true
			elevator.Dir = RequestChooseDirection(elevator)
			SetMotorDirection(elevator.Dir)
			elevator.Behaviour = EBmoving
		}
	}
	SetLights(elevator)
	//can print the state of elevator for debugg process
}

func OnFloorArrival(newFloor int) {
	//can print the new floor and the state of elevator

	elevator.Floor = newFloor
	SetFloorIndicator(elevator.Floor)
	switch elevator.Behaviour {
	case EBmoving:
		if RequestShouldStop(elevator) {
			SetMotorDirection(MD_Stop)
			SetDoorOpenLamp(true)
			elevator = RequestClearAtCurrentFloor(elevator)
			TimerDoor() //start timer for door
			SetLights(elevator)
			elevator.Behaviour = EBdoorOpen
		}
	}
	//can print state
}

func OnFloorTimeOut() {
	//can print elevator state and function

	switch elevator.Behaviour {
	case EBdoorOpen:
		elevator.Dir = RequestChooseDirection(elevator)
		SetDoorOpenLamp(false)
		SetMotorDirection(elevator.Dir)
		if elevator.Dir == MD_Stop {
			elevator.Behaviour = EBidle
		} else {
			elevator.Behaviour = EBmoving
		}
	}
	//print state?
}
