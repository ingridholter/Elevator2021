# runElevator.go
Most of the code used in this module is rewritten from C code given by the [project resources](https://github.com/TTK4145/Project-resources)


This is the package for the behaviour of one elevator. It contains code for handling order given by the syncElevators module and executes them. 

ElevState is the type for the elvator. It conatins the floor, direction, behaviour (door open, idle, stop, moving) and requests for one elevator. 
````
type ElevState struct {
	Floor     int
	Dir       MotorDirection
	Behaviour ElevBehaviour
	Requests  [NumFloors][NumButtons]bool
}
````

# SendStateTicker.C

Using golangs library "time", a ticker is used to send information over udp about my elevState. This happens every 500 milisecond.

`````
SendtStateTicker:=time.NewTicker(500*time.Millisecond)

for{
  select{
  case <-SendStateTicker.C:
    ElevSateMsgTx<-elevator
`````

# NewOrderCh

The new Order channel is used for communication with the syncElevators module. Once an order is put on this channel, it is executed by the fsm. 


# RequestShouldStop

The elevator does not stop if it is heading in the opposite direction of the button pushed on a floor. This is implemended in the RequestsShouldStop function. 
