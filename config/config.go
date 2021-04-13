package config

import (
	"time"
)

const (
	NumElevators     = 3
	NumFloors    int = 4
	NumButtons   int = 3
)

var Elevator ElevState //elevator state variable

var ElevStateArray [NumElevators]ElevState

var ElevatorLastMoved = make(map[int]time.Time)

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

type NewOrderMsg struct {
	SenderId   string
	RecieverId string
	Button     ButtonEvent
}

type ElevStateMsg struct {
	SenderId string
	Message  string
	Elevator ElevState
}
