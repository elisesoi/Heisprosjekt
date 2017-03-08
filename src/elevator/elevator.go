package elevator

import (
	. "../driver"
	. "../orders"
	"fmt"
	//"../Network"
	//. "../Network/network/localip"
	"time"
)

//var ip, _ = LocalIP()
//var id = ip[12:15]

func Initialize_elevator(id string) {
	Driver_init()
	fmt.Println("Press STOP button to stop elevator and exit program.\n")
	Driver_set_motor_direction(DIRN_STOP)

	fmt.Println(id)
	State_matrix[id] = Elevator_states{Floors: []int{-1, -1, -1, -1}, Current_direction: DIRN_STOP, Current_floor: Driver_get_floor_sensor_signal(), Alive: 1, Door_open: 0}

	for floor := 0; floor < N_FLOORS; floor++ {
		for button_type := 0; button_type < 2; button_type++ {
			External_orders[floor][button_type] = 0
		}
	}

	Internal_orders[id] = append(Internal_orders[id], 0)
	Internal_orders[id] = append(Internal_orders[id], 0)
	Internal_orders[id] = append(Internal_orders[id], 0)
	Internal_orders[id] = append(Internal_orders[id], 0)

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

	fmt.Println("State matrix: ", State_matrix) //for å se om tallene blir satt rett
	fmt.Println("External: ", External_orders)
	fmt.Println("Internal: ", Internal_orders[id])
}

func Elevator_loop(floor_reached_ch, order_new_state_ch, new_floor_ch, delete_order_ch chan int, new_dir_state_ch chan Driver_motor_dir, new_order_ch chan New_order) {
	state := Elevator_states{}
	go check_floors(floor_reached_ch)
	go check_buttons(new_order_ch)

	for {
		// Stop elevator and exit program if the stop button is pressed
		if Driver_get_stop_signal() != 0 {
			Driver_set_motor_direction(DIRN_STOP)
			fmt.Println("Dette er ny state matrix:", State_matrix)
			fmt.Println("Internal orders: ", Internal_orders)
			break
		}

		Driver_set_floor_indicator(Driver_get_floor_sensor_signal())
		Driver_set_motor_direction(DIRN_UP)
		//floor_reached_ch := make(chan int)
		//order_new_state_ch := make(chan int)

		select {
		case floor := <-floor_reached_ch: // this file
			state.Current_floor = floor
			new_floor_ch <- floor
			floor_reached(floor, delete_order_ch, new_dir_state_ch)

		case new_order := <-new_order_ch:
			//Driver_set_button_lamp(new_order.button, new_order.floor, 1)
			fmt.Println("Ny ordre på ordrekanal")
			fmt.Println("Ny ordre: ", new_order)
			//choose_elevator() + bcast til de andre
			/*
				case new_destination := <-destination_ch: // from Order()
					// ...
				case door := <-door_ch:
			*/
		}
		// Change direction when we reach top/bottom floor
		if Driver_get_floor_sensor_signal() == N_FLOORS-1 {
			Driver_set_motor_direction(DIRN_DOWN)
			//State_matrix[id].Current_direction = DIRN_DOWN
		} else if Driver_get_floor_sensor_signal() == 0 {
			Driver_set_motor_direction(DIRN_UP)
			//State_matrix[id].Current_direction = DIRN_UP
		}
	}
}

func floor_reached(current_floor int, delete_order_ch chan int, new_dir_state_ch chan Driver_motor_dir) {
	if Should_stop(current_floor) == true {
		dir := DIRN_STOP
		Driver_set_motor_direction(DIRN_STOP)
		new_dir_state_ch <- Driver_motor_dir(dir)
		open_door()
		//choose direction
		Delete_orders(current_floor, delete_order_ch)
	}
}

func check_floors(floor_reached_ch chan int) {
	for {
		floor := Driver_get_floor_sensor_signal()
		if floor >= 0 {
			floor_reached_ch <- floor
		}
		//time.Sleep(100 * time.Millisecond)
	}
}

func check_buttons(new_order_ch chan New_order) {
	for {
		fmt.Println("Check buttons kjøres")
		//new_order := New_order{}
		new_order := <-new_order_ch
		fmt.Println(new_order)
		for floor := 0; floor < N_FLOORS; floor++ {
			if Driver_get_button_signal(BUTTON_CALL_UP, floor) == 1 {
				//new_order.floor = floor
				//new_order.button = BUTTON_CALL_UP
				new_order_ch <- new_order
			} else if Driver_get_button_signal(BUTTON_CALL_DOWN, floor) == 1 {
				//new_order.floor = floor
				//new_order.button = BUTTON_CALL_DOWN
				new_order_ch <- new_order
			} else if Driver_get_button_signal(BUTTON_COMMAND, floor) == 1 {
				//new_order.floor = floor
				//new_order.button = BUTTON_COMMAND
				new_order_ch <- new_order
			}
			//fmt.Println("check_buttons funksjonen kjøres")
		}
	}
}

func open_door() {
	Driver_set_door_open_lamp(1)
	//State_matrix[id].Door_open = 1
	fmt.Println("Door open")
	//gi beskjed til de andre

	time.Sleep(3 * time.Second)
	Driver_set_door_open_lamp(0)
	fmt.Println("Door closed")
	//State_matrix[id].Door_open = 0
	//fmt.Println(State_matrix)
	// gi beskjed til de andre
}
