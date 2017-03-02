//Defining different structs we need in several modules

package driver

const N_FLOORS int = 4
const N_BUTTONS int = 3
const EXTERNAL_ORDER int = 9
const N_ELEVATORS int = 3


const MOTOR_SPEED = 2800


type Driver_button_type int
const (
	BUTTON_CALL_UP = 0
    BUTTON_CALL_DOWN = 1
    BUTTON_COMMAND = 2
)


type Driver_motor_dir int
const (
	DIRN_DOWN = -1
	DIRN_STOP = 0
	DIRN_UP = 1
)


type Elevator_states struct {
	Floors []int
	Current_direction int
	Current_floor int
	Alive int
	Door_open int
}

type External_order struct {
	Up int
	Down int
}

/*
type MsgType int
const (
    MSG_NEW_ORDER = iota
    MSG_ORDER_DONE
    MSG_ORDER_ACK
)

type NetworkMessage struct {
    type_of_message  MsgType
    state            Elevator_states
}
*/
