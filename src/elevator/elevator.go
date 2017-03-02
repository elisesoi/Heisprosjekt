package elevator

import (
	. "../driver"
	. "../orders"
	"fmt"
	//"../Network"
	. "../Network/network/localip"
	"time"
)

func Initialize_elevator() {
	Driver_init()
	fmt.Println("Press STOP button to stop elevator and exit program.\n")
	Driver_set_motor_direction(DIRN_STOP)

	//dette er nok noe network bør ta seg av, og sende på kanal til denne modulen
	ip, _ := LocalIP()
	id := ip[12:15]
	//skal denne være global, utenfor funksjonen??
	State_matrix := make(map[string]Elevator_states)
	State_matrix[id] = Elevator_states{0, 0, 0, 0, DIRN_STOP, Driver_get_floor_sensor_signal(), 1, 0}

	for {
		if Driver_get_floor_sensor_signal() == 0 {
			Driver_set_motor_direction(DIRN_STOP)
			open_door()
			break
		} else {
			Driver_set_motor_direction(DIRN_DOWN)
		}
	}

	//JUHUU! eg leve. Her e id-en min :) peerupdatech?
	//spør de andre etter oppdatering. if svar fra de andre: oppdater state_matrix

	//fmt.Println(State_matrix) //for å se om tallene blir satt rett
}

func Elevator_loop() {
	for {
		Driver_set_floor_indicator(Driver_get_floor_sensor_signal())

		// Change direction when we reach top/bottom floor
		if Driver_get_floor_sensor_signal() == N_FLOORS-1 {
			Driver_set_motor_direction(DIRN_DOWN)
		} else if Driver_get_floor_sensor_signal() == 0 {
			Driver_set_motor_direction(DIRN_UP)
		}

		// Stop elevator and exit program if the stop button is pressed
		if Driver_get_stop_signal() != 0 {
			Driver_set_motor_direction(DIRN_STOP)
			break
		}

		if Driver_get_floor_sensor_signal() == 1 {
			Should_stop()
			//send oppdatering til statematrix og til de andre at heis har ny current_state
		}
	}
}

func open_door() {
	Driver_set_door_open_lamp(1)
	//State_matrix[id] = Elevator_states{_, _, _, _, _, _, _, 1}
	//fmt.Println(State_matrix)
	//gi beskjed til de andre

	time.Sleep(3 * time.Second)
	Driver_set_door_open_lamp(0)
	//State_matrix[id][8] = 0
	//fmt.Println(State_matrix)
	// gi beskjed til de andre
}
