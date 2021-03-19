package config


const (
	PeerId = 100 //change for different computers? possible?
	NumElevators = 3
	_numFloors int = 4
	_numButtons int = 3
	AllRequests [_numFloors][_numButtons*NumElevators]bool
	ElevStateArray [NumElevators]elevState
)

type MotorDirection int

const (
	MD_Up   MotorDirection = 1
	MD_Down                = -1
	MD_Stop                = 0
)

type ButtonType int

const (
	BT_HallUp   ButtonType = 0
	BT_HallDown            = 1
	BT_Cab                 = 2
)

type ButtonEvent struct {
	Floor  int
	Button ButtonType
}

type elevBehaviour int

const (
	EBmoving elevBehaviour = 0
	EBstop = 1
	EBdoorOpen = 2
	EBidle = 3
)

type elevState struct {
	floor int
	dirn MotorDirection
	behaviour elevBehaviour
	requests [_numFloors][_numButtons]bool
}