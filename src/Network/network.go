package Network

import (
	. "../driver"
	"./network/bcast"
	"./network/localip"
	"./network/peers"
	//"flag"
	"fmt"
	//"os"
	//"time"
)

// We define some custom struct to send over the network.
// Note that all members we want to transmit must be public. Any private members
//  will be received as zero-values.
type HelloMsg struct {
	Message string
	Iter    int
	//her vil vi sende over en et element i matrisen eller hele matrisen
}

func GetLocalId() string {
	localIP, err := localip.LocalIP()
	if err != nil {
		return "DISCONNECTED"
	}
	return localIP[12:15]
}

func Network(local_id string, sender_ch, recv_ch, new_peer_ch chan string, lost_peer_ch chan []string, new_state_ch chan Elevator_states) {
	statesTx := make(chan Elevator_states) //kanal som (acc) = accnolage som svarer om en har fått melding
	statesRx := make(chan Elevator_states)

	go bcast.Transmitter(16585, statesTx, sender_ch)
	go bcast.Receiver(16585, statesRx, recv_ch)

	// We make a channel for receiving updates on the id's of the peers that are
	//  alive on the network
	peerUpdateCh := make(chan peers.PeerUpdate)
	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	peerTxEnable := make(chan bool)
	go peers.Transmitter(20319, local_id, peerTxEnable)
	go peers.Receiver(20319, peerUpdateCh)

	fmt.Println("Started")
	for {
		select {
		//case state_update := <- stateUpdateCh:

		case p := <-peerUpdateCh:
			if p.New != "" {
				new_id := p.New
				if new_id != local_id {
					new_peer_ch <- new_id //send ny id på kanal
					fmt.Println("Det er oppdaget en ny heis med id: ", new_id)
				}

			} else if p.Lost[0] != "" {
				lost_id := p.Lost
				lost_peer_ch <- lost_id
				fmt.Println("Vi har mistet kontakt med heis: ", lost_id)
			}

			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

		case state_update_tx := <-new_state_ch:
			//skal bcastes til alle på stateUpdateCh
			//fmt.Println("Dette er state_updates sendt til network via new_state_ch ", state_update_tx)
			statesTx <- state_update_tx

		case state_update_rx := <-statesRx:
			new_state_ch <- state_update_rx
		}
	}
}
