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

	fmt.Println("Key til state_matrix")
	fmt.Println("State matrix: ", State_matrix) //for å se om tallene blir satt rett
	fmt.Println("External: ", External_orders)
	fmt.Println("Internal: ", Internal_orders[id])
}

func Elevator_loop(floor_reached_ch, order_new_state_ch chan int, new_dir_state_ch chan Driver_motor_dir, new_order_ch chan Order_type) {

	//floor_reached_ch := make(chan int)
	//order_new_state_ch := make(chan int)

	stop_button_pushed_ch := make(chan int)
	go check_floors(floor_reached_ch)
	go check_stop_button(stop_button_pushed_ch)
	go check_buttons(new_order_ch)

	state := Elevator_states{}

	for {
		select {
		case stop := <-stop_button_pushed_ch:
			stop++
			Driver_set_motor_direction(DIRN_STOP)
			break

		case floor := <-floor_reached_ch: // this file
			state.Current_floor = floor
			order_new_state_ch <- floor
			fmt.Println(State_matrix)
			floor_reached(new_dir_state_ch)

		case new_order := <-new_order_ch:
			fmt.Println("New order!!", new_order)
			if Should_stop() == true {
				Driver_set_motor_direction(DIRN_STOP)
				open_door()
				//bcast til de andre
				//Delete_orders()
				//Choose_direction()
			}

		}

		/*
			case new_order := <-new_order_ch:
				// ...
			case new_destination := <-destination_ch: // from Order()
				// ...
			case door := <-door_ch:
		*/
		//}

		Driver_set_floor_indicator(Driver_get_floor_sensor_signal())

		// Change direction when we reach top/bottom floor
		if Driver_get_floor_sensor_signal() == N_FLOORS-1 {
			Driver_set_motor_direction(DIRN_DOWN)
			//State_matrix[id].Current_direction = DIRN_DOWN
		} else if Driver_get_floor_sensor_signal() == 0 {
			Driver_set_motor_direction(DIRN_UP)
			//State_matrix[id].Current_direction = DIRN_UP
		}

		// Stop elevator and exit program if the stop button is pressed
		//if Driver_get_stop_signal() != 0 {
		//	Driver_set_motor_direction(DIRN_STOP)
		//	fmt.Println("Dette er ny state matrix:", State_matrix)
		//	break
		//}
		fl := Driver_get_floor_sensor_signal()
		if fl == 0 || fl == 1 || fl == 2 || fl == 3 {
			/*
				if Should_stop() == true {
					Driver_set_motor_direction(DIRN_STOP)
					open_door()
					//send beskjed til de andre om hvilken ordre som er ferdig. Må ikke vente på svar.
					Delete_orders()
					//choose_direction()
					Driver_set_motor_direction(DIRN_UP)
				}
			*/
			//send oppdatering til statematrix og til de andre at heis har ny current_state
			//fmt.Println(State_matrix)
		}
	}

}

func floor_reached(new_dir_state_ch chan Driver_motor_dir) {
	if Should_stop() == true {
		dir := DIRN_STOP
		Driver_set_motor_direction(DIRN_STOP)
		new_dir_state_ch <- Driver_motor_dir(dir)
		open_door()
		//choose direction
	}
}

func check_floors(floor_reached_ch chan int) {
	for {
		floor := Driver_get_floor_sensor_signal()
		if floor >= 0 {
			floor_reached_ch <- floor
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func check_stop_button(stop_button_pushed_ch chan int) {
	for {
		if Driver_get_stop_signal() == 1 {
			stop_button_pushed_ch <- 1
		}
	}
}

func check_buttons(new_order_ch chan Order_type) {
	var new_order Order_type
	for {
		for floor := 0; floor < N_FLOORS; floor++ {
			if Driver_get_button_signal(BUTTON_COMMAND, floor) == 1 {
				new_order.Floor = floor
				new_order.Button = BUTTON_COMMAND
				new_order_ch <- new_order
			} /*else if Driver_get_button_signal(BUTTON_DOWN, floor) == 1 {
				new_order.Floor = floor
				new_order.Button = BUTTON_DOWN
				new_order_ch <- new_order
			} else if Driver_get_button_signal(BUTTON_UP, floor) == 1 {
				new_order.Floor = floor
				new_order.Button = BUTTON_UP
				new_order_ch <- new_order
			} */
		}
	}
}

func open_door() {
	Driver_set_door_open_lamp(1)
	//State_matrix[id].Door_open = 1
	fmt.Println(State_matrix)
	//gi beskjed til de andre

	time.Sleep(3 * time.Second)
	Driver_set_door_open_lamp(0)
	//State_matrix[id].Door_open = 0
	//fmt.Println(State_matrix)
	// gi beskjed til de andre
}
