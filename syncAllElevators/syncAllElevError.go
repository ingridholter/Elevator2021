package syncAllElevators

import (
	. "main/config"
	. "main/costFunc"
	"strconv"
	"time"
)

func anyRequestsExist(elevator ElevState) bool {
	for floor := 0; floor < NumFloors; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			if elevator.Requests[floor][btn] {
				return true
			}
		}
	}
	return false
}

func UpdateTimerElevatorLost(msg ElevStateMsg,ElevLastMovedCh chan map[int]time.Time, allElevators [NumElevators]ElevState, 
	LostElevArrayCh chan [NumElevators]string, NewOrderMsgTx chan<- NewOrderMsg) {
	
	elevLastMovedMap := <-ElevLastMovedCh
	lostElevArray := <-LostElevArrayCh

	SenderIdInt, _ := strconv.Atoi(msg.SenderId)
	newElevState := msg.Elevator

	oldElevState := allElevators[SenderIdInt]
	changedState := oldElevState.Floor != newElevState.Floor || oldElevState.Dir != newElevState.Dir || oldElevState.Behaviour != newElevState.Behaviour
	

	for lostId, _ := range lostElevArray  {
		
		if strconv.Itoa(lostId) == msg.SenderId && changedState { 
			lostElevArray[lostId] = "found"
		}

	}
	
	if changedState{
		elevLastMovedMap[SenderIdInt] = time.Now()
		
	}
	if !anyRequestsExist(oldElevState) {
		elevLastMovedMap[SenderIdInt] = time.Now()
		}

	ElevLastMovedCh <- elevLastMovedMap
	LostElevArrayCh <- lostElevArray 
}


func CheckTimerElevLost(ElevLastMovedCh chan map[int]time.Time, lostId chan int, lostElevators [NumElevators]string, LostElevArrayCh chan [NumElevators]string) {
	elevatorLastMoved := <-ElevLastMovedCh

	ElevLastMovedCh <- elevatorLastMoved

	for id, T := range elevatorLastMoved {

		if int(time.Now().Second())-int(T.Second()) > NoMotorTimeLimit {
			
			if lostElevators[id] != "lost" {
				lostElevators[id] = "lost"
				lostId <- id
			}
		}
	}
	LostElevArrayCh <- lostElevators
}

func DistibuteLostOrders(Id int, chanNewOrder chan<- ButtonEvent,oneLostId int, allElevators [NumElevators]ElevState, NewOrderMsgTx chan<- NewOrderMsg, 
	ElevStateArrayCh chan [NumElevators]ElevState, lostElevators [NumElevators]string) {

	if oneLostId == -1 {
		ElevStateArrayCh <- allElevators
		return
	}

	lostRequests := allElevators[oneLostId].Requests

	for floor := 0; floor < NumFloors; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			if lostRequests[floor][btn] && btn != BT_Cab{

				b:= ButtonEvent{Floor: floor, Button: ButtonType(btn)}
				msg := NewOrderDistributer(allElevators,b, strconv.Itoa(oneLostId),lostElevators)
					
				Reciever,_:=strconv.Atoi(msg.RecieverId)
				if Reciever == Id{
					chanNewOrder<-msg.Button
				}					
			}
		}
	}	
	ElevStateArrayCh <- allElevators
}