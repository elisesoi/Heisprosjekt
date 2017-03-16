package orders

import (
	. "../driver"
	"fmt"
	//. "math"
)

func Order(order_new_state_ch chan int, new_dir_state_ch chan Driver_motor_dir, new_order_ch, delete_order_ch chan Order_type, new_peer_ch chan string, lost_peer_ch chan []string, new_state_ch chan Elevator_states, id string) {
	for {
		select {
		case floor := <-order_new_state_ch:
			state := State_matrix[id]
			state.Current_floor = floor
			State_matrix[id] = state
			new_state_ch <- State_matrix[id]
			//fmt.Println("State_matrix i hver etg: ", State_matrix)

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
				state.Floors[new_order.Floor] = BUTTON_COMMAND
				State_matrix[id] = state
				Internal_orders[id][new_order.Floor] = 1
				Driver_set_button_lamp(new_order.Button, new_order.Floor, 1)

			} else if new_order.Button == BUTTON_CALL_UP {
				if State_matrix[id].Floors[new_order.Floor] != BUTTON_COMMAND { // || State_matrix[id].Floors[new_order.Floor] != BUTTON_COMMAND
					who_takes_order := choose_elevator(id, new_order.Floor) //spør kost hvem som skal ta bestilling
					fmt.Println("Har valgt heis: ", who_takes_order, " til å ta ordre ", new_order, " bestilt i ", id)
					state := State_matrix[who_takes_order]
					state.Floors[new_order.Floor] = BUTTON_CALL_UP
					State_matrix[who_takes_order] = state
				}
				Driver_set_button_lamp(new_order.Button, new_order.Floor, 1)

			} else if new_order.Button == BUTTON_CALL_DOWN {
				if State_matrix[id].Floors[new_order.Floor] != BUTTON_COMMAND {
					who_takes_order := choose_elevator(id, new_order.Floor) //spør kost hvem som skal ta bestilling
					fmt.Println("Har valgt heis: ", who_takes_order, " til å ta ordre ", new_order, " bestilt i ", id)
					state := State_matrix[who_takes_order]
					state.Floors[new_order.Floor] = BUTTON_CALL_DOWN
					State_matrix[who_takes_order] = state
					fmt.Println("State matrix etter kost beregnet: ", State_matrix)
				}
				Driver_set_button_lamp(new_order.Button, new_order.Floor, 1)
			}

		case delete_order := <-delete_order_ch:
			state := State_matrix[id]
			state.Floors[delete_order.Floor] = NO_ORDERS
			State_matrix[id] = state
			//bcast til de andre at ordre er slettet
			Internal_orders[id][delete_order.Floor] = 0
			External_orders[delete_order.Floor][0] = 0
			External_orders[delete_order.Floor][1] = 0

		case newPeer := <-new_peer_ch:
			//legg til newPeer til state_matrix
			State_matrix[newPeer] = Elevator_states{Id: newPeer, Floors: []Driver_button_type{NO_ORDERS, NO_ORDERS, NO_ORDERS, NO_ORDERS}, Current_direction: DIRN_STOP, Prev_direction: DIRN_STOP, Current_floor: 0, Alive: 1}
			if newPeer != id {
				if Internal_orders[newPeer] != nil {
					fmt.Println("indre ordre fra før")
					//få oppdatering fra de heis
				} else {
					Internal_orders[newPeer] = append(Internal_orders[newPeer], 0)
					Internal_orders[newPeer] = append(Internal_orders[newPeer], 0)
					Internal_orders[newPeer] = append(Internal_orders[newPeer], 0)
					Internal_orders[newPeer] = append(Internal_orders[newPeer], 0)
				}
			}

			//fmt.Println(State_matrix)
			//spør den nye heisen om den har bestilliger fra før? etg? oppdater?
			//legg til newPeer til Internal_orders
			//if Internal_orders[newPeer] ikke finnes{}

			// hvis internal orders allerede har en heis med den id-en, så send internal orders til denne heisen.

		case lostPeer := <-lost_peer_ch:
			if lostPeer[0] != id {
				//gi bestillinger fra lostPeer til en annen heis
				for floor := 0; floor < N_FLOORS; floor++ {
					if State_matrix[lostPeer[0]].Floors[floor] != NO_ORDERS {
						State_matrix[id].Floors[floor] = BUTTON_COMMAND
						if State_matrix[lostPeer[0]].Floors[floor] == BUTTON_COMMAND {
							Internal_orders[id][floor] = 1
						}
						fmt.Println("Heisen som døde ", lostPeer[0], " hadde en bestilling i etg ", floor, "som ble lagt til i heis: ", id)
					}
				}
				delete(State_matrix, lostPeer[0]) //slett lostPeer fra State_matrix

			}

		case new_state := <-new_state_ch:
			//fmt.Println("Den nye staten sendt på kanal til alle: ", new_state)
			//fmt.Println("***", new_state, "***")
			if new_state.Id != id {
				State_matrix[new_state.Id] = Elevator_states{Id: new_state.Id, Floors: new_state.Floors, Current_direction: new_state.Current_direction, Prev_direction: new_state.Prev_direction, Current_floor: new_state.Current_floor, Alive: 1}
				fmt.Println("Dette er oppdatert map fra new state: ", State_matrix)
			}

			no_order_var := 0

			for i := 0; i < N_FLOORS; i++ {
				no_order_var = 0
				if new_state.Floors[i] == BUTTON_CALL_UP {
					Driver_set_button_lamp(BUTTON_CALL_UP, i, 1)
				} else if new_state.Floors[i] == BUTTON_CALL_DOWN {
					Driver_set_button_lamp(BUTTON_CALL_DOWN, i, 1)
				}
				for key := range State_matrix {
					if State_matrix[key].Floors[i] != NO_ORDERS { // == BUTTON_CALL_DOWN || State_matrix[key].Floors[i] == BUTTON_CALL_UP {
						no_order_var++
					}
				}
				if no_order_var == 0 {
					Driver_set_button_lamp(BUTTON_CALL_UP, i, 0)
					Driver_set_button_lamp(BUTTON_CALL_DOWN, i, 0)
				}
			}
		}
	}
}

