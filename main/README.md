# MAIN

In this package the Init function called. Then all the channels and threads are initialized. 

NewOrderCh and LightsOfflineCh are the only channels used for communications between the two main threads, RunElevator and SyncAllElevators.

The functions for polling buttons and other hardware signals are spawned. And the functions for network communications is also spawned. They all run concurrently and the last select{} is used to stop main from exiting.
