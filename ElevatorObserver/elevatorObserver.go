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

/*
	NewOrderCh := make(chan ButtonEvent,1)
	ElevStateCh := make(chan ElevState, 1)
	ElevStateArrayCh := make(chan [NumElevators]ElevState, 1)
	ElevLastMovedCh := make(chan map[int]time.Time, 1)
	LostIdCh := make(chan int, NumElevators)
	LostElevArrayCh := make(chan [NumElevators]string, 1)
	LightsOfflineCh:=make(chan ElevState)
*/


func ElevatorObserver(id string, ElevStateMsgRx <-chan ElevStateMsg, ButtonPressed <-chan ButtonEvent,
	NewOrderMsgRx <-chan NewOrderMsg, NewOrderMsgTx chan<- NewOrderMsg,
	NewOrderCh chan<- ButtonEvent, chanElevator chan ElevState, chanElevatorArray chan [NumElevators]ElevState,
	ElevStateMsgTx chan ElevStateMsg, peerUpdateCh <-chan peers.PeerUpdate, lostId chan int,
	chanElevatorLastMoved chan map[int]time.Time, chanLostElevators chan [NumElevators]string, lightsNoNetwork chan ElevState) {

	var elevStateArray [NumElevators]ElevState
	Id, _ := strconv.Atoi(id)

	for i := 0; i < NumElevators; i++ {
		elevStateArray[i] = ElevState{
			Floor:     GetFloor(),
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

	elevStateArray[Id] = elevator
	chanElevatorArray <- elevStateArray

	var lostElevators [NumElevators]string
	for j := range lostElevators {
		lostElevators[j] = strconv.Itoa(j)
	}
	chanLostElevators <- lostElevators

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

			var LostId int = -2
			lostElevators = <-chanLostElevators

			
			for _,lostId := range p.Lost {
				if len(p.Lost) ==1{
					LostId,_ = strconv.Atoi(lostId)
					lostElevators[LostId] = strconv.Itoa(LostId)
				}else{
					elev,_ := strconv.Atoi(lostId)
					lostElevators[elev] = lostId
				}
			}
			chanLostElevators <- lostElevators
			fmt.Println("lost id: ", LostId)
			
			elevStateArray = <-chanElevatorArray

			DistibuteLostOrders(Id, NewOrderCh,LostId, elevStateArray, NewOrderMsgTx, chanElevatorArray,lostElevators)

			elevStateArray = <-chanElevatorArray
			ActiveElevatorStates(id, p.New, elevStateArray, chanElevatorArray, lostElevators, chanLostElevators, NewOrderMsgTx,chanElevatorLastMoved)
			

		case m := <-ElevStateMsgRx:
			fmt.Println("Recieving state message")
			elevStateArray = <-chanElevatorArray
			lostElevators = <-chanLostElevators
			elevatorLastMoved = <-chanElevatorLastMoved
			UpdateTimerElevatorLost(id, m, elevatorLastMoved, chanElevatorLastMoved, elevStateArray, lostElevators, chanLostElevators, NewOrderMsgTx)

			UpdateElevStateArray(m, elevStateArray, chanElevatorArray)

			elevStateArray = <-chanElevatorArray
			SyncAllLights(elevStateArray, id,lostElevators)
			chanElevatorArray <-elevStateArray


		case b := <-ButtonPressed:
			fmt.Printf("Button pressed %+v\n", b)

			elevStateArray= <-chanElevatorArray
			fmt.Println("elevator: ", elevStateArray[Id])
			msg := NewOrderDistributer(elevStateArray, b.Button, b.Floor, id, lostElevators) //ny mld med hvem som skal ha ordre!
			fmt.Println("elevator after order distr: ", elevStateArray[Id])
			chanElevatorArray <- elevStateArray

			lostElevators = <-chanLostElevators
			chanLostElevators <- lostElevators

			//ved nettverksfeil så skal vi tenne lys og ta ordre som blir trykket på			
			if msg.RecieverId == id && lostElevators[Id] == id {
				//ta den selv, ved nettverksfeil
				fmt.Println("msg ", msg)
				elevStateArray = <-chanElevatorArray
				UpdateOrders(id, msg, elevStateArray, chanElevatorArray)
				elevStateArray = <-chanElevatorArray
				SyncAllLights(elevStateArray, id,lostElevators)
				chanElevatorArray <- elevStateArray
				NewOrderCh <- b
				fmt.Println("elevator after order updated: ", elevStateArray[Id])
			} else {
				NewOrderMsgTx <- msg
			}

		case o := <-NewOrderMsgRx:

			fmt.Println("New order msg: ", o.RecieverId,"should accept? ", AcceptNewOrder(o, id,chanElevatorArray))
			lostElevators = <-chanLostElevators
			chanLostElevators <- lostElevators

			
			
			if AcceptNewOrder(o, id,chanElevatorArray) { //is the new order for this elevator?
				
				fmt.Println("update")
				elevStateArray = <-chanElevatorArray
				UpdateOrders(id, o, elevStateArray, chanElevatorArray)
				elevStateArray = <-chanElevatorArray
				SyncAllLights(elevStateArray, id,lostElevators)
				chanElevatorArray <- elevStateArray
				fmt.Print("before new order")
				NewOrderCh <- o.Button
				fmt.Print("sendt to new order")
			}

		case L := <-lostId:
			//alarm
			lostElevators = <-chanLostElevators
			chanLostElevators <- lostElevators

			fmt.Println("lostElevators: ", lostElevators)
			elevStateArray = <-chanElevatorArray
			DistibuteLostOrders(Id, NewOrderCh,L, elevStateArray, NewOrderMsgTx, chanElevatorArray,lostElevators)

		case <-checkElevatorLost.C:
			lostElevators = <-chanLostElevators

			CheckTimerElevatorLost(chanElevatorLastMoved, lostId, lostElevators, chanLostElevators)

		case e := <-lightsNoNetwork:
			//skrur av lys
			lostElevators = <-chanLostElevators
			chanLostElevators <- lostElevators
			fmt.Println("in case lights no network")
			if lostElevators[Id] == id{ // && changeState //sånn at vi ikke setter lys ved motorstopp
				fmt.Println("turning off lights when order done")
				elevStateArray = <-chanElevatorArray
				elevStateArray[Id] = e

				SyncAllLights(elevStateArray, id,lostElevators)
				chanElevatorArray <- elevStateArray
			}
		}

	}
}
