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

		}
	}
}

/*
func Order(sender_ch, recv_ch chan string) { //skal denne sende over Elevator_states? og ikke string?
    timer := time.NewTicker(time.Second).C
    var state_matrix map[string]Elevator_states //skal key være int? feks 3 siste element i ip-adresse?
    for {
        select {
        case str := <-recv_ch:
            fmt.Println("Got message: ", str)
            //switch msg.MsgType {
            //case MSG_NEW_ORDER:
            //    // ...
            //case MSG_ORDER_DONE:
            //    // ...
            //case MSG_ORDER_ACK:
            //    // ...
            //}
        //case update := <-peer_update_ch:
        //    // add to matrix
        case <-timer:
            sender_ch <- "This is my message!"
        }
    }
}

func order(){
	for {
    	select {
        	case str := <-recv_ch:
            	fmt.Println("Got message: ", str)
	     		switch msg.MsgType {
	            case MSG_NEW_ORDER:
	                // hvis får ny ordre fra. nr 2, må han få svar fra andre gjenlevende før han oppdaterer matrix
	            case MSG_ORDER_DONE:
	                // ...
	            case MSG_ORDER_ACK:
	                // ...
	            }
	        case update := <-peerUpdateCh:
	            // add to matrix
	            // new_key = update.New[12:15]
	            state_matrix[new_key] =
	        case <-timer:
	            sender_ch <- "This is my message!"
	        }
	    }
	}


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

/*
func Should_stop() bool{
	//fmt.Println("should stop?")

	//må sjekke om det er 1-tall i gjeldende etasje i state-matrix. Er det det må internal_orders sjekkes, er det en internal må
	// heisen stoppe. Er det ikke internal må external_order sjekkes. Er current_direction samme vei som bestillingen skal heisen stoppe
	// husk at når heisen stopper i en etg hopper alle på, så da slettes alle ordre i den etg til respektive heis (skjer i delete_orders())
	if State_matrix[id].Floors[current_floor] == 1 {
		if Internal_orders[id][current_floor] == 1 {
			return true
		} else if current_dir == 1 {//&& (External_order[current_floor][1] == 1) {
			return true
		} else if current_dir == -1 { //&& (External_order[current_floor][0] == 1) {
			return true
		}
	}
	return false
}*/

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

/*
func Delete_orders() {
	floor := State_matrix[id].Current_floor
	Internal_orders[id][floor] = 0
	//State_matrix[id] = Elevator_states{Floor_1: 1} MÅ FÅ TIL LISTE TIL FLOOR!!!
}*/