func calculate_cost(state Elevator_states, ordered_floor int) int { //må ta inn argument; statematrix og
	cost := 0
	for i := 0; i < N_FLOORS; i++ {
		if state.Floors[i] != NO_ORDERS {
			cost += 2
		}
	}
	cost1 := (state.Current_floor - ordered_floor) * 2
	if cost1 < 0 {
		cost -= cost1
	} else if cost1 >= 0 {
		cost += cost1
	}
	// if på vei i feil retning +2
	return cost
}

func choose_elevator(id string, ordered_floor int) string {
	min_cost := 100
	who_takes_order := id
	for id, v := range State_matrix {
		cost := calculate_cost(v, ordered_floor)
		if cost < min_cost {
			cost = min_cost
			who_takes_order = id
		}
	}
	return who_takes_order
}

func orders_above(current_floor int, id string) int {
	for i := current_floor + 1; i < N_FLOORS; i++ {
		if State_matrix[id].Floors[i] != NO_ORDERS {
			return 1
		}
	}
	return 0
}

func orders_under(current_floor int, id string) int {
	for j := current_floor - 1; j >= 0; j-- {
		if State_matrix[id].Floors[j] != NO_ORDERS {
			return 1
		}
	}
	return 0
}

func Should_stop(current_floor int, id string) bool {
	if State_matrix[id].Floors[current_floor] != NO_ORDERS {
		if State_matrix[id].Floors[current_floor] == BUTTON_COMMAND {
			return true
		} else if State_matrix[id].Floors[current_floor] == BUTTON_CALL_DOWN && State_matrix[id].Current_direction == -1 {
			return true
		} else if State_matrix[id].Floors[current_floor] == BUTTON_CALL_UP && State_matrix[id].Current_direction == 1 {
			return true
		} else if State_matrix[id].Floors[current_floor] == BUTTON_CALL_DOWN && State_matrix[id].Current_direction == 1 && orders_above(current_floor, id) == 0 {
			return true
		} else if State_matrix[id].Floors[current_floor] == BUTTON_CALL_UP && State_matrix[id].Current_direction == -1 && orders_under(current_floor, id) == 0 {
			return true
		}
	}
	return false
}

func Choose_direction(prev_dir, current_direction Driver_motor_dir, current_floor int, id string) Driver_motor_dir {
	switch current_direction {
	case DIRN_STOP:
		if prev_dir == 1 && orders_above(current_floor, id) == 1 {
			return DIRN_UP
		} else if prev_dir == -1 && orders_under(current_floor, id) == 1 {
			return DIRN_DOWN
		} else if prev_dir == 0 && orders_above(current_floor, id) == 1 {
			return DIRN_UP
		} else if prev_dir == 0 && orders_under(current_floor, id) == 1 {
			return DIRN_DOWN
		}
	case DIRN_UP:
		if orders_above(current_floor, id) == 1 {
			return DIRN_UP
		} else if orders_under(current_floor, id) == 1 {
			return DIRN_DOWN
		}
	case DIRN_DOWN:
		if orders_under(current_floor, id) == 1 {
			return DIRN_DOWN
		} else if orders_above(current_floor, id) == 1 {
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
