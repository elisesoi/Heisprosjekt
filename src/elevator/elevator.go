package elevator

import "../orders"


func elevator_loop(){
	for {
		//holder state_matrix oppdatert:
		state_matrix[id].Current_floor = Driver_get_floor_sensor_signal //er det endring må det sendes på kanal til de andre
		
	}
}


 
    //vil endre Current_floor i state_matrix hver gang den kommer til ny etg