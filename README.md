# Your own personal dynamic dns service

HomeDns is a very simple dynamic DNS service.

## Usage

### Starting the server

The command runs as a server when the `-bind` and `-password` command flags are
specified

    HomeDns \
		-bind=0.0.0.0:53 \
		-password=somesecret

### Updating the server's record with the HomeDns command

The command runs as a client with the `-server` command flag

	HomeDns \
		-password=somesecret \
		-server=mypersonaldyndns.example.com \
		-ttl=600 \
		-ipv4=1.2.3.4 \
		-name=myhome

This might seem odd, using the `-server` flag to indicate client mode, but that
indicates the server address, while the `-bind` is more common nomenclature for
a server parameter.

#### The UDP packet format to update

The UDP packet to update a DNS record follows this format:

    HOMEDNS;<your password>;<hostname>;<ttl>;<ip (optional)>;

If no IP is provided it will use the IP of the UDP client. You can use the 
netcat utility to send a UDP packet via shell.

    echo "HOMEDNS;somesecret;myhome;3600;" | nc -q 1 -u mypersonaldyndns.example.com 53

An A record query will return the proper DNS response (not implemented yet).

