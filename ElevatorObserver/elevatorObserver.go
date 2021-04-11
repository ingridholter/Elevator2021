package ElevatorObserver

//module for keeping track over states in other elevator, handling message error and setting lights to "take order"

import (
	"time"

	"fmt"
	. "main/config"
	. "main/costFunc"

	//"main/elevatorDriver"
	. "main/elevio"
	"strconv"
)

//Hvis fått mld om å ta ordre og den er din å ta, legg til i din requests.
func AcceptNewOrder(msg NewOrderMsg, id string, elevator ElevState) bool {

	if id == msg.RecieverId {

		//gjøre dette i fsm
		//elevator.Requests[msg.Button.Floor][msg.Button.Button] = true //legger til ordre i min lokale matrise
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

func UpdateElevStateArray(msg ElevStateMsg) {
	//fmt.Println("update Elev state array: ", msg.Elevator.Requests)

	id, _ := strconv.Atoi(msg.SenderId)

	ElevStateArray[id] = msg.Elevator
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
func ActiveElevatorStates(peers []string) {

	var ActiveElevatorStates [NumElevators]ElevState

	//error state
	var err = -2
	
	//0,1,2
	for i := 0; i < NumElevators; i++ {
		if elevatorActive(i, peers) {
			ActiveElevatorStates[i] = ElevStateArray[i]
		} else {
			ActiveElevatorStates[i].Floor = err
		}
	}
	ElevStateArray = ActiveElevatorStates
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


func TimerElevatorLost(msg ElevStateMsg,id string){
	//Id, _ := strconv.Atoi(id)
	SenderIdInt,_:=strconv.Atoi(msg.SenderId)

	oldElevState := ElevStateArray[SenderIdInt]
	newElevState := msg.Elevator

	//fmt.Println("any requests? ",AnyRequestsExist(oldElevState))

	//fmt.Println("old: ",oldElevState)
	//fmt.Println("new: ",newElevState)

	if oldElevState.Floor != newElevState.Floor ||  oldElevState.Dir != newElevState.Dir || oldElevState.Behaviour != newElevState.Behaviour{
		//fmt.Println("Resetting timer")
		ElevatorLastMoved[SenderIdInt].Reset(5*time.Second)
	}else if !AnyRequestsExist(oldElevState){
		//fmt.Println("dont have requests, resetting timer")
		ElevatorLastMoved[SenderIdInt].Reset(5*time.Second)
	}
	
}


func IdElevatorLost(LostId chan<- int){

	for{
	for i:= range ElevatorLastMoved{
		a :=ElevatorLastMoved[i]
	
		select{
			case <-a.C:
				LostId <- i
			}
		}
	}
}


//redistrubuere ordrene til en tapt heis. OBS IKKE CAB ORDERS!!
func DistibuteLostOrders(LostId int, lostOrders chan<- NewOrderMsg) {

	//lagrer lost id sine request i ny matrise
	LostRequests:= ElevStateArray[LostId].Requests

	//oppdate elevarray til at den er død
	ElevStateArray[LostId].Floor =  -2

	//slette ordrene til den heisen?

	//loop gjennom alle den dødes ordre og fordele de med newOrderDistributer.
	for floor := 0; floor < NumFloors; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			if LostRequests[floor][btn] {
				msg := NewOrderDistributer(ElevStateArray, ButtonType(btn), floor, strconv.Itoa(LostId), ElevStateArray[LostId])
				lostOrders <- msg
			}
		}
	}

}
