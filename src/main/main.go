package main

import (
	. "../driver"
	//"fmt"
	//"../Network/network/localip"
	. "../Network"
	//"time"
	. "../elevator"
	. "../orders"
)

func main() {

	sender_ch := make(chan string)
	recv_ch := make(chan string)
	new_peer_ch := make(chan string)
	lost_peer_ch := make(chan []string)
	new_state_ch := make(chan Elevator_states)

	floor_reached_ch := make(chan int)
	order_new_state_ch := make(chan int)
	new_dir_state_ch := make(chan Driver_motor_dir)

	new_order_ch := make(chan Order_type)
	delete_order_ch := make(chan Order_type)

	//fmt.Println("Har laget kanaler i main")
	localid := GetLocalId()
	Initialize_elevator(localid)

	go Network(localid, sender_ch, recv_ch, new_peer_ch, lost_peer_ch, new_state_ch)
	go Order(order_new_state_ch, new_dir_state_ch, new_order_ch, delete_order_ch, new_peer_ch, lost_peer_ch, new_state_ch, localid)
	go Elevator_loop(floor_reached_ch, order_new_state_ch, new_dir_state_ch, new_order_ch, delete_order_ch, localid)

	select {}

}
