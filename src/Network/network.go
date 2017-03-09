package Network

import (
	"./network/bcast"
	"./network/localip"
	"./network/peers"
	//"flag"
	"fmt"
	//"os"
)

// We define some custom struct to send over the network.
// Note that all members we want to transmit must be public. Any private members
//  will be received as zero-values.
type HelloMsg struct {
	Message string
	Iter    int
	//her vil vi sende over en et element i matrisen eller hele matrisen
}

func GetLocalId() string{
	localIP, err := localip.LocalIP()
	if err != nil {
		return "DISCONNECTED"
	}
	return localIP[12:15]
}

func Network(local_id string, sender_ch, recv_ch chan string) {
	helloTx := make(chan HelloMsg)
	helloRx := make(chan HelloMsg)
	sendTx := make(chan bool) //kanal som (acc) = accnolage som svarer om en har fått melding
	sendRx := make(chan bool)

	go bcast.Transmitter(16585, helloTx, sendTx, sender_ch)
	go bcast.Receiver(16585, helloRx, sendRx, recv_ch)
	// ... or alternatively, we can use the local IP address.
	// (But since we can run multiple programs on the same PC, we also append the
	//  process ID)
	if local_id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}

		//local_id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
		fmt.Println("Ip inni network", localIP)
		//fmt.Println("Er inni nettverk, prøver å printe id")
		id := GetLocalId()
		//fmt.Println("Id: ",id)
		sender_ch <- id
		//fmt.Println("Har sendt id på kanal sender_ch")
	}
	
	// We make a channel for receiving updates on the id's of the peers that are
	//  alive on the network
	peerUpdateCh := make(chan peers.PeerUpdate)
	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	peerTxEnable := make(chan bool)
	go peers.Transmitter(20019, local_id, peerTxEnable)
	go peers.Receiver(20019, peerUpdateCh)

	// We make channels for sending and receiving our custom data types

	// ... and start the transmitter/receiver pair on some port
	// These functions can take any number of channels! It is also possible to
	//  start multiple transmitters/receivers on the same port.
	

	// The example message. We just send one of these every second.
	//go func() {
	//	helloMsg := HelloMsg{"Hello from " + id, 0}
	//	for {
	//		helloMsg.Iter++
	//		helloTx <- helloMsg
	//		time.Sleep(1 * time.Second)
	//	}
	//}()

	fmt.Println("Started")
	for {
		select {
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

		case a := <-helloRx:
			fmt.Printf("Received: %#v\n", a)
		}
	

	}
}
