# Network module for Go (UDP) 

This network module is written by [TTK4145](https://github.com/TTK4145)
and can be found [here](https://github.com/TTK4145/Network-go). 

Features
--------

We use broadcasting with UDP for sending and receiving messages between elevators. By using broadcasting all peers receive all messages, even the sender itself. 

This is taken into account by adding the function AcceptNewOrder which ensures that the right receiver recieves message. See the elevatorObserver module for more.

The peers package is used for detecting Network errors and software crash. 

```
type PeerUpdate struct {
	Peers []string
	New   string
	Lost  []string
}

```
The PeerUpdate struct gives information about which peers being available and not and newcomers to the system. 
This is used for error handling and for updating the elevators when they return to normal operation and reconnects with the others. 




FJERNE????: FRA PER
Channel-in/channel-out pairs of (almost) any custom or built-in datatype can be supplied to a pair of transmitter/receiver functions. Data sent to the transmitter function is automatically serialized and broadcasted on the specified port. Any messages received on the receiver's port are deserialized (as long as they match any of the receiver's supplied channel datatypes) and sent on the corresponding channel. See [bcast.Transmitter and bcast.Receiver](network/bcast/bcast.go).

Peers on the local network can be detected by supplying your own ID to a transmitter and receiving peer updates (new, current and lost peers) from the receiver. See [peers.Transmitter and peers.Receiver](network/peers/peers.go).

Finding your own local IP address can be done with the [LocalIP](network/localip/localip.go) convenience function, but only when you are connected to the internet.


*Note: This network module does not work on Windows: Creating proper broadcast sockets on Windows is just not implemented yet in the Go libraries. See issues listed in [the comments here](network/conn/bcast_conn.go) for details.*
