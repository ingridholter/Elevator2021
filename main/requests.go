package main

import (
	. "main/config"
)

func RequestsAbove(elev ElevState) int {
	for f := elev.Floor + 1; f < NumFloors; f++ {
		for btn := 0; btn < NumButtons; btn++ {
			if elev.Requests[f][btn] {
				return 1
			}
		}
	}
	return 0
}

func RequestsBelow(elev ElevState) int {
	for f := 0; f < elev.Floor; f++ {
		for btn := 0; btn < NumButtons; btn++ {
			if elev.Requests[f][btn] {
				return 1
			}
		}
	}
	return 0
}

func RequestsChooseDirection(elev ElevState) MotorDirection {
	switch elev.Dirn {
	case MD_Up:
		if (RequestsAbove(elev)==1) {
			return MD_Up
		} else if (RequestsBelow(elev)==1) {
			return MD_Down
		} else {
			return MD_Stop
		}

	case MD_Down:
	case MD_Stop: /*....*/
	}
}

/*
BT_HallUp   ButtonType
	BT_HallDown
	BT_Cab
*/

func requestShouldStop(elev ElevState) int {
	switch elev.ElevState {
	case MD_Down:
		return (elev.Requests[elev.Floor][BT_HallDown] || elev.Requests[elev.Floor][BT_Cab] || !RequestsBelow(elev))

	case MD_Up:
		return (elev.requests[elev.floor][BT_HallUp] || elev.requests[elev.floor][BT_Cab] || !requestsAbove(elev))

	case MD_Stop: //do nothing?
	default:
	}
	return 1
}

//function for clearing orders??
