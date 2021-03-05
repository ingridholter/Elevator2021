package costFunc

func bestElevator() int {
	var costMap map[int]int //
	for elevNum:=0; i<NUMELEVATORS; elevNum++{
		Costmap[elevNum] = timeToServeRequest(eOld elevState, b Buttontype, f int)
	} 
	temp := costcostMap[0]
	for key, value := range costMap {
		
		if (value < temp){
			temp = value
		}
		
	}
	return temp.key + 100
}

func timeToServeRequest(eOld elevState, b Buttontype, f int) int{
	var e Elevator = eOld
	e.request[f][b] = 1

	var arrivedAtRequest = 0
		
	var duration int = 0

	switch(e.behaviour){
	case EBidle:
			e.dirn = requestChooseDirection(e)
			if (e.dirn == EBstop){
				return duration
			}
	case EBmoving:
		duration += TRAVELTIME/2
		e.floor += e.dirn
	case EBdoorOpen:
		duration -= DOOROPENTIME/2


	}


	

	for{
		if(requestsShouldStop(e)){
			e = requestClearAtCurrentFloor(e)
			if(inner_b == b && inner_f == f){ //////USIKKER: Erstatter funksjonskallet ifEqual fra D-koden
				arrivedAtRequest = 1
				}
			if(arrivedAtRequest){
				return duration
			}
			duration += DOOROPENTIME
			e.dirn = requestChooseDirection
		}
		e.floor += e.direction
		duration = TRAVELTIME
	}
}



    
