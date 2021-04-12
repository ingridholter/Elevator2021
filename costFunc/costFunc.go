package costFunc

import (
	"fmt"
	. "main/config"
	. "main/elevatorDriver"
	"strconv"
	"sync"
)

//Får inn alle states fra elevator observer
//får inn knappetrykk
//skal gi ut ny R=4x(heis*3)

//returnerer 100,101,102 basert på heis som skal ta ordren
func bestElevator(eOld [NumElevators]ElevState) int {
	CostMap := make(map[int]int)
	for elevNum := 0; elevNum < NumElevators; elevNum++ {
		if eOld[elevNum].Floor == -2 {
			CostMap[elevNum] = 99999999 //infinity
		} else {
			fmt.Println("elevatorNumber: ",elevNum)
			CostMap[elevNum] = timeToIdle(eOld[elevNum])
		}
	}
	minTime := CostMap[0] //se på mer
	for _, value := range CostMap {

		if value < minTime {
			minTime = value
			fmt.Println("min time: ", minTime)
		}

	}
	fmt.Println(CostMap)
	for key, value := range CostMap {
		fmt.Println("key ", key, "value ",value)
		if value == minTime {
			fmt.Println("key: ",key)
			return key
		}

	}
	//return -1 if could not find best elevator
	return -1
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
	const TRAVELTIME = 2500   //hva skal travel time være?
	const DOOROPENTIME = 3000 //hva skal door open time være?

	var duration int = 0
	fmt.Println("behaviour: ",e.Behaviour)
	switch e.Behaviour {
	case EBidle:
		e.Dir = RequestChooseDirection(e)
		if e.Dir == MD_Stop {
			return duration
		}
	case EBmoving:
		duration += TRAVELTIME / 2
		e.Floor += int(e.Dir)
		fmt.Println("duration: ", duration)
	case EBdoorOpen:
		duration -= DOOROPENTIME / 2
	}

	for {
		fmt.Println("elevator: ", e)
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

//send msg to the id that should take the order
func NewOrderDistributer(eOld [NumElevators]ElevState, btnType ButtonType, f int, id string, elevator ElevState, m *sync.Mutex) NewOrderMsg{

	//min id
	Id, _ := strconv.Atoi(id)

	bestElevatorId := bestElevator(eOld)

	b := ButtonEvent{Button: btnType, Floor: f}

	if ElevStateArray[bestElevatorId].Requests[f][btnType]{
		
		//da har vi ordren fra før
		msgNoOne := NewOrderMsg{
			SenderId:   id,
			RecieverId: "duplicate",
			Button:     b,
		}
		return  msgNoOne
	}

	//makes sure that cab orders are taken by owners
	if bestElevatorId == Id || btnType == 2 {

		//oppdatere egen request matrise og elevstatearray
		elevator.Requests[f][btnType] = true
		m.Lock()
		ElevStateArray[Id].Requests[f][btnType] = true
		m.Unlock()
		
		//send mld til meg seg
		msgMe := NewOrderMsg{
			SenderId:   id,
			RecieverId: id,
			Button:     b,
		}
		return msgMe
	}

	msg := NewOrderMsg{
		SenderId:   id,
		RecieverId: strconv.Itoa(bestElevatorId),
		Button:     b,
	}
	
	return msg
}
