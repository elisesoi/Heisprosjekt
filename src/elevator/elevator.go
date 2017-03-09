package elevator

import (
	. "../driver"
	. "../orders"
	"fmt"
	"../Network"
	//. "../Network/network/localip"
	"time"
)

//var ip, _ = LocalIP()
//var id = ip[12:15]

func Initialize_elevator(id string, ) {
	Driver_init()
	fmt.Println("Press STOP button to stop elevator and exit program.")
	Driver_set_motor_direction(DIRN_STOP)

	State_matrix[id] = Elevator_states{Floors: []int{0,0,0,0}, Current_direction: DIRN_STOP, Prev_direction: DIRN_STOP, Current_floor: 0, Alive: 1}

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

func Elevator_loop(floor_reached_ch, order_new_state_ch chan int, new_dir_state_ch chan Driver_motor_dir, new_order_ch, delete_order_ch chan Order_type) {
	//floor_reached_ch := make(chan int)
	//order_new_state_ch := make(chan int)

	stop_button_pushed_ch := make(chan int)
	go check_floors(floor_reached_ch)
	go check_stop_button(stop_button_pushed_ch)
	go check_buttons(new_order_ch)

	state := Elevator_states{}

	for {
		Driver_set_floor_indicator(Driver_get_floor_sensor_signal())

		select {
		case stop := <-stop_button_pushed_ch:
			stop++
			Driver_set_motor_direction(DIRN_STOP)
			break

		case floor := <-floor_reached_ch: // this file
			state.Current_floor = floor
			order_new_state_ch <- floor
			//fmt.Println(State_matrix)
			floor_reached(new_dir_state_ch, delete_order_ch, floor)
			if floor >= N_FLOORS -1 {
				new_dir_state_ch <- DIRN_DOWN
				Driver_set_motor_direction(DIRN_DOWN) // egentlig Choose_direction()
			} else if floor <= 0 {
				new_dir_state_ch <- DIRN_UP
				Driver_set_motor_direction(DIRN_UP) // egentlig Choose_direction()
			}

		case new_order := <-new_order_ch:
			//fmt.Println("New order!!", new_order)
			//id := Network.GetLocalId()
			Driver_set_button_lamp(new_order.Button, new_order.Floor, 1)
			//fmt.Println("Choose direction: ", Choose_direction(State_matrix[id].Prev_direction, new_dir_state_ch, State_matrix[id].Current_floor, id))
			//Driver_set_motor_direction(Choose_direction(State_matrix[id].Prev_direction, State_matrix[id].Current_floor, id))
			//fmt.Println(State_matrix[id].Prev_direction)
			//fmt.Println(State_matrix[id].Current_floor)
			//Choose_direction(prev_dir Driver_motor_dir, current_floor int, id string) Driver_motor_dir
			//sender bare beskjed til de andre om at det er kommet bestilling

		}
	}
}

func floor_reached(new_dir_state_ch chan Driver_motor_dir, delete_order_ch chan Order_type, current_floor int) {
	if Should_stop(current_floor) == true {
		id := Network.GetLocalId()
		prev_dir := State_matrix[id].Prev_direction
		//fmt.Println("Previous direction ", prev_dir)
		Driver_set_motor_direction(DIRN_STOP)
		new_dir_state_ch <- DIRN_STOP
		go open_door()
		if open_door() == 0 {
			Delete_orders(delete_order_ch)
			Driver_set_button_lamp(BUTTON_COMMAND, current_floor, 0)
			Driver_set_button_lamp(BUTTON_CALL_DOWN, current_floor, 0)
			Driver_set_button_lamp(BUTTON_CALL_UP, current_floor, 0)
			//Driver_set_motor_direction(DIRN_UP) //foreløpig
			new_dir_state_ch <- Choose_direction(prev_dir, current_floor, id)
			Driver_set_motor_direction(Choose_direction(prev_dir, current_floor, id))
		}
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
			} else if Driver_get_button_signal(BUTTON_CALL_DOWN, floor) == 1 {
				new_order.Floor = floor
				new_order.Button = BUTTON_CALL_DOWN
				new_order_ch <- new_order
			} else if Driver_get_button_signal(BUTTON_CALL_UP, floor) == 1 {
				new_order.Floor = floor
				new_order.Button = BUTTON_CALL_UP
				new_order_ch <- new_order
			} 
		}
	}
}

func open_door() int{
	Driver_set_door_open_lamp(1)
	time.Sleep(3 * time.Second)
	Driver_set_door_open_lamp(0)
	return 0
}
