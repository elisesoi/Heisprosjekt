package main

import (
    //"fmt"
    ."../driver"
    //"../Network/network/localip"
    //"../Network"
    //"time"
    "../elevator"
)



func main() {
    Initialize_elevator()
    Elevator_loop()
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









// vil ha en ordrematrise, som inneholder ordrene. 


var orders  = [N_FLOORS][2] int{      //verdier i denne settes når en knapp trykkes. Sett først et random nr. feks 9, og endre til 1 når en er sikker på at de andre har fått bestillingen de også. Når bestilling slettes settes 0 i ruten som passer. 
        {0, 0},
        {0, 0},
        {0, 0},
        {0, 0},
    }

//var orders map[int]updown_order


    //order map som tar seg av indre ordre
var internal_order map[string]int 

