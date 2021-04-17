#SyncAllElevators



#Variables Arrays and Matrices

This module uses an array of Elevstates for all the available elevators called elevStateArray. 
The Elevstate type has information about the state of the elevator and its requests. 
````
type ElevState struct {
	Floor     int
	Dir       MotorDirection
	Behaviour ElevBehaviour
	Requests  [NumFloors][NumButtons]bool
}
````
The floors are 0-indexed. 
The Requests matrix is build in this way:

NOTHING| UP | DOWN | CAB
--------------- | ---------- | ---------- | ----------
0 | :arrow_forward: / :zzz: / :clock10: | :arrow_forward: / :zzz: / :clock10: |  :arrow_forward: / :zzz: / :clock10:
1     | :arrow_forward: / :zzz: / :clock10: | :arrow_forward: / :zzz: / :clock10: | :arrow_forward: / :zzz: / :clock10:
2     | :arrow_forward: / :zzz: / :clock10: | :arrow_forward: / :zzz: / :clock10: | :arrow_forward: / :zzz: / :clock10:
3   | :arrow_forward: / :zzz: / :clock10: | :arrow_forward: / :zzz: / :clock10: |  :arrow_forward: / :zzz: / :clock10:
