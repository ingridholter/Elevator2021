package main

import (
	"sync"
	"flag"
	"fmt"
	. "main/ElevatorObserver"
	"main/Network/bcast"
	"main/Network/peers"
	. "main/config"
	. "main/costFunc"
	. "main/elevatorDriver"
	. "main/elevio"
	"os"
	"time"
)

func main() {
	var _mtx sync.Mutex
	var _mtx2 sync.Mutex
	var _mtx3 sync.Mutex
	var _mtx4 sync.Mutex

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
	fmt.Println("init state: ",Elevator.Behaviour)
	
	drv_buttons := make(chan ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
	go PollButtons(drv_buttons)
	go PollFloorSensor(drv_floors)
	go PollObstructionSwitch(drv_obstr)
	go PollStopButton(drv_stop)

	chan_timer := make(chan bool,1)
	go Timer(chan_timer)

	Timer := time.NewTimer(2 * time.Second)
	Timer.Stop()

	SendStateTicker := time.NewTicker(1*time.Second)

	
	for i:= range ElevatorLastMoved{
		ElevatorLastMoved[i] =  time.Now()
	}

	
	lostOrders := make(chan NewOrderMsg)
	//go DistibuteLostOrders(id,powerloss)

	//drive down if between floors
	

	elevstate := ElevStateMsg{
		SenderId: id,
		Message:  "State Update",
		Elevator: Elevator,
	}
	lostId := make(chan int) //only one id lost a time
	//go IdElevatorLost(lostId)

	//newestElevStateMsg := make(chan ElevStateMsg)

	
	
	if Between {
		fmt.Println("on init between floors")
		OnInitBetweenFloors()
	}
	
	/*
	go func() {
		for {
			_mtx.Lock()
			ElevStateMsgTx <- elevstate
			_mtx.Unlock()
			time.Sleep(1 * time.Second)
		
		}
		
	}()
	*/

	for {
		select {
		
		case <-SendStateTicker.C:
			//send my state every 1 seconds
			_mtx.Lock()
			ElevStateMsgTx <- elevstate
			_mtx.Unlock()

		case p := <-peerUpdateCh:

			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

			//evnt gjøre noe med new mtp software kræsj

			//Ta lost sine cab orders og lagre dem. 

			//når koblet på igjen: send msg om cab ordren til den som er new.

			//i en annen case: ta new sin cab order kolonne og or'e med lagret cab orders

			ActiveElevatorStates(p.Peers,&_mtx4)
			//fmt.Println("update active e: ", ElevStateArray)
			
		case l := <- lostOrders:
			//update active elevators
				_mtx2.Lock() 
				NewOrderMsgTx <- l
				_mtx2.Unlock()
			//redistribuere alle ordre
/*
			//sende disse til de som skal ha dem
			go func() {
				//send newOrder message for 2 seconds then stop.
				for timeout := time.After(500 * time.Millisecond); ; {
					select {
					case <-timeout:
						return
					default:
					}
					_mtx2.Lock() 
					NewOrderMsgTx <- l
					_mtx2.Unlock() 
					time.Sleep(100 * time.Millisecond)
				}
			}()
*/
		case r := <-ElevStateMsgRx:
			//fmt.Println("new msg from", r.SenderId)

			
			TimerElevatorLost(r,id,lostId,&_mtx3,&_mtx4)

			UpdateElevStateArray(r,&_mtx4)
			
			SyncAllLights(ElevStateArray, id)
			

		case b := <-drv_buttons:
			fmt.Printf("%+v\n", b)

			msg := NewOrderDistributer(ElevStateArray, b.Button, b.Floor, id, Elevator,&_mtx4) //ny mld med hvem som skal ha ordre!
			NewOrderMsgTx <- msg
			/*
			go func() {
				//send newOrder message for 2 seconds then stop.
				for timeout := time.After(1 * time.Second); ; {
					select {
					case <-timeout:
						return
					default:
					}
					_mtx2.Lock() 
					NewOrderMsgTx <- msg
					_mtx2.Unlock() 
					time.Sleep(100 * time.Millisecond)
				}
			}()
			*/

		case m := <-NewOrderMsgRx:
			fmt.Println("new order recieved: ", m)

			if AcceptNewOrder(m, id, Elevator) { //is the new order for this elevator?
				SyncAllLights(ElevStateArray, id)
				OnRequestButtonPress(m.Button.Floor, m.Button.Button, Timer) //sets button request == true on wanted elevator
			}
			
			_mtx.Lock()
			elevstate = ElevStateMsg{
				SenderId: id,
				Message:  "State Update",
				Elevator: Elevator,
			} //oppdatere requests basert på knappetrykk
			_mtx.Unlock()
			
		case a := <-drv_floors:
			fmt.Printf("floor: %+v\n", a)
			//Reset timer for powerloss

			OnFloorArrival(a, id, Timer)
			SyncAllLights(ElevStateArray, id)

			_mtx.Lock()
			
			elevstate = ElevStateMsg{
				SenderId: id,
				Message:  "State Update",
				Elevator: Elevator, //sende over min state hele tiden!
			}
			_mtx.Unlock()
		

		case a := <-drv_obstr:
			fmt.Printf("Obstuction! %+v\n", a)
			if a && Elevator.Behaviour==EBdoorOpen{
				Timer.Reset(3*time.Second)
				Timer.Stop()
			}

			if !a{
				Timer.Stop()
				OnDoorTimeOut()
			}
			
		case <-Timer.C:
			OnDoorTimeOut()

		case id:= <-lostId:
			//alarm
			fmt.Println("In lostId: ", id)
			DistibuteLostOrders(id,&_mtx4,NewOrderMsgTx)

		case a := <-drv_stop:
			fmt.Printf("stop button: %+v\n", a)
			SetStopLamp(true)
			for a{
				for f := 0; f < numFloors; f++ {
					for b := ButtonType(0); b < 3; b++ {
					SetButtonLamp(b, f, false)
					}				
				}
				time.Sleep(500*time.Millisecond)
				for f := 0; f < numFloors; f++ {
					for b := ButtonType(0); b < 3; b++ {
					SetButtonLamp(b, f, true)

					}
				}
				time.Sleep(500*time.Millisecond)
			}
		}
	}

}
