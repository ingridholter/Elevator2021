sudo iptables -A INPUT -p tcp --dport 15657 -j ACCEPT

sudo iptables -A INPUT -p tcp --sport 15657 -j ACCEPT

sudo iptables -A INPUT -p udp --dport 20007 -m statistic --mode random --probability 0.4 -j DROP
sudo iptables -A INPUT -p udp --dport 20008 -m statistic --mode random --probability 0.4 -j DROP
sudo iptables -A INPUT -p udp --dport 20009 -m statistic --mode random --probability 0.4 -j DROP