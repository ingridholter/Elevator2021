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
	ElevStateMsgTx chan ElevStateMsg, peerUpdateCh <-chan peers.PeerUpdate, lostId chan int,
	chanElevatorLastMoved chan map[int]time.Time, chanLostElevators chan [NumElevators]string, lightsNoNetwork chan ElevState) {

	var elevatorArray [NumElevators]ElevState
	Id, _ := strconv.Atoi(id)

	for i := 0; i < NumElevators; i++ {
		elevatorArray[i] = ElevState{
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

	elevatorArray[Id] = elevator
	chanElevatorArray <- elevatorArray

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

			//antar at p.Lost kun har en id på seg
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
			//evnt gjøre noe med new mtp software kræsj

			//Ta lost sine cab orders og lagre dem.

			//når koblet på igjen: send msg om cab ordren til den som er new.

			//i en annen case: ta new sin cab order kolonne og or'e med lagret cab orders
			elevatorArray = <-chanElevatorArray

			DistibuteLostOrders(LostId, elevatorArray, NewOrderMsgTx, chanElevatorArray,lostElevators)

			elevatorArray = <-chanElevatorArray
			ActiveElevatorStates(id, p.New, elevatorArray, chanElevatorArray, lostElevators, chanLostElevators, NewOrderMsgTx,chanElevatorLastMoved)
			//fmt.Println("elevatorStateArray: ", elevatorArray)

		case m := <-ElevStateMsgRx:
			fmt.Println("Recieving state message")
			elevatorArray = <-chanElevatorArray
			lostElevators = <-chanLostElevators
			elevatorLastMoved = <-chanElevatorLastMoved
			UpdateTimerElevatorLost(id, m, elevatorLastMoved, chanElevatorLastMoved, elevatorArray, lostElevators, chanLostElevators, NewOrderMsgTx)

			UpdateElevStateArray(m, elevatorArray, chanElevatorArray)

			elevatorArray = <-chanElevatorArray
			SyncAllLights(elevatorArray, id,lostElevators)
			chanElevatorArray <- elevatorArray

			//fmt.Println("out of case")

		case b := <-ButtonPressed:
			fmt.Printf("Button pressed %+v\n", b)

			elevatorArray = <-chanElevatorArray
			fmt.Println("elevator: ", elevatorArray[Id])
			msg := NewOrderDistributer(elevatorArray, b.Button, b.Floor, id, lostElevators) //ny mld med hvem som skal ha ordre!
			fmt.Println("elevator after order distr: ", elevatorArray[Id])
			chanElevatorArray <- elevatorArray

			lostElevators = <-chanLostElevators
			chanLostElevators <- lostElevators

			if msg.RecieverId == id && lostElevators[Id] == id {
				//ta den selv, ved nettverksfeil
				fmt.Println("msg ", msg)
				elevatorArray = <-chanElevatorArray
				UpdateOrders(id, msg, elevatorArray, chanElevatorArray)
				elevatorArray = <-chanElevatorArray
				SyncAllLights(elevatorArray, id,lostElevators)
				chanElevatorArray <- elevatorArray
				chanNewOrder <- b
				fmt.Println("elevator after order updated: ", elevatorArray[Id])
			} else {
				NewOrderMsgTx <- msg
			}

		case o := <-NewOrderMsgRx:

			fmt.Println("New order msg: ", o.RecieverId)
			lostElevators = <-chanLostElevators
			chanLostElevators <- lostElevators

			
			
			if AcceptNewOrder(o, id,chanElevatorArray) { //is the new order for this elevator?
				
				fmt.Println("update")
				elevatorArray = <-chanElevatorArray
				UpdateOrders(id, o, elevatorArray, chanElevatorArray)
				elevatorArray = <-chanElevatorArray
				SyncAllLights(elevatorArray, id,lostElevators)
				chanElevatorArray <- elevatorArray
				fmt.Print("before new order")
				chanNewOrder <- o.Button
				fmt.Print("sendt to new order")
			}

		case Id := <-lostId:
			//alarm
			lostElevators = <-chanLostElevators
			chanLostElevators <- lostElevators

			fmt.Println("lostElevators: ", lostElevators)
			elevatorArray = <-chanElevatorArray
			DistibuteLostOrders(Id, elevatorArray, NewOrderMsgTx, chanElevatorArray,lostElevators)

		case <-checkElevatorLost.C:
			lostElevators = <-chanLostElevators

			CheckTimerElevatorLost(chanElevatorLastMoved, lostId, lostElevators, chanLostElevators)

		case e := <-lightsNoNetwork:
			//skrur av lys
			lostElevators = <-chanLostElevators
			chanLostElevators <- lostElevators
			fmt.Println("in case lights no network")
			if lostElevators[Id] == id {
				fmt.Println("turning off lights when order done")
				elevatorArray = <-chanElevatorArray
				elevatorArray[Id] = e

				SyncAllLights(elevatorArray, id,lostElevators)
				chanElevatorArray <- elevatorArray
			}
		}

	}
}
