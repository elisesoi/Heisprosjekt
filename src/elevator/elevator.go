package elevator

import (
	//"../orders"
	."../driver"
	//"../Network"
	."../Network/network/localip"
	"time"
)



func Initialize_elevator(){
	Driver_init()
	Driver_set_motor_direction(DIRN_DOWN)
	if Driver_get_floor_sensor_signal() == 1{
		Driver_set_motor_direction(DIRN_STOP)
		open_door()
	}
	//sett alle verdier i state_matrix til 0
	var State_matrix map[string]Elevator_states
	//spør de andre etter oppdatering. if svar fra de andre: oppdater state_matrix
	//else: 
	ip_adress, error := LocalIP()
	var id string = ip_adress[14:16]
	//State_matrix["id"].Floors = [0]
	State_matrix["id"].Current_direction = DIRN_DOWN
	State_matrix["id"].Current_floor = Driver_get_floor_sensor_signal()
	State_matrix["id"].Alive = 1
	State_matrix["id"].Door_open = 0
	fmt.Println(State_matrix) //for å se om tallene blir satt rett
}

func Elevator_loop(){
	for {
		Driver_set_floor_indicator(Driver_get_floor_sensor_signal())

        // Change direction when we reach top/bottom floor
        if Driver_get_floor_sensor_signal() == N_FLOORS - 1 {
            Driver_set_motor_direction(DIRN_DOWN)
        } else if Driver_get_floor_sensor_signal() == 0 {
            Driver_set_motor_direction(DIRN_UP)
        }

        // Stop elevator and exit program if the stop button is pressed
        if (Driver_get_stop_signal() != 0) {
            Driver_set_motor_direction(DIRN_STOP)
            //fmt.Println(internal_order)
            break
        }

		if Driver_get_floor_sensor_signal() == 1{
			Should_stop()
			//send oppdatering til statematrix og til de andre at heis har ny current_state
		}
	}
}

func open_door(){
	Driver_set_door_open_lamp(1)
	//set State_matrix["id"].Door_open = 1
	//gi beskjed til de andre
	time.Sleep(3 * time.Second)
	Driver_set_door_open_lamp(0)
	// set State_matrix["id"].Door_open = 0
	// gi beskjed til de andre
}