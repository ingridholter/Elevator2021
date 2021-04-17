package elevatorDriver

import (
	. "main/config"
	. "main/elevio"
	"time"
)

func OnInitBetweenFloors(elevator ElevState, ElevStateCh chan ElevState) {
	elevator.Dir = MD_Down
	SetMotorDirection(elevator.Dir)
	elevator.Behaviour = EBmoving
	ElevStateCh <- elevator
}

func OnRequestButtonPress(elevator ElevState, btnFloor int, btnType ButtonType, timer *time.Timer, ElevStateCh chan ElevState, LightsOfflineCh chan ElevState) {

	switch elevator.Behaviour {
	case EBdoorOpen:
		if elevator.Floor == btnFloor {
			timer.Reset(DoorOpenTime * time.Second)
			elevator = RequestClearAtCurrentFloor(elevator)
		} else {
			elevator.Requests[btnFloor][btnType] = true
		}
	case EBmoving:
		elevator.Requests[btnFloor][btnType] = true
	case EBidle:
		if elevator.Floor == btnFloor {
			SetDoorOpenLamp(true)
			timer.Reset(DoorOpenTime * time.Second) 
			elevator.Behaviour = EBdoorOpen

		} else {
			elevator.Requests[btnFloor][btnType] = true
			elevator.Dir = RequestChooseDirection(elevator)
			SetMotorDirection(elevator.Dir)
			elevator.Behaviour = EBmoving
		}	
	}
	LightsOfflineCh <- elevator
	ElevStateCh <- elevator
}

func OnFloorArrival(elevator ElevState, newFloor int, id string, timer *time.Timer, ElevStateCh chan ElevState, LightsOfflineCh chan ElevState) {
	
	elevator.Floor = newFloor
	SetFloorIndicator(elevator.Floor)

	switch elevator.Behaviour {
	case EBmoving:
		if RequestShouldStop(elevator) {
			SetMotorDirection(MD_Stop)
			SetDoorOpenLamp(true)

			elevator = RequestClearAtCurrentFloor(elevator)

			timer.Reset(DoorOpenTime * time.Second)

			elevator.Behaviour = EBdoorOpen
			LightsOfflineCh <- elevator
		}
	}
	ElevStateCh <- elevator
}

func OnDoorTimeOut(elevator ElevState, ElevStateCh chan ElevState) {

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
	ElevStateCh <- elevator
}
