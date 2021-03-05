package main

func requestsAbove(elev elevState) int{
	for f := elev.floor+1; f <_numFloors; f++{
		for btn := 0; btn < _numButtons; btn++{
			if(elev.requests[f][btn]){
				return 1
			}
		}
	}
	return 0
}

func requestsBelow(elev elevState) int{
	for f := 0; f <elev.floor; f++{
		for btn := 0; btn < _numButtons; btn++{
			if(elev.requests[f][btn]){
				return 1
			}
		}
	}
	return 0
}

func requestsChooseDirection(elev elevState) dirn{
	switch (elev.elevState){
	case MD_Up:
		if requestsAbove(elev) {
			return MD_Up
		} else if requestsBelow(elev) {
			return MD_Down
		}else { return MD_Stop}

	case MD_Down: 
	case MD_Stop:/*....*/		
	}
}
/*
BT_HallUp   ButtonType 
	BT_HallDown          
	BT_Cab 
*/

func requestShouldStop(elev elevState) int{
	switch elev.elevState {
	case MD_Down:
		return (elev.requests[elev.floor][BT_HallDown] || elev.requests[elev.floor][BT_Cab] || !requestsBelow(elev))
			
	case MD_Up:
		return (elev.requests[elev.floor][BT_HallUp] || elev.requests[elev.floor][BT_Cab] || !requestsAbove(elev))
			
	case MD_Stop: //do nothing?
	default:
	}
	return 1;
}

//function for clearing orders??