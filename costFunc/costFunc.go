package costFunc

import (
	//. "../config"
	//c "github.com/ingridholter/Elevator2021/config"
	. "main/config"
)

//Får inn alle states fra elevator observer
//får inn knappetrykk
//skal gi ut ny R=4x(heis*3)

//returnerer 100,101,102 basert på heis som skal ta ordren
func bestElevator(eOld [NumElevators]ElevState, b Buttontype, f int) int {
	var CostMap map[int]int
	for elevNum := 0; elevNum < NumElevators; elevNum++ {
		CostMap[elevNum] = timeToServeRequest(eOld[elevNum], b, f)
	}
	temp := CostMap[0] //se på mer
	for key, value := range CostMap {

		if value < temp {
			temp = value
		}

	}
	return temp.key + 100
}

func timeToServeRequest(eOld elevState, b Buttontype, f int) int {
	var e Elevator = eOld
	e.request[f][b] = 1

	var arrivedAtRequest = 0

	var duration int = 0

	switch e.behaviour {
	case EBidle:
		e.dirn = requestChooseDirection(e)
		if e.dirn == EBstop {
			return duration
		}
	case EBmoving:
		duration += TRAVELTIME / 2
		e.floor += e.dirn
	case EBdoorOpen:
		duration -= DOOROPENTIME / 2
	}

	for {
		if requestsShouldStop(e) {
			e = requestClearAtCurrentFloor(e)
			if inner_b == b && inner_f == f { //////USIKKER: Erstatter funksjonskallet ifEqual fra D-koden
				arrivedAtRequest = 1
			}
			if arrivedAtRequest {
				return duration
			}
			duration += DOOROPENTIME
			e.dirn = requestChooseDirection
		}
		e.floor += e.direction
		duration = TRAVELTIME
	}
}

func RequestMatrix(eOld elevStateArray, b Buttontype, f int) AllRequests {
	//gives 100,101,102
	id := bestElevator(eOld, b, f)
	elevatorIndex := id - 100
	AllRequests[f][b+3*elevatorIndex] = true //rad, kolonne
	return AllRequests
}
