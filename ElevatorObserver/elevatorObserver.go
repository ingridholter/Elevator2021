package ElevatorObserver

import (
	"fmt"
	. "main/config"
	. "main/costFunc"
	"main/Network/peers"
	"strconv"
)


func ElevatorObserver(id string,ElevStateMsgRx <-chan ElevStateMsg, ButtonPressed <-chan ButtonEvent,
	NewOrderMsgRx <-chan NewOrderMsg, NewOrderMsgTx chan<- NewOrderMsg,
	chanNewOrder chan<- ButtonEvent,chanElevator chan ElevState,chanElevatorArray chan [NumElevators]ElevState,
	ElevStateMsgTx chan ElevStateMsg, peerUpdateCh <-chan peers.PeerUpdate, lostId <-chan int){

	var elevatorArray[NumElevators]ElevState
	Id, _ := strconv.Atoi(id)

	for i:=0; i<NumElevators;i++{
		elevatorArray[i] =ElevState{
			Floor:    	-2,
			Dir:       MD_Stop,
			Behaviour: EBidle,
			Requests:  [NumFloors][NumButtons]bool{},
		}
	}
	elevatorArray[Id] = <-chanElevator
	chanElevatorArray <- elevatorArray


	for{
		select{
	
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

			//evnt gjøre noe med new mtp software kræsj

			//Ta lost sine cab orders og lagre dem. 

			//når koblet på igjen: send msg om cab ordren til den som er new.

			//i en annen case: ta new sin cab order kolonne og or'e med lagret cab orders
			elevatorArray= <-chanElevatorArray
			ActiveElevatorStates(p.Peers,elevatorArray,chanElevatorArray)
		

		case m:= <-ElevStateMsgRx:

			//TimerElevatorLost(r,id,lostId) //fikse denne senere
			elevatorArray= <-chanElevatorArray

			UpdateElevStateArray(m,elevatorArray,chanElevatorArray)
			
			elevatorArray= <-chanElevatorArray
			SyncAllLights(elevatorArray, id) 
			chanElevatorArray <- elevatorArray
		
		case b:= <-ButtonPressed:
			fmt.Printf("%+v\n", b)

			elevatorArray= <-chanElevatorArray
			msg := NewOrderDistributer(elevatorArray, b.Button, b.Floor, id) //ny mld med hvem som skal ha ordre!
			NewOrderMsgTx <- msg
			chanElevatorArray<-elevatorArray

		case o:= <-NewOrderMsgRx:

			if AcceptNewOrder(o, id) { //is the new order for this elevator?
				elevatorArray= <-chanElevatorArray
				SyncAllLights(elevatorArray, id)
				chanNewOrder <- o.Button
				chanElevatorArray<-elevatorArray
			}
		case id:= <-lostId:
			//alarm
			fmt.Println("LOST ID: ", id)
			elevatorArray= <-chanElevatorArray
			DistibuteLostOrders(id,elevatorArray,NewOrderMsgTx,chanElevatorArray)
		}
	}
}

