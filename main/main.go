package main

import (
	"fmt"
	//. "./elevio"
	. "main/config"
	//"./Network/bcast"
	//"./Network/localip"
	//"time"
	"net"
	"strings"

	. "main/costFunc"
)

//var elevator elevState
/*func Answerfunc() {
	answerMsg:= AnswerMsg{
		Answer: true,
		Message:"Answer from " + id,
		R:[2][3]bool{
				{true,false,false},
				{false,true,false},
		},
		Iter: 0}
}
*/

type HelloMsg struct {
	Answer  bool
	Message string
	R       [2][3]bool
	Iter    int
}

type AnswerMsg struct {
	Answer  bool
	Message string
	R       [2][3]bool
	Iter    int
}

var localIP string

func LocalIP() (string, error) {
	if localIP == "" {
		conn, err := net.DialTCP("tcp4", nil, &net.TCPAddr{IP: []byte{8, 8, 8, 8}, Port: 53})
		if err != nil {
			return "", err
		}
		defer conn.Close()
		localIP = strings.Split(conn.LocalAddr().String(), ":")[0]
	}
	return localIP, nil
}

var e1 = ElevState{
	Floor:     1,
	Dir:       MD_Up,
	Behaviour: EBmoving,
	Requests: [4][3]bool{
		{true, false, false},
		{false, true, false},
		{true, false, false},
		{false, false, false},
	},
}

var e2 = ElevState{
	Floor:     2,
	Dir:       MD_Stop,
	Behaviour: EBdoorOpen,
	Requests: [4][3]bool{
		{true, false, false},
		{false, false, false},
		{true, false, false},
		{false, false, false},
	},
}

var e3 = ElevState{
	Floor:     3,
	Dir:       MD_Stop,
	Behaviour: EBidle,
	Requests: [4][3]bool{
		{false, false, false},
		{false, false, false},
		{false, false, false},
		{false, false, false},
	},
}

var eOld [3]ElevState

func main() {
	eOld[0] = e1
	eOld[1] = e2
	eOld[2] = e3

	//timerDoor()
	fmt.Println("hello")
	/*/var id ="101"
	var id string
	if id == "" {
		localIP, err := LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}
	*/
	fmt.Println(1)
	R := RequestMatrix(eOld, 2, 2)

	for f := 0; f < 4; f++ {
		for i := 0; i < 9; i++ {
			fmt.Println("Request matrix", "floor: ", f, "buttontype: ", ButtonType(i), "value: ", R[f][i])
		}
	}
	/*
		numFloors := 4
		Init("localhost:15657", numFloors)
		var d MotorDirection = MD_Up
		//elevio.SetMotorDirection(d)
		drv_buttons := make(chan ButtonEvent)
		drv_floors := make(chan int)
		drv_obstr := make(chan bool)
		drv_stop := make(chan bool)
		go PollButtons(drv_buttons)
		go PollFloorSensor(drv_floors)
		go PollObstructionSwitch(drv_obstr)
		go PollStopButton(drv_stop)
		//drive down if between floors
		//onInitBetweenFloors()
	*/
	/*
		// We make channels for sending and receiving our custom data types
		helloTx := make(chan HelloMsg)
		helloRx := make(chan HelloMsg)
		// ... and start the transmitter/receiver pair on some port
		// These functions can take any number of channels! It is also possible to
		//  start multiple transmitters/receivers on the same port.
		go bcast.Transmitter(16569, helloTx)
		go bcast.Receiver(16569, helloRx) //10.100.23.209

		//go Answerfunc()

		// The example message. We just send one of these every second.
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
					//svare at vi har motatt mld
					if !b.Answer{
					answerMsg:= HelloMsg{
						Answer: true,
						Message:"Answer from " + id,
						R:[2][3]bool{
								{true,false,false},
								{false,true,false},
						},
						Iter: b.Iter+1}
					helloTx <- answerMsg
					time.Sleep(1 * time.Millisecond)
					}
				}
			}
	*/

	/*
	   case a := <-helloTx:
	   				fmt.Printf("Transmitted: %#v\n", a)
	   	for {
	   		select {
	   		case a := <-drv_buttons:
	   			fmt.Printf("%+v\n", a)
	   			SetButtonLamp(a.Button, a.Floor, true)
	   			SetLights(elevState)

	   		case a := <-drv_floors:
	   			fmt.Printf("%+v\n", a)
	   			if a == numFloors-1 {
	   				d = MD_Down
	   			} else if a == 0 {
	   				d = MD_Up
	   			}
	   			SetMotorDirection(d)

	   		case a := <-drv_obstr:
	   			fmt.Printf("%+v\n", a)
	   			if a {
	   				SetMotorDirection(MD_Stop)
	   			} else {
	   				SetMotorDirection(d)
	   			}

	   		case a := <-drv_stop:
	   			fmt.Printf("%+v\n", a)
	   			for f := 0; f < numFloors; f++ {
	   				for b := ButtonType(0); b < 3; b++ {
	   					SetButtonLamp(b, f, false)
	   				}
	   			}
	   		}
	   	}
	*/
}
