package costFunc


import (
	. "main/config"
	. "main/elevatorDriver"
	"strconv"
)

func bestElevator(eOld [NumElevators]ElevState, lostElevators [NumElevators]string) int {
	CostMap := make(map[int]int)
	for elevId := 0; elevId < NumElevators; elevId++ {
		if lostElevators[elevId] == "lost" {
			CostMap[elevId] = 99999999 //infinity
		} else {
			CostMap[elevId] = timeToIdle(eOld[elevId])
		}
	}
	minTime := CostMap[0]
	for _, value := range CostMap {

		if value < minTime {
			minTime = value
		}
	}
	
	for elevator, value := range CostMap {
		if value == minTime {
			return elevator
		}

	}
	return -1 
}

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
	const TRAVELTIME = 2500
	const DOOROPENTIME = 3000

	var duration int = 0
	switch e.Behaviour {
	case EBidle:
		e.Dir = RequestChooseDirection(e)
		if e.Dir == MD_Stop {
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

func checkDuplicate(Id int, allElevators [NumElevators]ElevState, button ButtonEvent, lostElevators [NumElevators] string) bool {

	temp := false
	for i, elevator := range allElevators {
		if lostElevators[i] == "found" {
			if button.Button != BT_Cab {
				temp = temp || elevator.Requests[button.Floor][button.Button]
			} else {
				return allElevators[Id].Requests[button.Floor][BT_Cab]
			}
		}
	}
	return temp
}

func NewOrderDistributer(eOld [NumElevators]ElevState, btn ButtonEvent, id string, lostElevators [NumElevators]string) NewOrderMsg {
	
	Id, _ := strconv.Atoi(id)

	f := btn.Floor
	btnType := btn.Button

	if checkDuplicate(Id, eOld, btn,lostElevators) {
		msgNoOne := NewOrderMsg{
			SenderId:   id,
			RecieverId: "duplicate",
			Button:     btn,
		}
		return msgNoOne
	}

	for i := 0; i < NumElevators; i++ {
		eOld[i].Requests[f][btnType] = true
	}

	bestElevatorId := bestElevator(eOld, lostElevators)

	if bestElevatorId == Id || btnType == 2 {

		msgMe := NewOrderMsg{
			SenderId:   id,
			RecieverId: id,
			Button:     btn,
		}
		return msgMe
	}

	msg := NewOrderMsg{
		SenderId:   id,
		RecieverId: strconv.Itoa(bestElevatorId),
		Button:     btn,
	}

	return msg
}
