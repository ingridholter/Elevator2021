package elevatorDriver

import (
	. "main/config"
	"fmt"
)

func RequestsAbove(elev ElevState) bool {
	for f := elev.Floor + 1; f < NumFloors; f++ {
		for btn := 0; btn < NumButtons; btn++ {
			fmt.Println("Requests above: ", elev.Requests[f][btn])
			if elev.Requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func RequestsBelow(elev ElevState) bool {
	for f := 0; f < elev.Floor; f++ {
		for btn := 0; btn < NumButtons; btn++ {
			fmt.Println("Requests below: ", elev.Requests[f][btn])
			if elev.Requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func RequestChooseDirection(elev ElevState) MotorDirection {
	switch elev.Dir {
	case MD_Up:
		if RequestsAbove(elev) {
			return MD_Up
		} else if RequestsBelow(elev) {
			return MD_Down
		} else {
			return MD_Stop
		}

	case MD_Down:
		if RequestsBelow(elev) {
			return MD_Down
		} else if RequestsAbove(elev) {
			return MD_Up
		} else {
			return MD_Stop
		}
	case MD_Stop:
		if RequestsBelow(elev) {
			return MD_Down
		} else if RequestsAbove(elev) {
			return MD_Up
		} else {
			return MD_Stop
		}
	default:
		return MD_Stop
	}
}

/*
BT_HallUp   ButtonType
	BT_HallDown
	BT_Cab
*/

func RequestShouldStop(elev ElevState) bool {
	switch elev.Dir {
	case MD_Down:
		return (elev.Requests[elev.Floor][BT_HallDown] || elev.Requests[elev.Floor][BT_Cab] || !RequestsBelow(elev))

	case MD_Up:
		fmt.Println("case MD_UP: ", (elev.Requests[elev.Floor][BT_HallUp] || elev.Requests[elev.Floor][BT_Cab] || !RequestsAbove(elev)))
		return (elev.Requests[elev.Floor][BT_HallUp] || elev.Requests[elev.Floor][BT_Cab] || !RequestsAbove(elev))

	case MD_Stop: //do nothing?

	}
	return true
}

//metod 2 for clearing requests, only clearing in same direction
func RequestClearAtCurrentFloor(elev ElevState) ElevState {
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
