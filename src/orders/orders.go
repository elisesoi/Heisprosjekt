package orders

import (
	"fmt"
	//"timer"
	. "../Network"
	. "../driver"
)

func Order(order_new_state_ch chan int, new_dir_state_ch chan Driver_motor_dir, new_order_ch, delete_order_ch chan Order_type, new_peer_ch chan string, lost_peer_ch chan []string, new_state_ch chan Elevator_states, id string) {
	for {
		select {
		case floor := <-order_new_state_ch:
			state := State_matrix[id]
			state.Current_floor = floor
			State_matrix[id] = state
			new_state_ch <- State_matrix[id]

		case dir := <-new_dir_state_ch:
			state := State_matrix[id]
			state.Prev_direction = State_matrix[id].Current_direction
			state.Current_direction = dir
			State_matrix[id] = state
			//fmt.Println("dir: ", dir)

		case new_order := <-new_order_ch:
			//sjekk om det er greit for de andre
			if new_order.Button == BUTTON_COMMAND {
				state := State_matrix[id]
				state.Floors[new_order.Floor] = 1
				State_matrix[id] = state

				Internal_orders[id][new_order.Floor] = 1
				Driver_set_button_lamp(new_order.Button, new_order.Floor, 1)
				//bcast til de andre
			} else if new_order.Button == BUTTON_CALL_UP {
				External_orders[new_order.Floor][1] = EXTERNAL_ORDER
				//fmt.Println("call up", External_orders)
				//spør kost hvem som skal ta bestilling
				//send til de andre, vent på svar
				// MIDLERTIDIG: (bare for at den skal stoppe for ytre bestillinger også.)
				state := State_matrix[id]
				state.Floors[new_order.Floor] = 1
				State_matrix[id] = state
				External_orders[new_order.Floor][1] = 1
				Driver_set_button_lamp(new_order.Button, new_order.Floor, 1)
				//når svar fra alle: legg til i state_matrix og endre External_order til heis som tar bestilling til 1 (istede for 9)
			} else if new_order.Button == BUTTON_CALL_DOWN {
				External_orders[new_order.Floor][0] = EXTERNAL_ORDER
				//fmt.Println("call down", External_orders)
				//spør kost hvem som skal ta bestilling
				//send til de andre, vent på svar
				// MIDLERTIDIG: (bare for at den skal stoppe for ytre bestillinger også.)
				state := State_matrix[id]
				state.Floors[new_order.Floor] = 1
				State_matrix[id] = state
				External_orders[new_order.Floor][0] = 1
				Driver_set_button_lamp(new_order.Button, new_order.Floor, 1)
				//når svar fra alle: legg til i state_matrix og endre External_order til heis som tar bestilling til 1 (istede for 9)
			}

		case delete_order := <-delete_order_ch:
			state := State_matrix[id]
			state.Floors[delete_order.Floor] = 0
			State_matrix[id] = state
			//bcast til de andre at ordre er slettet
			Internal_orders[id][delete_order.Floor] = 0
			External_orders[delete_order.Floor][0] = 0
			External_orders[delete_order.Floor][1] = 0

		case newPeer := <-new_peer_ch:
			//legg til newPeer til state_matrix
			State_matrix[newPeer] = Elevator_states{Id: newPeer, Floors: []int{0, 0, 0, 0}, Current_direction: DIRN_STOP, Prev_direction: DIRN_STOP, Current_floor: 0, Alive: 1}
			if Internal_orders[newPeer] != nil {
				fmt.Println("ingen indre ordre fra før")
				//få oppdatering fra de heis
			} else {
				Internal_orders[newPeer] = append(Internal_orders[newPeer], 0)
				Internal_orders[newPeer] = append(Internal_orders[newPeer], 0)
				Internal_orders[newPeer] = append(Internal_orders[newPeer], 0)
				Internal_orders[newPeer] = append(Internal_orders[newPeer], 0)
			}

			//fmt.Println(State_matrix)
			//spør den nye heisen om den har bestilliger fra før? etg? oppdater?
			//legg til newPeer til Internal_orders
			//if Internal_orders[newPeer] ikke finnes{}

			// hvis internal orders allerede har en heis med den id-en, så send internal orders til denne heisen.
			/*
				case lostPeer := <- lost_peer_ch:
					//gi bestillinger fra lostPeer til en annen heis

					for floor := 0; floor<N_FLOORS; floor++{
						if State_matrix[lostPeer[0]].Floors[floor] == 1{
							local_id := GetLocalId()
							State_matrix[local_id].Floors[floor] = 1
							Internal_orders[local_id][floor] = 1
							fmt.Println("Heisen som døde ", lostPeer[0]," hadde en bestilling i etg ", floor, "som ble lagt til i heis: ", local_id)
						}
					}

					delete(State_matrix, lostPeer) //slett lostPeer fra State_matrix
			*/
		case new_state := <-new_state_ch:
			//fmt.Println("Den nye staten sendt på kanal til alle: ", new_state)
			fmt.Println("Dette er oppdatert map: ", State_matrix)
			fmt.Println("***", new_state, "***")

		}
	}
}

