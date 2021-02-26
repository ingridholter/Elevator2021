package fsm

import {
	"fmt"
	"elevio"
}


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
	requests[N_FLOORS][N_BUTTONS] int
}

func onInitBetweenFloors(){
	SetMotorDirection(MD_Down)
	elevState.dirn = MD_Down
	elevState.behaviour = EBmoving
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

}