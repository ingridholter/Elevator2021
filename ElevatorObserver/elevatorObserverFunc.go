package ElevatorObserver

//module for keeping track over states in other elevator, handling message error and setting lights to "take order"

import (
	"fmt"
	. "main/config"
	. "main/costFunc"
	. "main/elevio"
	"strconv"
	"time"
)

//Hvis fått mld om å ta ordre og den er din å ta, legg til i din requests.
func AcceptNewOrder(msg NewOrderMsg, id string, chanElevatorArray chan [NumElevators]ElevState) bool {
	elevatorArray:= <-chanElevatorArray
	chanElevatorArray<-elevatorArray

	Id,_:=strconv.Atoi(id)
	if id == msg.RecieverId {
		if !elevatorArray[Id].Requests[msg.Button.Floor][msg.Button.Button]{
		fmt.Println("my order: ", msg.Button)
		return true
		}
	}
	return false
}


//tror ikke det blir problem når jeg selv skal ta ordre og sette lys og så får feil.
func SyncAllLights(allElevators [NumElevators]ElevState, id string,lostElevators [NumElevators]string) {

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
		if lostElevators[index] == "!L" {

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

func UpdateElevStateArray(msg ElevStateMsg, allElevators [NumElevators]ElevState, chanElevatorArray chan [NumElevators]ElevState) {

	id, _ := strconv.Atoi(msg.SenderId)

	
	allElevators[id] =msg.Elevator

	chanElevatorArray <- allElevators
}

func UpdateOrders(id string, newMsg NewOrderMsg, allElevators [NumElevators]ElevState, chanElevatorArray chan [NumElevators]ElevState) {

	Id, _ := strconv.Atoi(id)

	allElevators[Id].Requests[newMsg.Button.Floor][newMsg.Button.Button] = true

	chanElevatorArray <- allElevators

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
func ActiveElevatorStates(id string, new string, allElevators [NumElevators]ElevState, chanElevatorArray chan [NumElevators]ElevState, lostElevators [NumElevators]string, chanLostElevators chan [NumElevators]string, NewOrderMsgTx chan<- NewOrderMsg, chanElevatorLastMoved chan map[int]time.Time) {


	lostElevators = <-chanLostElevators
	
	//0,1,2
	
	if len(new) != 0{
		fmt.Println("ine new if: ", new)
		New,_ := strconv.Atoi(new)	
		lostElevators[New] = "!L"
		elevatorLastMoved:= <-chanElevatorLastMoved
		elevatorLastMoved[New] = time.Now()
		chanElevatorLastMoved<-elevatorLastMoved
		for f := 0; f < NumFloors; f++ {
			if allElevators[New].Requests[f][BT_Cab] {
				msg := NewOrderMsg{
						SenderId:   id,
						RecieverId: new,
						Button:     ButtonEvent{Floor: f, Button: BT_Cab},
				}
				fmt.Println("adding msg to chan")
				for i:=0;i<10;i++{
					fmt.Println("send cab order! ", i)
					NewOrderMsgTx <- msg
				}
			}
		}
	} 

	chanElevatorArray <- allElevators
	chanLostElevators <- lostElevators
}

func AnyRequestsExist(elevator ElevState) bool {

	for floor := 0; floor < NumFloors; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			if elevator.Requests[floor][btn] {
				return true
			}
		}
	}
	return false
}

func UpdateTimerElevatorLost(id string, msg ElevStateMsg, elevatorLastMoved map[int]time.Time, chanElevLastMoved chan map[int]time.Time, allElevators [NumElevators]ElevState, lostElevators [NumElevators]string, chanLostElevators chan [NumElevators]string, NewOrderMsgTx chan<- NewOrderMsg) {
	//Id, _ := strconv.Atoi(id)

	SenderIdInt, _ := strconv.Atoi(msg.SenderId)
	newElevState := msg.Elevator

	oldElevState := allElevators[SenderIdInt]
	changedState := oldElevState.Floor != newElevState.Floor || oldElevState.Dir != newElevState.Dir || oldElevState.Behaviour != newElevState.Behaviour
	//changedState := oldElevState.Dir != newElevState.Dir || oldElevState.Behaviour != newElevState.Behaviour

	for index, lostId := range lostElevators {
		
		if lostId == msg.SenderId && changedState {
			fmt.Println("endret state, ikke motorstopp: ", msg.SenderId)
			lostElevators[index] = "!L"
		}
	}
	fmt.Println("lostelevator in update timer: ", lostElevators)

	//fmt.Println("senderid: ", SenderIdInt)
	//fmt.Println("any requests? ",AnyRequestsExist(oldElevState))

	//fmt.Println("old: ",oldElevState)
	//fmt.Println("new: ",newElevState)

	if changedState{
		fmt.Println(SenderIdInt, "Resetting timer to ", SenderIdInt)
		elevatorLastMoved[SenderIdInt] = time.Now()
	}
	if !AnyRequestsExist(oldElevState) {
		fmt.Println(SenderIdInt, "dont have requests, resetting timer")
		elevatorLastMoved[SenderIdInt] = time.Now()
	}
	chanElevLastMoved <- elevatorLastMoved
	chanLostElevators <- lostElevators
}

func CheckTimerElevatorLost(elevLastMoved chan map[int]time.Time, lostId chan int, lostElevators [NumElevators]string, chanLostElevators chan [NumElevators]string) {
	elevatorLastMoved := <-elevLastMoved

	elevLastMoved <- elevatorLastMoved

	for i, t := range elevatorLastMoved {

		if int(time.Now().Second())-int(t.Second()) > 7 {
			id := strconv.Itoa(i)
			fmt.Println("elevator lost: ", i)
			if lostElevators[i] != id {
				fmt.Println("in if")
				lostElevators[i] = id
				fmt.Println("lost elevator array:", lostElevators[i])
				lostId <- i
			}
		}
	}
	chanLostElevators <- lostElevators
	//fmt.Println("lostelevators in checktimer: ", lostElevators)
}

//redistrubuere ordrene til en tapt heis. OBS IKKE CAB ORDERS!!
func DistibuteLostOrders(Id int, chanNewOrder chan<- ButtonEvent,LostId int, allElevators [NumElevators]ElevState, NewOrderMsgTx chan<- NewOrderMsg, chanElevatorArray chan [NumElevators]ElevState, lostElevators [NumElevators]string) {

	if LostId == -2 {
		chanElevatorArray <- allElevators
		return
	}
	//lagrer lost id sine request i ny matrise
	lostRequests := allElevators[LostId].Requests
	fmt.Println("lostrequests: ", lostRequests)
	//oppdate elevarray til at den er død

	fmt.Println("lost elevators: ",lostElevators)

	//slette ordrene til den heisen?

	//loop gjennom alle den dødes ordre og fordele de med newOrderDistributer.
	for floor := 0; floor < NumFloors; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			if lostRequests[floor][btn] {
				if lostRequests[floor][BT_Cab] {
					//b := ButtonEvent{Button: BT_Cab, Floor: floor}
					//chanCabOrders <- b
				} else {
					fmt.Println("DISTRIBUTE ORDER")
					msg := NewOrderDistributer(allElevators, ButtonType(btn), floor, strconv.Itoa(LostId),lostElevators)
					fmt.Println("before newordermsgTx")
					//bare legg det til i arrayet med en gang
					Reciever,_:=strconv.Atoi(msg.RecieverId)
					if Reciever == Id{
						chanNewOrder<-msg.Button
						fmt.Println("my lost order, adding to my chan")
					}
					//allElevators[Reciever].Requests[floor][btn]=true
					fmt.Println("add to all elevators: ", allElevators[Reciever])
					
					//NewOrderMsgTx <- msg 
					
					fmt.Println("after newordermsgTx")
				}
			}
		}
	}
	fmt.Println("add to all elevators")
	chanElevatorArray <- allElevators
	fmt.Println("out of func")
}
