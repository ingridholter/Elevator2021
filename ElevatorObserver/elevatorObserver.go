package ElevatorObserver

import (
	"fmt"
	"main/Network/peers"
	. "main/config"
	. "main/costFunc"
	. "main/elevio"
	"strconv"
	"time"
)

func ElevatorObserver(id string, ElevStateMsgRx <-chan ElevStateMsg, ButtonPressed <-chan ButtonEvent,
	NewOrderMsgRx <-chan NewOrderMsg, NewOrderMsgTx chan<- NewOrderMsg,
	chanNewOrder chan<- ButtonEvent, chanElevator chan ElevState, chanElevatorArray chan [NumElevators]ElevState,
	ElevStateMsgTx chan ElevStateMsg, peerUpdateCh <-chan peers.PeerUpdate, lostId chan int, chanElevatorLastMoved chan map[int]time.Time) {

	var elevatorArray [NumElevators]ElevState
	Id, _ := strconv.Atoi(id)

	for i := 0; i < NumElevators; i++ {
		elevatorArray[i] = ElevState{
			Floor:     -2,
			Dir:       MD_Stop,
			Behaviour: EBidle,
			Requests:  [NumFloors][NumButtons]bool{},
		}
	}
	elevator := ElevState{
		Floor:     GetFloor(),
		Dir:       MD_Stop,
		Behaviour: EBidle,
		Requests:  [NumFloors][NumButtons]bool{},
	}

	elevatorArray[Id] = elevator
	chanElevatorArray <- elevatorArray

	var lostElevators [NumElevators]string
	for j := range lostElevators {
		lostElevators[j] = "q"
	}

	checkElevatorLost := time.NewTicker(80 * time.Millisecond)

	elevatorLastMoved := make(map[int]time.Time)

	for i := 0; i < NumElevators; i++ {
		fmt.Println("in foor loop: ", i)
		elevatorLastMoved[i] = time.Now()
	}

	chanElevatorLastMoved <- elevatorLastMoved

	for {
		select {

		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

			//evnt gjøre noe med new mtp software kræsj

			//Ta lost sine cab orders og lagre dem.

			//når koblet på igjen: send msg om cab ordren til den som er new.

			//i en annen case: ta new sin cab order kolonne og or'e med lagret cab orders
			elevatorArray = <-chanElevatorArray
			ActiveElevatorStates(p.Peers, elevatorArray, chanElevatorArray)
			//fmt.Println("elevatorStateArray: ", elevatorArray)

		case m := <-ElevStateMsgRx:
			//fmt.Println("Recieving state message ")
			elevatorArray = <-chanElevatorArray

			UpdateTimerElevatorLost(id, m, chanElevatorLastMoved, elevatorArray, lostElevators)

			UpdateElevStateArray(m, elevatorArray, chanElevatorArray)

			elevatorArray = <-chanElevatorArray
			SyncAllLights(elevatorArray, id)
			chanElevatorArray <- elevatorArray

			//fmt.Println("out of case")

		case b := <-ButtonPressed:
			fmt.Printf("Button pressed %+v\n", b)

			elevatorArray = <-chanElevatorArray
			msg := NewOrderDistributer(elevatorArray, b.Button, b.Floor, id) //ny mld med hvem som skal ha ordre!
			chanElevatorArray <- elevatorArray
			NewOrderMsgTx <- msg

		case o := <-NewOrderMsgRx:
			fmt.Println("New order msg: ", o.RecieverId)
			if AcceptNewOrder(o, id) { //is the new order for this elevator?
				elevatorArray = <-chanElevatorArray
				fmt.Println("update")
				UpdateOrders(id, o, elevatorArray, chanElevatorArray)
				elevatorArray = <-chanElevatorArray
				SyncAllLights(elevatorArray, id)
				chanElevatorArray <- elevatorArray
				chanNewOrder <- o.Button
			}

		case Id := <-lostId:
			//alarm
			if Id == 1 {
				fmt.Println("LOST ID: ", Id)
				fmt.Println("lostElevators: ", lostElevators)
			}

			id := strconv.Itoa(Id)
			if lostElevators[Id] != id {
				lostElevators[Id] = id
				fmt.Println("restributing orders", Id)
				elevatorArray = <-chanElevatorArray
				DistibuteLostOrders(Id, elevatorArray, NewOrderMsgTx, chanElevatorArray)
			}

		case <-checkElevatorLost.C:
			CheckTimerElevatorLost(chanElevatorLastMoved, lostId)
		}

	}
}
