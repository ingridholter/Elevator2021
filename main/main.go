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
	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	peerTxEnable := make(chan bool)
	go peers.Transmitter(20008, id, peerTxEnable)
	go peers.Receiver(20008, peerUpdateCh)

	ElevStateMsgTx := make(chan ElevStateMsg)
	ElevStateMsgRx := make(chan ElevStateMsg)
	NewOrderMsgTx := make(chan NewOrderMsg)
	NewOrderMsgRx := make(chan NewOrderMsg)

	// ... and start the transmitter/receiver pair on some port
	// These functions can take any number of channels! It is also possible to
	//  start multiple transmitters/receivers on the same port.
	go bcast.Transmitter(20009, ElevStateMsgTx)
	go bcast.Receiver(20009, ElevStateMsgRx) //10.100.23.209
	go bcast.Transmitter(20007, NewOrderMsgTx)
	go bcast.Receiver(20007, NewOrderMsgRx)

	numFloors := 4
	Init(simport, numFloors)

	drvButtons := make(chan ButtonEvent)
	drvFloors := make(chan int)
	drvObstr := make(chan bool)
	drvStop := make(chan bool,1)
	go PollButtons(drvButtons)
	go PollFloorSensor(drvFloors)
	go PollObstructionSwitch(drvObstr)
	go PollStopButton(drvStop)

	chanNewOrder := make(chan ButtonEvent,1)
	chanElevator := make(chan ElevState, 1)
	chanElevatorArray := make(chan [NumElevators]ElevState, 1)
	chanElevatorLastMoved := make(chan map[int]time.Time, 1)
	lostId := make(chan int, NumElevators) //only one id lost a time
	chanLostElevators := make(chan [NumElevators]string, 1)
	lightsNoNetwork:=make(chan ElevState)

	go DrvElevator(id, chanNewOrder, drvFloors, drvObstr, drvStop, chanElevator, ElevStateMsgTx,lightsNoNetwork)
	go ElevatorObserver(id, ElevStateMsgRx, drvButtons, NewOrderMsgRx, NewOrderMsgTx, chanNewOrder, chanElevator, chanElevatorArray, ElevStateMsgTx, peerUpdateCh, lostId, chanElevatorLastMoved, chanLostElevators,lightsNoNetwork)

	select {}
}
