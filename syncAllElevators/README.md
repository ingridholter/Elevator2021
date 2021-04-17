# SyncAllElevators

This module ensures that all elevators cooperate and work together. It also makes sure that all orders are taken. If an error occur, this module ensures that another elevator serves the UP and DOWN orders for the faulty elevator. However, this system is only sustainable for one error to happen at a time (motor power loss/network disconnection/software crash). 

# Error detection and handling
As we never want any order to be lost we have choosen to share all orders with all elevators at all time. This make the system sustainable to errors because other elevators can serve the orders of a faulty elevator. <br/>
<br/>
The peers package detects network disconnection or software crash. When an network error occur, the UP and DOWN orders meant for the lost elevator is recalculated and redistributed by the ones remaining on the network. The disconnected elevator will still work in a single-elevator-mode. <br/>
<br/>
The peers package also detects when a software crash occurs and handles this in a similar way to network error. The UP and DOWN orders redistributed between the remaining elevators and the CAB orders are saved and attempted resent to the elevator when it rejoins the system. <br/>
<br/>
When a motor power loss occures, the elevator is still able to communicate with the others over the network. The peers package is not able to detect this error. We have therefore implemented a timer to detect motor power loss, when an elevator is unable to take on of its orders. This timer tells when the other elevators should redistribute and serve the lost elevators orders. The CAB orders of the lost elevator is in this error state never lost and the elevator will serve theese when the motor is repaired. <br/>
<br/>

## Lost elevator variables and arrays
For keeping track of the elevators motor power loss an array of times called elevLastMovedMap are added. The times are reset every time the elevators change their stated. This array helps us detect motor power loss and take action. <br/>
<br/>
A variable called oneLostId keeps track of which id beeing lost. Since this system is only sustainable for one error at a time this id is set to the value of the faulty elevatoris id. When no id's are lost, oneLostId have the value -1. <br/>
<br/>
The lostElevators array keeps track of the systems available elevators. Every elements index and content corresponds to an elevator id. The value in this array is either "lost" or "found" telling the system which elevators to rely on. 


## Weaknesses of our system
Our system reach a vulnerable state when one elevator rejoins after software crash and needs to get its cab orders from the others. It has no way to get theese ordes back if not from the others so this message sending is critical. As we are supposed to be sustainable for package loss, we have no guarantee that this message will arrive. We found a way around this problem using a for loop and send the message serveral times. However, this would not work in a real-life scenario and should be solved in a robust way. <br/>
<br/>
This message need to be sent until we have confirmed its arrival by an answer message. Since there is not another need for this confirmation we stayed with out solution aware of its weakness. 


# Variables, arrays and matrices

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

