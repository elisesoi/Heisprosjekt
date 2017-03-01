//ta inspirasjon fra 
//elev.c. skriv det en vil ha fra elev.c om til go

package driver

import (
		//"os" //errorhandling
		"fmt"
		//"time"
		//"defines"
		)

//define MOTOR_SPEED 2800

var lamp_channel_matrix  = [N_FLOORS][N_BUTTONS] int{
		{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
    	{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
    	{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
    	{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
	}

var	button_channel_matrix = [N_FLOORS][N_BUTTONS] int{
    	{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
    	{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
    	{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
    	{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
	}

//Her er det en feil
func Driver_init() error{
	init_success := io_init()
    if init_success == 0 {
        return fmt.Errorf("Unable to initialize elevator hardware!")
    }
	//assert(init_success && "Unable to initialize elevator hardware!")


	for f := 0; f < N_FLOORS; f++ {
		for b := 0; b < N_BUTTONS; b++ {
			Driver_set_button_lamp(Driver_button_type(b),f,0)
		}
	}
	Driver_set_stop_lamp(0)
	Driver_set_door_open_lamp(0)
	Driver_set_floor_indicator(0)
	//Spør de andre etter oppdateringer
    return nil
}


func Driver_set_motor_direction(dirn Driver_motor_dir) {
	if dirn == 0 {
		io_write_analog(MOTOR, 0)
	} else if dirn > 0 {
		io_clear_bit(MOTORDIR)
		io_write_analog(MOTOR, MOTOR_SPEED)
	} else if dirn < 0 {
		io_set_bit(MOTORDIR)
		io_write_analog(MOTOR, MOTOR_SPEED)
	}

}

func Driver_set_button_lamp(button Driver_button_type, floor int, value int) error {
    if floor < 0 || floor > N_FLOORS  { //|| button < 0 || button > N_BUTTONS
        return fmt.Errorf("Wrong input: %d")
    }

    if value == 1 { //endret til == 1. For å sette lampen sett value = 1
        io_set_bit(lamp_channel_matrix[floor][button])
    } else {
        io_clear_bit(lamp_channel_matrix[floor][button])
    }
    return nil
}


func Driver_set_floor_indicator(floor int) error {
    if floor < 0 || floor > N_FLOORS{
        return fmt.Errorf("Wrong input to floor_indicator: %d")
    }

    // Binary encoding. One light must always be on.
    if floor & 0x02 != 0 {
        io_set_bit(LIGHT_FLOOR_IND1)
    } else {
        io_clear_bit(LIGHT_FLOOR_IND1)
    }    

    if floor & 0x01 != 0 {
        io_set_bit(LIGHT_FLOOR_IND2)
    } else {
        io_clear_bit(LIGHT_FLOOR_IND2)
    }    
    return nil
}


func Driver_set_door_open_lamp(value int) {
    if value == 1 {
        io_set_bit(LIGHT_DOOR_OPEN)
    } else {
        io_clear_bit(LIGHT_DOOR_OPEN)
    }
}

func Driver_set_stop_lamp(value int) {
    if value == 1 {
        io_set_bit(LIGHT_STOP)
    } else {
        io_clear_bit(LIGHT_STOP)
    }
}


func Driver_get_button_signal(button Driver_button_type, floor int) int {

    if floor < 0 || floor > N_FLOORS || button < 0 || int(button) > N_BUTTONS{
        fmt.Println("Wrong floor input: %d")
    }

    return io_read_bit(button_channel_matrix[floor][button])
}


func Driver_get_floor_sensor_signal() int {
    if io_read_bit(SENSOR_FLOOR1) == 1 {
        return 0
    } else if io_read_bit(SENSOR_FLOOR2) == 1 {
        return 1
    } else if io_read_bit(SENSOR_FLOOR3) == 1 {
        return 2
    } else if io_read_bit(SENSOR_FLOOR4) == 1 {
        return 3
    } else {
        return -1
    }
}


func Driver_get_stop_signal() int {
    return io_read_bit(STOP)
}

/* kanksje vi ikke trenger å gjøre noe med obstruction?
func driver_get_obstruction_signal() int {
    return io_read_bit(OBSTRUCTION)
}
*/