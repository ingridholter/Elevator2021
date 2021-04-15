# Elevator21
Elevator project w/networking

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

sudo iptables -A INPUT -p PROTOCOL --dport 20009 -m statistic --mode random --probability 0.2 -j DROP

sudo iptables -A INPUT -p PROTOCOL --dport 20007 -m statistic --mode random --probability 0.2 -j DROP



