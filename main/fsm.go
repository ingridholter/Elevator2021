package main

import (
	"time"
	. "main/elevio"
	. "main/config"
)

var elevator ElevState //elevator state variable

func timerDoor() { //funker som en sleep i 3 sekunder

	timer := time.NewTimer(3*time.Second)
	<-timer.C
}

func onInitBetweenFloors(){
	SetMotorDirection(MD_Down)
	elevator.Dirn = MD_Down
	elevator.Behaviour = EBmoving
}

func setLights(elev ElevState){
	for floor := 0; floor < NumFloors; floor++{
		for btn := 0; btn < NumButtons; btn++{
			SetButtonLamp(btn, floor,elev.Requests[floor][btn]) //requests == 0/1 is this false/true?
		}
	}
}

func onRequestButtonPress(btnFloor int,btnType ButtonType){
	switch elevator.Behaviour{
	case EBdoorOpen:
		if(elevator.Floor==btnFloor){
			//start timer for door
		}else{
			elevator.Requests[btnFloor][btnType] =true
		}
	case EBmoving:
		elevator.Requests[btnFloor][btnType] =true
	case EBidle:
		if(elevator.Floor ==btnFloor){
			SetDoorOpenLamp(true)
			//timer start
			elevator.Behaviour = EBdoorOpen
		}else{
			elevator.Requests[btnFloor][btnType] =true
			elevator.Dirn = Requests_chooseDirection(elevator)
			SetMotorDirection(elevator.Dirn)
			elevator.Behaviour = EBmoving
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


