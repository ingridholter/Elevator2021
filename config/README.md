# Configuration package

This package contains definitions of types and global constants used by the other modules. 
Our system can be customized into using different elevator hardware. By changing NumElevators we can set the maximum number of elevators operating. The number of floors is also possible to modify and the system will work as normal. It is also possible to change the time the door stays open when ariving at a floor and the limit for when an error is detected during motor stop. However, the number of buttons should not be changed as this requires more adjustments elsewhere in the code. 
```
const (
	NumElevators     = 3
	NumFloors    int = 4
	NumButtons   int = 3
	må endre her!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
)
```
