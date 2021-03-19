package elevatorObserver
//module for keeping track over states in other elevator, handling message error and setting lights to "take order"

import(
	. "./config"
)

type ElevStateMsg struct {
	id int
	PeerId int
	Message string
	AllRequests [numFloors][NumElevators*3]bool //3 is number of button types
	elevator elevState
	Iter int
}

type PowerLoss struct{
	id int
	Message string
	elevator elevState
}

//useful functions

func UpdateElevStateArray(msg ElevStateMsg){
	ElevStateArray[msg.PeerId-100] = msg.elevator //id== 100-101-102
}

func UpdateAllRequestsMatrix(msg ElevStateMsg){
	AllRequests = msg.AllRequests
}

func AllRequestToRequests(){
	for f:=0;f<4;f++{
		for i:=3*PeerId; i<3*PeerId+3; i++{
			elevator.requests[f][i-3*PeerId] = AllRequests[f][i] 
		}
	}
}

//sette lys og ta ordre, det kommer fra motatt R matrise 
func AcceptOrder(elevator elevState){
	setLights(elevator);
}

/*
MANGLER:
- feilhåndering

*/

newElevState:= make(chan ElevStateMsg)
elevPowerLoss:= make(chan PowerLoss)

//newElevState messages
go bcast.Transmitter(16569, newElevState)
go bcast.Receiver(16569, elevPowerLoss) //10.100.23.209



go func() {
	helloMsg := HelloMsg{
			Answer: false,
			Message:"Hello from " + id,
			R:[2][3]bool{
				{true,false,false},
				{false,true,false},
			},
			Iter: 0}
	for {
		helloMsg.Iter++
		helloTx <- helloMsg
		time.Sleep(3 * time.Second)
	}
}()



for {
	select{
	
	case b := <-helloRx:
		fmt.Printf("Received: %#v\n", b)
		
		}
	}
}



