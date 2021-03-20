package costFunc

import (
	. "main/config"
	. "main/elevatorDriver"
)

//Får inn alle states fra elevator observer
//får inn knappetrykk
//skal gi ut ny R=4x(heis*3)

//returnerer 100,101,102 basert på heis som skal ta ordren
func bestElevator(eOld [NumElevators]ElevState) int {
	CostMap := make(map[int]int)
	for elevNum := 0; elevNum < NumElevators; elevNum++ {
		CostMap[elevNum] = timeToIdle(eOld[elevNum])
	}
	minTime := CostMap[0] //se på mer
	for _, value := range CostMap {

		if value < minTime {
			minTime = value
		}

	}
	for key, value := range CostMap {

		if value == minTime {
			return key + 100
		}

	}
	//return 0 if could not find best elevator
	return 0
}

//simualte clearing orders
func SimualtionRequestClearAtCurrentFloor(eOld ElevState) ElevState {
	elev := eOld
	elev.Requests[elev.Floor][BT_Cab] = false
	switch elev.Dir {
	case MD_Up:
		elev.Requests[elev.Floor][BT_HallUp] = false
		if !RequestsAbove(elev) {
			elev.Requests[elev.Floor][BT_HallDown] = false
		}
	case MD_Down:
		elev.Requests[elev.Floor][BT_HallDown] = false
		if !RequestsBelow(elev) {
			elev.Requests[elev.Floor][BT_HallUp] = false
		}
	case MD_Stop:
		elev.Requests[elev.Floor][BT_HallUp] = false
		elev.Requests[elev.Floor][BT_HallDown] = false
	default:
		elev.Requests[elev.Floor][BT_HallUp] = false
		elev.Requests[elev.Floor][BT_HallDown] = false
	}
	return elev
}

func timeToIdle(eOld ElevState) int {
	e := eOld
	var TRAVELTIME = 3   //hva skal travel time være?
	var DOOROPENTIME = 3 //hva skal door open time være?

	var duration int = 0

	switch e.Behaviour {
	case EBidle:
		e.Dir = RequestChooseDirection(e)
		if e.Dir == EBstop {
			return duration
		}
	case EBmoving:
		duration += TRAVELTIME / 2
		e.Floor += int(e.Dir)
	case EBdoorOpen:
		duration -= DOOROPENTIME / 2
	}

	for {
		if RequestShouldStop(e) {
			e = SimualtionRequestClearAtCurrentFloor(e)

			duration += DOOROPENTIME
			e.Dir = RequestChooseDirection(e)
			if e.Dir == MD_Stop {
				return duration
			}
		}
		e.Floor += int(e.Dir)
		duration += TRAVELTIME
	}
}

//buttonType is int now
func RequestMatrix(eOld [NumElevators]ElevState, btnType int, f int) [NumFloors][NumButtons * NumElevators]bool {
	//gives 100,101,102

	id := bestElevator(eOld)
	//makes sure that cab orders are taken by owners
	if btnType == 2 {
		id = PeerId
	}
	elevatorIndex := id - 100
	AllRequests[f][btnType+3*elevatorIndex] = true //rad, kolonne
	return AllRequests
}
