package ElevatorObserver

//module for keeping track over states in other elevator, handling message error and setting lights to "take order"

import (
	"fmt"
	. "main/config"

	//. "main/costFunc"
	"main/elevatorDriver"
	. "main/elevio"
	"strconv"
)

//Hvis fått mld om å ta ordre og den er din å ta, legg til i din requests.
func AcceptNewOrder(msg NewOrderMsg, id string, elevator ElevState) bool {
	//senderId, _ := strconv.Atoi(msg.SenderId)
	//recieverId, _ := strconv.Atoi(msg.RecieverId)
	if id == msg.SenderId {
		//sjekke om noen andre vet om denne ordren? Nei tror det går fint
		return true
	}
	if id == msg.RecieverId {

		//gjøre dette i fsm
		//elevator.Requests[msg.Button.Floor][msg.Button.Button] = true //legger til ordre i min lokale matrise
		return true
	}
	return false
}

//tror ikke det blir problem når jeg selv skal ta ordre og sette lys og så får feil. OBS CAB lights!!
func SyncAllLights(allElevators [NumElevators]ElevState, id string) {

	Id, _ := strconv.Atoi(id)
	//fmt.Println(allElevators)
	for floor := 0; floor < NumFloors; floor++ {
		SetButtonLamp(BT_Cab, floor, allElevators[Id].Requests[floor][2]) //this is for cab orders
	}

	for _, elevator := range allElevators {
		if elevator.Floor != -1 {
			elevatorDriver.SetLights(elevator)
			fmt.Println("active elevator")
		}

		//fmt.Println("setting lights in elevatorObs")
	}

}

func UpdateElevStateArray(msg ElevStateMsg) {
	fmt.Println("update Elev state array: ", msg)

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
	var err = ElevState{
		Floor:     -1,
		Dir:       MD_Stop,
		Behaviour: EBmoving,
		Requests: [4][3]bool{
			{false, false, false},
			{false, false, false},
			{false, false, false},
			{false, false, false},
		},
	}
	//0,1,2
	for i := 0; i < NumElevators; i++ {
		if elevatorActive(i, peers) {
			ActiveElevatorStates[i] = ElevStateArray[i]
		} else {
			ActiveElevatorStates[i] = err
		}
	}
	ElevStateArray = ActiveElevatorStates
}

//redistrubuere ordrene til en tapt heis. OBS IKKE CAB ORDERS!!
func DistibuteLostOrders(lost []string) {
	for _, id := range lost {
		Id, _ := strconv.Atoi(id)
		for floor := 0; floor < NumFloors; floor++ {
			for btn := 0; btn < NumButtons; btn++ {
				if ElevStateArray[Id].Requests[floor][btn] {
					//msg := NewOrderDistributer(ElevStateArray, ButtonType(btn), floor, id, Elevator)
					//send <- msg //send to channel :)
				}
			}
		}
	}
}