/*

func calculate_cost(state elevator_states) int{ //må ta inn argument; statematrix og
	cost := 0
	if state.Alive == 0{
		cost = 1000
		return cost
	}
	for i:=0; i<N_FLOORS; i++{
		if state.Floors[i] == 1{
			cost += 2
		}
	}
	// if på vei i feil retning +2
	return cost
}

func choose_elevator(){
	min_cost := 100
	for id, v := range state_matrix{
		cost := calculate_cost(v)
		if cost < min_cost {
			cost = min_cost
			who_takes_order = id
		}
	// bruk nettverk og send til de andre hvem som tar bestilling
	}
	//if //svar_fra_alle_heiser_over_acc_channel{
		//sett tall i state_matrix
		//sett tall i orders
		//tenn lys / kall lysfunksjon
	//}
}
*/

func orders_above(current_floor int) int {
	id := GetLocalId()
	for i := current_floor; i < N_FLOORS; i++ {
		if State_matrix[id].Floors[i] == 1 {
			return 1
		}
	}
	return 0
}

func orders_under(current_floor int) int {
	id := GetLocalId()
	for j := current_floor; j >= 0; j-- {
		if State_matrix[id].Floors[j] == 1 {
			return 1
		}
	}
	return 0
}

func Should_stop(current_floor int) bool {
	id := GetLocalId()
	if State_matrix[id].Floors[current_floor] == 1 {
		if Internal_orders[id][current_floor] == 1 {
			return true
		} else if External_orders[current_floor][0] == 1 && State_matrix[id].Current_direction == -1 {
			return true
		} else if External_orders[current_floor][1] == 1 && State_matrix[id].Current_direction == 1 {
			return true
		}
	}
	return false
}

func Choose_direction(prev_dir, current_direction Driver_motor_dir, current_floor int, id string) Driver_motor_dir {
	switch current_direction {
	case DIRN_STOP:
		//fmt.Println("Valgte å stå i ro. Min forrige retning var: ", prev_dir)
		if prev_dir == 1 && orders_above(current_floor) == 1 {
			return DIRN_UP
		} else if prev_dir == -1 && orders_under(current_floor) == 1 {
			return DIRN_DOWN
		} else if prev_dir == 0 && orders_above(current_floor) == 1 {
			return DIRN_UP
		} else if prev_dir == 0 && orders_under(current_floor) == 1 {
			return DIRN_DOWN
		}
	case DIRN_UP:
		if orders_above(current_floor) == 1 {
			return DIRN_UP
		} else if orders_under(current_floor) == 1 {
			return DIRN_DOWN
		}
	case DIRN_DOWN:
		if orders_under(current_floor) == 1 {
			return DIRN_DOWN
		} else if orders_above(current_floor) == 1 {
			return DIRN_UP
		}
	}
	return DIRN_STOP
}

func Delete_orders(delete_order_ch chan Order_type, floor int) {
	var order_to_delete Order_type
	order_to_delete.Floor = floor
	delete_order_ch <- order_to_delete
}
