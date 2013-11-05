package main

import (
	"fmt"
	"flag"
	"net"
	"strings"
	"strconv"
)

var db map[string]*record

// server mode flags
var flagBind *string
var flagPassword *string

// client mode flags
var flagServer *string
var flagHostname *string
var flagIpv4 *string
var flagTtl *string

var homeDNSPrefix string
var homeDNSPrefixLength int

type record struct {
	name	string
	ttl		int64
	ipv4	net.IP
}

func init() {
	db = make(map[string]*record)

	flagBind = flag.String("bind", "", "IP:PORT to bind to")
	flagPassword = flag.String("password", "", "password to be used for authentication")
	flagServer = flag.String("server", "", "server address to send updated record to")
	flagHostname = flag.String("hostname", "", "A record to update")
	flagIpv4 = flag.String("ipv4", "", "IP address of the record. Leave alone to use IP of the client")
	flagTtl = flag.String("ttl", "3600", "TTL setting for the DNS record")
	flag.Parse()

	homeDNSPrefix = "HOMEDNS;" + *flagPassword + ";"
	homeDNSPrefixLength = len(homeDNSPrefix)

}

func main() {
	if *flagBind != "" {
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
	} else if *flagServer != "" {
		var msg string
		if *flagIpv4 != "" {
			msg = fmt.Sprintf("%s%s;%s;%s;", homeDNSPrefix, *flagHostname, *flagTtl, *flagIpv4)
		} else {
			msg = fmt.Sprintf("%s%s;%s;", homeDNSPrefix, *flagHostname, *flagTtl)
		}
		ip, port := ipPort(*flagServer)
		netIP := net.ParseIP(ip)
		addr := net.UDPAddr{IP: netIP, Port: port, Zone: ""}
		conn, err := net.DialUDP("udp", nil, &addr)
		defer conn.Close()
		if err != nil {
			fmt.Println(err)
		}
		conn.Write([]byte(msg))
	}
}

func ipPort(address string) (string, int) {
	parts := strings.Split(address, ":")
	ip := parts[0]
	var port int
	if len(parts) == 2 {
		port64, _ := strconv.ParseInt(parts[1], 10, 0)
		port = int(port64)
	} else {
		port = 53
	}
	return ip, port
}

/**
 *
 */
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
			fmt.Println("Storing", rec)
		} else {
			fmt.Printf("%x", input[0:size])
		}
	}
}
