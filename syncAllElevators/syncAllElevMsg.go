package syncAllElevators

import (
	. "main/config"
	. "main/elevio"
	"strconv"
	"time"
)

func AcceptNewOrder(msg NewOrderMsg, id string, ElevStateArrayCh chan [NumElevators]ElevState) bool {
	Id,_:=strconv.Atoi(id)

	elevatorArray:= <-ElevStateArrayCh
	ElevStateArrayCh<-elevatorArray
	
	orderExist:=elevatorArray[Id].Requests[msg.Button.Floor][msg.Button.Button]

	if id == msg.RecieverId {
		if !orderExist{
			return true
		}
	}
	return false
}

func SyncAllLights(allElevators [NumElevators]ElevState, Id int,lostElevators [NumElevators]string) {

	lightsMatrix := [4][3]bool{
		{false, false, false},
		{false, false, false},
		{false, false, false},
		{false, false, false},
	}

	for floor := 0; floor < NumFloors; floor++ {
		SetButtonLamp(BT_Cab, floor, allElevators[Id].Requests[floor][2]) //lights for cab orders
	}


	for index, elevator := range allElevators {
		if lostElevators[index] == "found" {

			for floor := 0; floor < NumFloors; floor++ {
				for btn := 0; btn < NumButtons-1; btn++ {
					lightsMatrix[floor][btn] = lightsMatrix[floor][btn] || elevator.Requests[floor][btn]
					SetButtonLamp(ButtonType(btn), floor, lightsMatrix[floor][btn])
				}
			}
			
		}else if Id == index {
			for floor := 0; floor < NumFloors; floor++ {
				for btn := 0; btn < NumButtons; btn++ {
					SetButtonLamp(ButtonType(btn), floor, elevator.Requests[floor][btn])
				}
			}
		} 
	}
}

func UpdateElevStateArray(msg ElevStateMsg, allElevators [NumElevators]ElevState, ElevStateArrayCh chan [NumElevators]ElevState) {

	Id, _ := strconv.Atoi(msg.SenderId)

	allElevators[Id] =msg.Elevator

	ElevStateArrayCh <- allElevators
}

func SyncMyOrders(Id int, newMsg NewOrderMsg, allElevators [NumElevators]ElevState, ElevStateArrayCh chan [NumElevators]ElevState) {

	allElevators[Id].Requests[newMsg.Button.Floor][newMsg.Button.Button] = true

	ElevStateArrayCh <- allElevators
}


func SyncNewElevator(id string, new string, allElevators [NumElevators]ElevState, ElevStateArrayCh chan [NumElevators]ElevState, lostElevators [NumElevators]string,
	 LostElevArrayCh chan [NumElevators]string, NewOrderMsgTx chan<- NewOrderMsg, ElevLastMovedCh chan map[int]time.Time) {

	lostElevators = <-LostElevArrayCh
	
	if len(new) != 0{
		New,_ := strconv.Atoi(new)	
		lostElevators[New] = "found"
		elevatorLastMoved:= <-ElevLastMovedCh
		elevatorLastMoved[New] = time.Now()
		ElevLastMovedCh<-elevatorLastMoved
		for f := 0; f < NumFloors; f++ {
			if allElevators[New].Requests[f][BT_Cab] {
				msg := NewOrderMsg{
						SenderId:   id,
						RecieverId: new,
						Button:     ButtonEvent{Floor: f, Button: BT_Cab},
				}
				for i:=0;i<10;i++{
					NewOrderMsgTx <- msg
				}
			}
		}
	} 

	ElevStateArrayCh <- allElevators
	LostElevArrayCh <- lostElevators
}