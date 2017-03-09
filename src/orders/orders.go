package orders

import (
	"fmt"
	//"timer"
	. "../Network"
	. "../driver"
	//. "../Network/network/localip"
)

/*
func Order_default(){
	for i:=0; i<N_FLOORS; i++{
		if Driver_get_button_signal(BUTTON_COMMAND, i) == 1{
			Driver_set_button_lamp(BUTTON_COMMAND, i, 1)
			Internal_orders[id][i] = 1
			//State_matrix[id].Floors[i] = 1 //går ikke fordi Floors[] er tom...
		}
	}
}
*/

func Order(order_new_state_ch chan int, new_dir_state_ch chan Driver_motor_dir, new_order_ch, delete_order_ch chan Order_type, id string) {
	for {
		select {
		case floor := <-order_new_state_ch:
			state := State_matrix[id]
			state.Current_floor = floor
			State_matrix[id] = state
			//Bør si i fra til de andre hvilken etg han er i
			//
		case dir := <-new_dir_state_ch:
			state := State_matrix[id]
			state.Current_direction = dir
			State_matrix[id] = state

		case new_order := <-new_order_ch:
			//sjekk om det er greit for de andre
			if new_order.Button == BUTTON_COMMAND {
				state := State_matrix[id]
				state.Floors[new_order.Floor] = 1
				State_matrix[id] = state

				Internal_orders[id][new_order.Floor] = 1
				//bcast til de andre
			}else if new_order.Button == BUTTON_CALL_UP{
				External_orders[new_order.Floor][1] = EXTERNAL_ORDER
				fmt.Println("call up", External_orders)
				//spør kost hvem som skal ta bestilling
				//send til de andre, vent på svar
				//når svar fra alle: legg til i state_matrix og endre External_order til heis som tar bestilling til 1 (istede for 9)
			}else if new_order.Button == BUTTON_CALL_DOWN{
				External_orders[new_order.Floor][0] = EXTERNAL_ORDER
				fmt.Println("call down", External_orders)
				//spør kost hvem som skal ta bestilling
				//send til de andre, vent på svar
				//når svar fra alle: legg til i state_matrix og endre External_order til heis som tar bestilling til 1 (istede for 9)
			}
		case delete_order := <-delete_order_ch:
			state := State_matrix[id]
			state.Floors[delete_order.Floor] = 0
			State_matrix[id] = state
			//sette Internal_orders og External_orders til 0
			//bcast til de andre at ordre er slettet
			//Internal_orders[id][delete_order.Floor] = 0
			//External_order[delete_order.Floors][0] = 0
			//External_order[delete_order.Floors][1] = 0
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

func Should_stop(current_floor int) bool {
	id := GetLocalId()
	fmt.Println(id)
	fmt.Println("Current floor", current_floor)
	fmt.Println("ordre i etg fra matrise: ",State_matrix[id].Floors[current_floor])
	if State_matrix[id].Floors[current_floor] == 1 {
		if Internal_orders[id][current_floor] == 1{
			return true
		}
		//må sjekke mot Internal_orders og External_orders
		//sjekk tallet i matrisen opp mot dir
		return true
	}
	fmt.Println("Should stop?")
	return false
}

func choose_direction() {
	//Lik som i 1.klasse?
}


func Delete_orders(delete_order_ch chan Order_type) {
	id := GetLocalId()
	var order_to_delete Order_type
	order_to_delete.Floor = State_matrix[id].Current_floor
	delete_order_ch <- order_to_delete
}













