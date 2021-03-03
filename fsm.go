package fsm

import (
	. "fmt"
	. "./elevio" //how to import this correctly?
)

type elevBehaviour int


const (
	EBmoving elevBehaviour = 0
	EBstop = 1
	EBdoorOpen = 2
	EBidle = 3
)

type elevState struct {
	floor int
	dirn elevio.MotorDirection
	behaviour elevBehaviour
	requests[_numFloors][_numButtons] int
}

func onInitBetweenFloors(){
	SetMotorDirection(MD_Down)
	elevState.dirn = MD_Down
	elevState.behaviour = EBmoving
}

func setLights(elev elevState){
	for floor := 0; i <_numFloors; floor++{
		for btn := 0; btn < _numButtons; btn++{
			SetButtonLamp(btn, floor,true) //only turn on lights not off
		}
	}


func onRequestButtonPress(btnFloor int,btnType elevio.ButtonType){
	switch elevState.behaviour{
	case EBdoorOpen:
		if(elevState.floor==btnFloor){
			//start timer for door
		}else{
			elevState.requests[btnFloor][btnType] =1
		}
	case EBmoving:
		elevState.requests[btnFloor][btnType] =1
	case EBidle:
		if(elevState.floor ==btnFloor){
			SetDoorOpenLamp(true)
			//timer start
			elevState.behaviour = EBdoorOpen
		}else{
			elevState.requests[btnFloor][btnType] =1
			elevState.dirn = //need logic for deciding direction
			SetMotorDirection(elevState.dirn)
			elevState.behaviour = EBmoving
		}
	}
	//sett the correct lights
	//can print the state of elevator for debugg process
}

func onFloorArrival(newFloor int){
	//can print the new floor and the state of elevator

	elevator.floor= newFloor
	elevio.SetFloorIndicator(elevator.floor)
	switch elevator.behaviour {
	case EBmoving:
		if(){//if elevator is requested 
			elevState.dirn = MD_Down
			elevio.SetDoorOpenLamp(true)
			//clear the order we have done
			//start timer for door
			//set lights for the orders
			elevState.behaviour = EBdoorOpen
		}
	} 
	//can print state 
}

func onFloorTimeOut(){
	//can print elevator state and function

	switch elevator.behaviour {
	case EBdoorOpen:
		elevState.dirn =//choose direction based on reqests
		elevio.SetDoorOpenLamp(false)
		elevState.dirn = elevator.dirn
		if(elevator.dirn == MD_Stop){
			elevState.behavoior = EBidle 
		} else{
			elevState.behaviour= EBmoving
		}
	}
	//print state?
}

