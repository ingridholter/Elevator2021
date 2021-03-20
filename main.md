package main

import(
	"fmt"
	. "./elevio"
	. "./config"
)

func main() {
	//timerDoor()
	//fmt.Println("time out")

	

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
	OnInitBetweenFloors()

	for {
		select {
		case a := <-drv_buttons:
			fmt.Printf("%+v\n", a)
			OnRequestButtonPress(a.Floor,a.Button)
			//SetButtonLamp(a.Button, a.Floor, true)
			

		case a := <-drv_floors:
			fmt.Printf("%+v\n", a)
			if a == numFloors-1 {
				d = MD_Down
			} else if a == 0 {
				d = MD_Up
			}
			OnFloorArrival(a)
			OnFloorTimeOut()
			//SetMotorDirection(d)

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
	
}
