package config

const (
	PeerId           = 100 //change for different computers? possible?
	NumElevators     = 3
	NumFloors    int = 4
	NumButtons   int = 3
)

var ActiveElevators = 3
var AllRequests [NumFloors][NumButtons * NumElevators]bool
var ElevStateArray [NumElevators]ElevState

type ButtonType int

const (
	BT_HallUp   ButtonType = 0
	BT_HallDown            = 1
	BT_Cab                 = 2
)

type MotorDirection int

const (
	MD_Up   MotorDirection = 1
	MD_Down                = -1
	MD_Stop                = 0
)

type ButtonEvent struct {
	Floor  int
	Button ButtonType
}

type ElevBehaviour int

const (
	EBmoving   ElevBehaviour = 0
	EBstop                   = 1
	EBdoorOpen               = 2
	EBidle                   = 3
)

type ElevState struct {
	Floor     int
	Dir       MotorDirection
	Behaviour ElevBehaviour
	Requests  [NumFloors][NumButtons]bool
}
