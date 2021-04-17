# SyncAllElevators



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
