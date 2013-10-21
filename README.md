# Your own personal dynamic dns service

HomeDns is a very simple dynamic DNS service.

## Usage

### Starting the server

Start the server with a password specified.

    HomeDns -password=somesecret

Add a dynamic DNS entry to the server by sending a specially formatted UDP 
packet on port 53. It must follow the format:

    HOMEDNS;<your password>;<hostname>;<ttl>;<ip (optional)>;

If no IP is provided it will use the IP of the UDP client. You can use the 
netcat utility to send a UDP packet via shell.

    echo "HOMEDNS;somesecret;myhome;3600;" | nc -q 1 -u mypersonaldyndns.example.com 53

An A record query will return the proper DNS response (not implemented yet).
