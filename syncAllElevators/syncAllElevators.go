package syncAllElevators

import (
	"fmt"
	"main/Network/peers"
	. "main/config"
	. "main/costFunc"
	. "main/elevio"
	"strconv"
	"time"
)

func SyncElevators(id string, ElevStateMsgRx <-chan ElevStateMsg, ButtonPressed <-chan ButtonEvent,
	NewOrderMsgRx <-chan NewOrderMsg, NewOrderMsgTx chan<- NewOrderMsg,
	NewOrderCh chan<- ButtonEvent, ElevStateArrayCh chan [NumElevators]ElevState,
	ElevStateMsgTx chan ElevStateMsg, peerUpdateCh <-chan peers.PeerUpdate, LostIdCh chan int,
	ElevLastMovedCh chan map[int]time.Time, LostElevArrayCh chan [NumElevators]string, LightsOfflineCh chan ElevState) {

	Id, _ := strconv.Atoi(id)	

	var elevStateArray [NumElevators]ElevState

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
	ElevStateArrayCh <- elevStateArray

	var lostElevArray [NumElevators]string

	for j := range lostElevArray {
		lostElevArray[j] = "lost"
	}
	LostElevArrayCh <- lostElevArray

	
	checkElevLost := time.NewTicker(80 * time.Millisecond)
	elevLastMovedMap := make(map[int]time.Time)

	for i := 0; i < NumElevators; i++ {
		elevLastMovedMap[i] = time.Now()
	}

	ElevLastMovedCh <- elevLastMovedMap

	for {
		select {

		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

			var oneLostId int = -1
			lostElevArray = <-LostElevArrayCh

			
			for _,lostId := range p.Lost {
				if len(p.Lost) ==1{
					oneLostId,_ = strconv.Atoi(lostId)
					lostElevArray[oneLostId] = "lost"
				}else{
					elev,_ := strconv.Atoi(lostId)
					lostElevArray[elev] = "lost"
				}
			}

			LostElevArrayCh <- lostElevArray
			
			elevStateArray = <-ElevStateArrayCh

			DistibuteLostOrders(Id, NewOrderCh,oneLostId, elevStateArray, NewOrderMsgTx, ElevStateArrayCh,lostElevArray)

			elevStateArray = <-ElevStateArrayCh

			SyncNewElevator(id, p.New, elevStateArray, ElevStateArrayCh, lostElevArray, LostElevArrayCh, NewOrderMsgTx,ElevLastMovedCh)
			

		case eMsg := <-ElevStateMsgRx:
			
			elevStateArray = <-ElevStateArrayCh
			
			UpdateTimerElevatorLost(eMsg, ElevLastMovedCh, elevStateArray, LostElevArrayCh, NewOrderMsgTx)

			UpdateElevStateArray(eMsg, elevStateArray, ElevStateArrayCh)

			elevStateArray = <-ElevStateArrayCh
			SyncAllLights(elevStateArray, Id,lostElevArray)
			ElevStateArrayCh <-elevStateArray


		case btn := <-ButtonPressed:
			

			elevStateArray= <-ElevStateArrayCh
			newMsg := NewOrderDistributer(elevStateArray, btn, id, lostElevArray)
			
			ElevStateArrayCh <- elevStateArray

			lostElevArray = <-LostElevArrayCh
			LostElevArrayCh <- lostElevArray
		
			if newMsg.RecieverId == id && lostElevArray[Id] == "lost" {	
				elevStateArray = <-ElevStateArrayCh
				SyncMyOrders(Id, newMsg, elevStateArray, ElevStateArrayCh)
				elevStateArray = <-ElevStateArrayCh
				SyncAllLights(elevStateArray, Id,lostElevArray)
				ElevStateArrayCh <- elevStateArray
				NewOrderCh <- btn
			} else {
				NewOrderMsgTx <- newMsg
			}

		case oMsg := <-NewOrderMsgRx:
		
			lostElevArray = <-LostElevArrayCh
			LostElevArrayCh <- lostElevArray

			if AcceptNewOrder(oMsg, id,ElevStateArrayCh) { 
				elevStateArray = <-ElevStateArrayCh
				SyncMyOrders(Id, oMsg, elevStateArray, ElevStateArrayCh)
				elevStateArray = <-ElevStateArrayCh
				ElevStateArrayCh <- elevStateArray
				SyncAllLights(elevStateArray, Id,lostElevArray)
				NewOrderCh <- oMsg.Button
			}

		case lostId := <-LostIdCh:

			lostElevArray = <-LostElevArrayCh
			LostElevArrayCh <- lostElevArray

			elevStateArray = <-ElevStateArrayCh
			DistibuteLostOrders(Id, NewOrderCh,lostId, elevStateArray, NewOrderMsgTx, ElevStateArrayCh,lostElevArray)

		case <-checkElevLost.C:
			lostElevArray = <-LostElevArrayCh

			CheckTimerElevLost(ElevLastMovedCh, LostIdCh, lostElevArray, LostElevArrayCh)

		case elev := <-LightsOfflineCh:
			
			lostElevArray = <-LostElevArrayCh
			LostElevArrayCh <- lostElevArray
			
			if lostElevArray[Id] == "lost"{ 
				elevStateArray = <-ElevStateArrayCh
				elevStateArray[Id] = elev

				SyncAllLights(elevStateArray, Id,lostElevArray)
				ElevStateArrayCh <- elevStateArray
			}
		}

	}
}
