# SyncAllElevators

This module ensures that all elevators cooperate and work together. It also makes sure that all orders are taken. If an error occur, this module ensures that another elevator serves the UP and DOWN orders for the faulty elevator. However, this system is only sustainable for one error to happen at a time (motor power loss/network disconnection/software crash). 

# Error detection and handling
As we never want any order to be lost we have choosen to share all orders with all elevators at all time. This make the system sustainable to errors because other elevators can serve the orders of a faulty elevator. <br/>
<br/>
The peers package detects network disconnection or software crash. When an network error occur, the UP and DOWN orders meant for the lost elevator is sendt over to the ones remaining on the network. Theese are then recalculated and redistributed betweeen the remaining elevators. All the elevators on the network is at every point aware of all the other elevators orders, so that when an error occurs they are able to save the faulty elevators CAB orders and send it back to the faulty elevator when it rejoin the system. <br/>
<br/>
The peers package also detects when a software crash occurs and handles this in a similar way to network error. The UP and DOWN orders redistributed between the remaining elevators and the CAB orders are saved and attempted resent to the elevator when it rejoins the system. <br/>
<br/>
When a motor power loss occures, the elevator is still able to communicate with the others over the network. The peers package is not able to detect this error. We have therefore implemented a timer to detect motor power loss. This timer makes the other ones able to redistribute and serve the lost elevators orders. The CAB orders of the lost elevator is in this error state never lost and the elevator will serve theese when the motor is repaired. <br/>
<br/>



# Variables Arrays and Matrices

This module uses an array of Elevstates for all the available elevators called elevStateArray. 
The Elevstate type has information about the state of the elevator and its requests. 
````
type ElevState struct {
	Floor     int
	Dir       MotorDirection
	Behaviour ElevBehaviour
	Requests  [NumFloors][NumButtons]bool
}
````
The floors are 0-indexed. 
The Requests matrix is boolean and built in this way:

FLOOR\BUTTON| UP | DOWN | CAB
--------------- | ---------- | ---------- | ----------
0 | false| false |  false
1     | false | false|  false
2     | false | false| false
3   | false| false |  false

 Cab orders are id specific and should never be taken by another id. 
