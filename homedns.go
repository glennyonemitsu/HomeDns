package main

import (
	"fmt"
	"flag"
	"net"
	"strings"
	"strconv"
)

var db map[string]*record
var flagBind *string
var flagPassword *string
var homeDNSPrefix string
var homeDNSPrefixLength int

type record struct {
	name	string
	ttl		int64
	ipv4	net.IP
}

func init() {
	db = make(map[string]*record)

	flagBind = flag.String("bind", "0.0.0.0:53", "IP:PORT to bind to")
	flagPassword = flag.String("password", "", "password to be used for authentication")
	flag.Parse()

	homeDNSPrefix = "HOMEDNS;" + *flagPassword + ";"
	homeDNSPrefixLength = len(homeDNSPrefix)

}

func main() {
	ip, port := ipPort(*flagBind)
	netIP := net.ParseIP(ip)
	addr := net.UDPAddr{IP: netIP, Port: port, Zone: ""}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println(err)
	}
	for {
		handleUDPConnection(conn)
	}
}

func ipPort(address string) (ip string, port int) {
	parts := strings.Split(address, ":")
	ip = parts[0]
	port64, _ := strconv.ParseInt(parts[1], 10, 0)
	port = int(port64)

	return ip, port
}

func handleUDPConnection(conn *net.UDPConn) {
	input := make([]byte, 1024)
	size, source, err := conn.ReadFromUDP(input[0:])
	if err != nil {
		fmt.Println(err)
	} else {
		if string(input[0:homeDNSPrefixLength]) == homeDNSPrefix {
			message := string(input[homeDNSPrefixLength:size])
			messageParts := strings.Split(message, ";")
			rec := new(record)
			rec.name = messageParts[0]
			rec.ttl, err = strconv.ParseInt(messageParts[1], 10, 64)
			if err != nil {
				fmt.Println(err)
				return
			}
			if len(messageParts) > 2 && strings.TrimSpace(messageParts[2]) != "" {
				rec.ipv4 = net.ParseIP(messageParts[2])
			} else {
				rec.ipv4 = source.IP
			}
			db[messageParts[0]] = rec
		} else {
		}
	}
}
