#SyncAllElevators



#Variables Arrays and Matrices

This module uses an array of Elevstates for all the available elevators called allElevators. 
The Elevstate type has information about the state of the elevator and its requests. 
````
type ElevState struct {
	Floor     int
	Dir       MotorDirection
	Behaviour ElevBehaviour
	Requests  [NumFloors][NumButtons]bool
}
```
Our re

Elevator State | Elevator 1 | Elevator 2 | Elevator 3
--------------- | ---------- | ---------- | ----------
Floor 4 | :arrow_forward: / :zzz: / :clock10: | :arrow_forward: / :zzz: / :clock10: |  :arrow_forward: / :zzz: / :clock10:
Floor 3     | :arrow_forward: / :zzz: / :clock10: | :arrow_forward: / :zzz: / :clock10: | :arrow_forward: / :zzz: / :clock10:
Floor 2     | :arrow_forward: / :zzz: / :clock10: | :arrow_forward: / :zzz: / :clock10: | :arrow_forward: / :zzz: / :clock10:
Floor 1     | :arrow_forward: / :zzz: / :clock10: | :arrow_forward: / :zzz: / :clock10: |  :arrow_forward: / :zzz: / :clock10:
