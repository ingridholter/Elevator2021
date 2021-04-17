# MAIN

In this package the Init function called. Then all the channels and threads are initialized. 

We mainly used channels for communications between threads/modules and use golangs goroutines for concurrent programming. 

NewOrderCh and LightsOfflineCh are the only channels used for communications between the two main gorotines, RunElevator and SyncAllElevators.

The functions for polling buttons and other hardware signals are spawned. And the functions for network communications is also spawned. They all run concurrently and the last select{} is used to stop main from exiting. The channels used in these goroutines are also used in the two main gorotines, RunElevator and SyncAllElevators.
