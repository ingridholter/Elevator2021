package ElevatorObserver

//module for keeping track over states in other elevator, handling message error and setting lights to "take order"

import (
	. "main/config"
	. "fmt"
	"strconv"
	. "main/elevatorDriver"
)




type PowerLoss struct {
	id       int
	Message  string
	elevator ElevState
}


//sette lys og ta ordre, det kommer fra motatt R matrise
func AcceptNewOrder(msg NewOrderMsg,id int,elevator ElevState) {
	senderId,_ := strconv.Atoi(msg.SenderId)
	recieverId,_ := strconv.Atoi(msg.RecieverId)
	if (id == senderId){
		//sjekke om noen andre vet om denne ordren? Nei tror det går fint
	}
	if (id == recieverId){

		elevator.Requests[msg.Button.Floor][msg.Button.Button] = true //legger til ordre i min lokale matrise
		
	}
}

//tror ikke det blir problem når jeg selv skal ta ordre og sette lys og så får feil
func SyncAllLights(allElevators [NumElevators]ElevState){

	for _,elevator:= range allElevators{
		SetLights(elevator)
	}

}


func UpdateElevStateArray(msg ElevStateMsg) {
	Println(msg)

	id,_ := strconv.Atoi(msg.SenderId)
	ElevStateArray[id] =  msg.Elevator
}



func elevatorActive(id int, peers []string) bool{
	for _, b:= range peers{
		a,_ := strconv.Atoi(b)
		if a == id{
			return true
		}
	}
	return false
}

//oppdatere elevStateArray slik at de peers som ikke er på nettverket har floor=-1 og ingen requests, tar inn p.Peers
func ActiveElevatorStates(peers []string){

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
	for i:=0; i<NumElevators; i++{
		if elevatorActive(i+100,peers){
			ActiveElevatorStates[i] = ElevStateArray[i]
		} else {
			ActiveElevatorStates[i] = err
		}
	}
	ElevStateArray = ActiveElevatorStates
}

