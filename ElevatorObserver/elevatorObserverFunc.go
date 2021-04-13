package ElevatorObserver

//module for keeping track over states in other elevator, handling message error and setting lights to "take order"

import (
	"time"
	"fmt"
	. "main/config"
	. "main/costFunc"
	"strconv"
	. "main/elevio"

)

//Hvis fått mld om å ta ordre og den er din å ta, legg til i din requests.
func AcceptNewOrder(msg NewOrderMsg, id string) bool {
	
	if id == msg.RecieverId {

		fmt.Println("my order: ", msg.Button)
		return true
	}
	return false
}

//tror ikke det blir problem når jeg selv skal ta ordre og sette lys og så får feil.
func SyncAllLights(allElevators [NumElevators]ElevState, id string) {

	Id, _ := strconv.Atoi(id)

	for floor := 0; floor < NumFloors; floor++ {
		SetButtonLamp(BT_Cab, floor, allElevators[Id].Requests[floor][2]) //this is for cab orders
	}

	lightsMatrix := [4][3]bool{
		{false, false, false},
		{false, false, false},
		{false, false, false},
		{false, false, false},
	}

	for index, elevator := range allElevators {
		if elevator.Floor != -2 {

			//får en lights matrix with all lights for all elevators
			for floor := 0; floor < NumFloors; floor++ {
				for btn := 0; btn < NumButtons; btn++ {
					lightsMatrix[floor][btn] = lightsMatrix[floor][btn] || elevator.Requests[floor][btn]
				}
			}

			for floor := 0; floor < NumFloors; floor++ {
				for btn := 0; btn < NumButtons-1; btn++ {
					SetButtonLamp(ButtonType(btn), floor, lightsMatrix[floor][btn])
				}
			}
			//network loss, if i am off the network is my floor -2? så vet ikke om dette er nødvendig?
		} else if Id == index {
			for floor := 0; floor < NumFloors; floor++ {
				for btn := 0; btn < NumButtons; btn++ {
					SetButtonLamp(ButtonType(btn), floor, elevator.Requests[floor][btn])
				}
			}
		}
	}
}

func UpdateElevStateArray(msg ElevStateMsg,allElevators [NumElevators]ElevState, chanElevatorArray chan [NumElevators]ElevState) {

	id, _ := strconv.Atoi(msg.SenderId)
	
	allElevators[id] = msg.Elevator

	chanElevatorArray<-allElevators
}


func elevatorActive(id int, peers []string) bool {
	for _, b := range peers {
		a, _ := strconv.Atoi(b)
		if a == id {
			return true
		}
	}
	return false
}

//oppdatere elevStateArray slik at de peers som ikke er på nettverket har floor=-1 og ingen requests, tar inn p.Peers
func ActiveElevatorStates(peers []string,allElevators [NumElevators]ElevState, chanElevatorArray chan [NumElevators]ElevState) {

	var activeElevatorStates [NumElevators]ElevState

	//error state
	var err = -2
	
	//0,1,2
	for i := 0; i < NumElevators; i++ {
		if elevatorActive(i, peers) {
			activeElevatorStates[i] = ElevStateArray[i]
		} else {
			activeElevatorStates[i].Floor = err
		}
	}
	chanElevatorArray <-activeElevatorStates
}

func AnyRequestsExist(elevator ElevState) bool{

	for floor := 0; floor < NumFloors; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			if elevator.Requests[floor][btn] {
				return true
			}
		}
	}
	return false
}


func TimerElevatorLost(msg chan ElevStateMsg,id string,LostId chan<- int,allElevators [NumElevators]ElevState){
	//Id, _ := strconv.Atoi(id)
	var newMsg ElevStateMsg
	for{
		
		select{
		case r:= <- msg:
			
			newMsg= r
			
		}
		
		SenderIdInt,_:=strconv.Atoi(newMsg.SenderId)
		newElevState:=newMsg.Elevator

	
		oldElevState := ElevStateArray[SenderIdInt]
	

		//fmt.Println("senderid: ", SenderIdInt)
		//fmt.Println("any requests? ",AnyRequestsExist(oldElevState))

		//fmt.Println("old: ",oldElevState)
		//fmt.Println("new: ",newElevState)
		if oldElevState.Floor != newElevState.Floor ||  oldElevState.Dir != newElevState.Dir || oldElevState.Behaviour != newElevState.Behaviour{
			//fmt.Println("Resetting timer")
			ElevatorLastMoved[SenderIdInt]=time.Now()
		}
		if !AnyRequestsExist(oldElevState){
			//fmt.Println(SenderIdInt,"dont have requests, resetting timer")
			ElevatorLastMoved[SenderIdInt]=time.Now()
		}

		if int(time.Now().Second())-int(ElevatorLastMoved[SenderIdInt].Second())>7 && AnyRequestsExist(oldElevState){
			LostId <- SenderIdInt
		}

		time.Sleep(100*time.Millisecond)
	}
}

/*
func IdElevatorLost(LostId chan<- int){

	for{
		fmt.Println("in idElevatorLost")
	for i:= range ElevatorLastMoved{
		a :=ElevatorLastMoved[i]
		fmt.Println("time since last moved: ",a)
		
		select{
			case <-a.C:
				LostId <- i
			}
		}
		time.Sleep(100*time.Millisecond)
	}
}
*/

//redistrubuere ordrene til en tapt heis. OBS IKKE CAB ORDERS!!
func DistibuteLostOrders(LostId int,allElevators [NumElevators]ElevState,NewOrderMsgTx chan<- NewOrderMsg,chanElevatorArray chan [NumElevators]ElevState) {

	//lagrer lost id sine request i ny matrise
	LostRequests:= allElevators[LostId].Requests

	//oppdate elevarray til at den er død
	
	allElevators[LostId].Floor =  -2
	

	//slette ordrene til den heisen?

	//loop gjennom alle den dødes ordre og fordele de med newOrderDistributer.
	for floor := 0; floor < NumFloors; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			if LostRequests[floor][btn] {
				msg:=NewOrderDistributer(allElevators, ButtonType(btn), floor, strconv.Itoa(LostId))
				NewOrderMsgTx <- msg
			}
		}
	}
	chanElevatorArray<-allElevators
}
