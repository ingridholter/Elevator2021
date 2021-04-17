# Elevator Project

This the code for the elevator project made by group67. For more details se the different folder for more specific README's


## General implementation methods

### Peer-to-peer/fleeting master
We have chosen a peer-to-peer/fleeting master implementation. This means that if a button is pressed on an elevator panel, the same elevator decides who should take that order and broadcasts it over the network. Likewise if another button panel is pressed, that elevator decides who to take the order

### broadcasting with udp 

We use broadcasting with udp, and all messages are sent to every elevator listening on the same ports. 

### packet loss handling

Our system is vulnerable to packet loss. To gurantee that no orders are lost over the network communication we only sync the lights if more than one elevator know about the order. If the message with the order is lost, no lights will be turned on an so a person will have to press the button until the message arrives and they see the light. In addition, if somewthing were to happen to one elevator, the order would always be redistributed and never lost. 

## Minor details

When starting the project a flag -id is given with the terminal command, this should be 0, 1 or 2 when working with three elevators. This is given as a string in our program. We consistetly use id with small i as a string and Id with a big I as an int. 


## Til installering og kontrollering på andre maskiner på sanntidslabben
ElevatorServer
go run -race main.go -id=0 

cargo install --force ttk4145_elevator_server
export PATH="/home/student/.cargo/bin:$PATH"


 ssh student@ip   
 scp -r /home/student/group67/Elevator2021 student@10.100.23.171:/home/student
 
 vår ip: 10.100.23.209

to andre pc: 
ssh student@10.100.23.131
10.100.23.131, scp -r /home/student/group67/Elevator2021 student@10.100.23.131:/home/student

ssh student@10.100.23.139
10.100.23.139, scp -r /home/student/group67/Elevator2021 student@10.100.23.139:/home/student

for pc 157:
ssh student@10.100.23.153
10.100.23.153, scp -r /home/student/group67/Elevator2021 student@10.100.23.157:/home/student

go run main.go -id=0
go run main.go -id=1


update go: https://khongwooilee.medium.com/how-to-update-the-go-version-6065f5c8c3ec




