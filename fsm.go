package main

import (
	"time"
	. "fmt"
	. "./elevio"
	. "./config"
)

var elevator elevState //elevator state variable

func timerDoor() { //funker som en sleep i 3 sekunder

	timer := time.NewTimer(3*time.Second)
	<-timer.C
}

func onInitBetweenFloors(){
	SetMotorDirection(MD_Down)
	elevator.dirn = MD_Down
	elevator.behaviour = EBmoving
}

func setLights(elev elevState){
	for floor := 0; floor <_numFloors; floor++{
		for btn := 0; btn < _numButtons; btn++{
			SetButtonLamp(btn, floor,elev.requests[floor][btn]) //requests == 0/1 is this false/true?
		}
	}
}

func onRequestButtonPress(btnFloor int,btnType ButtonType){
	switch elevator.behaviour{
	case EBdoorOpen:
		if(elevator.floor==btnFloor){
			//start timer for door
		}else{
			elevator.requests[btnFloor][btnType] =1
		}
	case EBmoving:
		elevator.requests[btnFloor][btnType] =1
	case EBidle:
		if(elevator.floor ==btnFloor){
			SetDoorOpenLamp(true)
			//timer start
			elevator.behaviour = EBdoorOpen
		}else{
			elevator.requests[btnFloor][btnType] =1
			elevator.dirn = requests_chooseDirection(elevator)
			SetMotorDirection(elevator.dirn)
			elevator.behaviour = EBmoving
		}
	}
	setLights(elevator)
	//can print the state of elevator for debugg process
}

func onFloorArrival(newFloor int){
	//can print the new floor and the state of elevator

	elevator.floor= newFloor
	SetFloorIndicator(elevator.floor)
	switch elevator.behaviour {
	case EBmoving:
		if(requestShouldStop(elevator)){
			elevator.dirn = EBstop 
			SetDoorOpenLamp(true)
			//clear the order we have done elevator = request_clearAtCurrentFloor(elevator)
			//start timer for door
			setLights(elevator)
			elevator.behaviour = EBdoorOpen
		}
	} 
	//can print state 
}

func onFloorTimeOut(){
	//can print elevator state and function

	switch elevator.behaviour {
	case EBdoorOpen:
		elevator.dirn =requests_chooseDirection(elevator)
		elevator.SetDoorOpenLamp(false)
		elevator.dirn = elevator.dirn
		if(elevator.dirn == MD_Stop){
			elevator.behavoior = EBidle 
		} else{
			elevator.behaviour= EBmoving
		}
	}
	//print state?
}


