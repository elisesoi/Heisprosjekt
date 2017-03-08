package main

import (
	"fmt"
	. "../driver"
	//"../Network/network/localip"
	."../Network"
	//"time"
	. "../elevator"
	."../orders"
)


func main() {
	sender_ch := make(chan string)
	recv_ch := make(chan string)
	floor_reached_ch := make(chan int)
	order_new_state_ch := make(chan int)
	new_dir_state_ch := make(chan Driver_motor_dir)
	delete_order_ch := make(chan int)
	new_order_ch := make(chan Driver_button_type)

	fmt.Println("Har laget kanaler i main")
	localid := ""

	go Network(localid, sender_ch, recv_ch)
	local_id := <- recv_ch
	fmt.Println(local_id)

	Initialize_elevator(local_id)
	go Order(order_new_state_ch, new_dir_state_ch, new_order_ch, local_id)
	go Elevator_loop(floor_reached_ch, order_new_state_ch, new_dir_state_ch)
	select {}

	/*
	   Driver_init()

	   fmt.Println("Press STOP button to stop elevator and exit program.\n")

	   Driver_set_motor_direction(DIRN_UP)

	   //kalle nettverk-initialisering
	   sender_ch := make(chan string)
	   recv_ch := make(chan string)
	   //go order(sender_ch, recv_ch)
	   go Network.Network(sender_ch, recv_ch)

	   // this is actually a function: elevator_loop()
	   for {
	   	Driver_set_floor_indicator(Driver_get_floor_sensor_signal())

	       // Change direction when we reach top/bottom floor
	       if Driver_get_floor_sensor_signal() == N_FLOORS - 1 {
	           Driver_set_motor_direction(DIRN_DOWN)
	       } else if Driver_get_floor_sensor_signal() == 0 {
	           Driver_set_motor_direction(DIRN_UP)
	       }

	       // Stop elevator and exit program if the stop button is pressed
	       if (Driver_get_stop_signal() != 0) {
	           Driver_set_motor_direction(DIRN_STOP)
	           //fmt.Println(internal_order)
	           break
	       }
	       // test for order, må gjøre det samme for alle etg og alle knapper. Kanskje i en annen funksjon/modul???
	       if (Driver_get_button_signal(BUTTON_CALL_UP,2) != 0){
	       	orders[2][BUTTON_CALL_UP] = EXTERNAL_ORDER
	       }
	       // test for internal order map, må gjøres for alle etg
	       if (Driver_get_button_signal(BUTTON_COMMAND,2) != 0){
	       	internal_order = map[string]int{} //initialisere mapen slik at en ikke aksesserer tom map
	       	localIP, err := localip.LocalIP();
	       	if err == nil{
	       		internal_order[localIP] = 1
	       		//Ubs. må sette value = en liste med 4 plasser = etasjene
	       	}
	       }
	       //test for state_matrix

	   }
	*/

}