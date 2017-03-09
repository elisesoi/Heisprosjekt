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

func Order(order_new_state_ch chan int, new_dir_state_ch chan Driver_motor_dir, new_order_ch chan Order_type, id string) {
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
				state.Floors[new_order.Floor] = 2
				State_matrix[id] = state
			}
		case delete_order := <- delete_order_ch:
				state := State_matrix[id]
				state.Floors[delete_order.Floor] = 0
				State_matrix[id] = state
				Internal_orders[id][delete_order.Floor] = 0
				External_order[delete_order.Floors][0] = 0
				External_order[delete_order.Floors][1] = 0
				//må bcaste at den har slettet ordre.

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

func Should_stop() bool {
	id := GetLocalId()
	fmt.Println(id)
	current_floor := State_matrix[id].Current_floor
	if current_floor != -1 {
		if State_matrix[id].Floors[current_floor] >= 0 {
			//sjekk tallet i matrisen opp mot dir
			return true
		}
	}
	return false
}

func choose_direction() {
	//Lik som i 1.klasse?
}

delete_order_ch := make(chan Order_type)
func Delete_orders() {
	var order_to_delete Order_type
	order_to_delete.Floor = State_matrix[id].Current_floor
	delete_order_ch <- order_to_delete

}













