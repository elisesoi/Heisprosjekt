package orders

import (
	"timer"
)

func order(sender_ch, recv_ch chan string) { //skal denne sende over Elevator_states? og ikke string?
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
	if //svar_fra_alle_heiser_over_acc_channel{
		//sett tall i state_matrix
		//sett tall i orders
		//tenn lys / kall lysfunksjon
	}
}

func should_stop(){

}

func choose_direction(){

}

func delete_orders(){

}