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
The PeerUpdate struct gives information about which peers being disconnected or not to the nwtwork and newcomers to the system. 
This is used for error handling and for updating the elevators when they return to normal operation and reconnects with the others. 





