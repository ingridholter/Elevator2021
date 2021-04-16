package main

import (
	"flag"
	"fmt"
	. "main/ElevatorObserver"
	"main/Network/bcast"
	"main/Network/peers"
	. "main/config"
	. "main/elevatorDriver"
	. "main/elevio"
	"os"
	"time"
)

func main() {
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")

	if id == "" {
		localIP, err := peers.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}

	var simport string
	flag.StringVar(&simport, "simport", "", "simport for this terminal")
	flag.Parse()

	fmt.Println(id)
	if simport == "" {
		simport = "localhost:15657"
	}

	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)
	go peers.Transmitter(20008, id, peerTxEnable)
	go peers.Receiver(20008, peerUpdateCh)

	ElevStateMsgTx := make(chan ElevStateMsg)
	ElevStateMsgRx := make(chan ElevStateMsg)
	NewOrderMsgTx := make(chan NewOrderMsg)
	NewOrderMsgRx := make(chan NewOrderMsg)

	go bcast.Transmitter(20009, ElevStateMsgTx)
	go bcast.Receiver(20009, ElevStateMsgRx) 
	go bcast.Transmitter(20007, NewOrderMsgTx)
	go bcast.Receiver(20007, NewOrderMsgRx)

	Init(simport, NumFloors)

	drvButtons := make(chan ButtonEvent)
	drvFloors := make(chan int)
	drvObstr := make(chan bool)
	drvStop := make(chan bool,1)
	go PollButtons(drvButtons)
	go PollFloorSensor(drvFloors)
	go PollObstructionSwitch(drvObstr)
	go PollStopButton(drvStop)

	NewOrderCh := make(chan ButtonEvent,1)
	ElevStateCh := make(chan ElevState, 1)
	ElevStateArrayCh := make(chan [NumElevators]ElevState, 1)
	ElevLastMovedCh := make(chan map[int]time.Time, 1)
	LostIdCh := make(chan int, NumElevators)
	LostElevArrayCh := make(chan [NumElevators]string, 1)
	LightsOfflineCh:=make(chan ElevState)

	go DrvElevator(id, NewOrderCh, drvFloors, drvObstr, drvStop, ElevStateCh, ElevStateMsgTx,LightsOfflineCh)
	go ElevatorObserver(id, ElevStateMsgRx, drvButtons, NewOrderMsgRx, NewOrderMsgTx, NewOrderCh, ElevStateCh, ElevStateArrayCh, ElevStateMsgTx, peerUpdateCh, LostIdCh, ElevLastMovedCh, LostElevArrayCh,LightsOfflineCh)

	select {}
}
