package elevatorDriver

import (
	. "main/config"
	. "main/elevio"
	"time"
)

func RunElevator(id string, NewOrderCh <-chan ButtonEvent, FloorCh <-chan int, ObstrCh <-chan bool,
	StopCh chan bool, ElevStateCh chan ElevState, ElevStateMsgTx chan<- ElevStateMsg, LightsOfflineCh chan ElevState) {

	elevator := ElevState{
		Floor:     GetFloor(),
		Dir:       MD_Stop,
		Behaviour: EBidle,
		Requests:  [NumFloors][NumButtons]bool{},
	}

	if Between {
		OnInitBetweenFloors(elevator, ElevStateCh)
	} else {
		ElevStateCh <- elevator
	}

	DoorTimer := time.NewTimer(DoorOpenTime * time.Second)
	DoorTimer.Stop()
	SendStateTicker := time.NewTicker(500 * time.Millisecond)
	
	
	for {

		select {
		case <-SendStateTicker.C:

			elevator = <-ElevStateCh

			ElevStateMsgTx <- ElevStateMsg{
				SenderId: id,
				Elevator: elevator,
			}

			ElevStateCh <- elevator

		case b := <-NewOrderCh: 
			
			elevator = <-ElevStateCh
			
			OnRequestButtonPress(elevator, b.Floor, b.Button, DoorTimer, ElevStateCh, LightsOfflineCh)

		case f := <-FloorCh: 

			elevator = <-ElevStateCh
			OnFloorArrival(elevator, f, id, DoorTimer, ElevStateCh, LightsOfflineCh)

		case <-DoorTimer.C:
			elevator = <-ElevStateCh
			OnDoorTimeOut(elevator, ElevStateCh)

		case obstr := <-ObstrCh:
			elevator = <-ElevStateCh
			ElevStateCh <- elevator
			if obstr && elevator.Behaviour == EBdoorOpen {
				DoorTimer.Reset(3 * time.Second)
				DoorTimer.Stop()
			}
			if !obstr {
				DoorTimer.Stop()
				elevator = <-ElevStateCh
				OnDoorTimeOut(elevator, ElevStateCh)
			}

		case stop := <-StopCh:
			SetStopLamp(true)
			if stop {
				elevator = <-ElevStateCh
				elevator.Dir = MD_Stop
				SetMotorDirection(elevator.Dir)
				ElevStateCh <- elevator
			} else {
				SetStopLamp(false)
				elevator = <-ElevStateCh
				if !(elevator.Behaviour == EBdoorOpen) {
					elevator.Dir = RequestChooseDirection(elevator)
					SetMotorDirection(elevator.Dir)
				}
				ElevStateCh <- elevator
			}
		}
	}

}
