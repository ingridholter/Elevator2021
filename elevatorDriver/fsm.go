package elevatorDriver

import (
	"fmt"
	. "main/config"
	. "main/elevio"
	"time"
)

func TimerDoor() { //funker som en sleep i 3 sekunder

	timer := time.NewTimer(3 * time.Second)
	<-timer.C
}

func OnInitBetweenFloors() {
	SetMotorDirection(MD_Down)
	Elevator.Dir = MD_Down
	Elevator.Behaviour = EBmoving
}

func SetLights(elev ElevState) {
	for floor := 0; floor < NumFloors; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			if elev.Requests[floor][btn] {
				fmt.Println(true)
			}
			SetButtonLamp(ButtonType(btn), floor, elev.Requests[floor][btn]) //requests == 0/1 is this false/true?
		}
	}
}

func OnRequestButtonPress(btnFloor int, btnType ButtonType) {
	switch Elevator.Behaviour {
	case EBdoorOpen:
		if Elevator.Floor == btnFloor {
			TimerDoor() //start timer for door
		} else {
			Elevator.Requests[btnFloor][btnType] = true
		}
	case EBmoving:
		Elevator.Requests[btnFloor][btnType] = true
	case EBidle:
		if Elevator.Floor == btnFloor {
			SetDoorOpenLamp(true)
			TimerDoor() //timer start
			Elevator.Behaviour = EBdoorOpen
		} else {
			Elevator.Requests[btnFloor][btnType] = true
			Elevator.Dir = RequestChooseDirection(Elevator)
			SetMotorDirection(Elevator.Dir)
			Elevator.Behaviour = EBmoving
		}
	}
	//SetLights(Elevator)
	//fmt.Println("Elevator ", Elevator)
	//can print the state of elevator for debugg process
}

func OnFloorArrival(newFloor int) {
	//can print the new floor and the state of elevator

	Elevator.Floor = newFloor
	SetFloorIndicator(Elevator.Floor)
	switch Elevator.Behaviour {
	case EBmoving:
		if RequestShouldStop(Elevator) {
			SetMotorDirection(MD_Stop)
			SetDoorOpenLamp(true)
			Elevator = RequestClearAtCurrentFloor(Elevator)
			TimerDoor() //start timer for door
			SetLights(Elevator)
			Elevator.Behaviour = EBdoorOpen
		}
	}
	//can print state
}

func OnFloorTimeOut() {
	//can print elevator state and function

	switch Elevator.Behaviour {
	case EBdoorOpen:
		Elevator.Dir = RequestChooseDirection(Elevator)
		SetDoorOpenLamp(false)
		SetMotorDirection(Elevator.Dir)
		if Elevator.Dir == MD_Stop {
			Elevator.Behaviour = EBidle
		} else {
			Elevator.Behaviour = EBmoving
		}
	}
	//print state?
}
