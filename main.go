package main
import(
	"time"
	"fmt"
	. "./elevio"
	"fmt"
	. "./config"
)
func timerDoor() {

	timer := time.NewTimer(3*time.Second)
	<-timer.C
}



func main() {
	timerDoor()
	fmt.Println("time out")

	for {
		select {
		case a := <-drv_buttons:
			fmt.Printf("%+v\n", a)
			SetButtonLamp(a.Button, a.Floor, true)

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
	
}
