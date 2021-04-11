package elevatorDriver

import (
	"fmt"
	. "main/config"
	. "main/elevio"

	//. "main/lights"
	"strconv"
	"time"
)

func OnInitBetweenFloors() {
	SetMotorDirection(MD_Down)
	Elevator.Dir = MD_Down
	Elevator.Behaviour = EBmoving
}

func SetLights(elev ElevState) {
	for floor := 0; floor < NumFloors; floor++ {
		for btn := 0; btn < NumButtons-1; btn++ {

			SetButtonLamp(ButtonType(btn), floor, elev.Requests[floor][btn]) //requests == 0/1 is this false/true?
		}
	}
}

func OnRequestButtonPress(btnFloor int, btnType ButtonType, timer *time.Timer) {

	//timer.Reset(3 * time.Second)
	fmt.Println("adding to my requests")
	switch Elevator.Behaviour {
	case EBdoorOpen:
		if Elevator.Floor == btnFloor {
			fmt.Println("RESET O")
			timer.Reset(3 * time.Second) //TimerDoor() //start timer for door
			//så slette ordren?
			Elevator = RequestClearAtCurrentFloor(Elevator)
			//OnFloorTimeOut()

		} else {
			Elevator.Requests[btnFloor][btnType] = true
		}
	case EBmoving:
		Elevator.Requests[btnFloor][btnType] = true
	case EBidle:
		if Elevator.Floor == btnFloor {
			SetDoorOpenLamp(true)
			fmt.Println("RESET I")
			timer.Reset(3 * time.Second) //TimerDoor() //timer start
			Elevator.Behaviour = EBdoorOpen
			//OnFloorTimeOut()
		} else {
			Elevator.Requests[btnFloor][btnType] = true
			Elevator.Dir = RequestChooseDirection(Elevator)
			SetMotorDirection(Elevator.Dir)
			Elevator.Behaviour = EBmoving
		}
	}
	//fmt.Println("LIGHTS IN fsm")
	//SetLights(Elevator)
	//fmt.Println("Elevator ", Elevator)
	//can print the state of elevator for debugg process
}

func OnFloorArrival(newFloor int, id string, timer *time.Timer) {
	//can print the new floor and the state of elevator
	Id, _ := strconv.Atoi(id)

	Elevator.Floor = newFloor
	SetFloorIndicator(Elevator.Floor)

	switch Elevator.Behaviour {
	case EBmoving:
		if RequestShouldStop(Elevator) {
			SetMotorDirection(MD_Stop)
			fmt.Println("set open door lamp")
			SetDoorOpenLamp(true)

			Elevator = RequestClearAtCurrentFloor(Elevator)
			ElevStateArray[Id] = Elevator //vet ikke om denne er nødvendig

			//SyncAllLights(ElevStateArray, id)
			fmt.Println("RESET M")
			timer.Reset(3 * time.Second) //TimerDoor() //start timer for door
			//SetLights(Elevator) vet ikke om det går fint at den er med her eller ikke
			Elevator.Behaviour = EBdoorOpen
		}
	}
	//can print state
}

func OnFloorTimeOut() {
	//can print elevator state and function
	fmt.Println("IN ON FLOOR TIME OUT")
	switch Elevator.Behaviour {
	case EBdoorOpen:

		Elevator.Dir = RequestChooseDirection(Elevator)
		fmt.Println("Close door")
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
