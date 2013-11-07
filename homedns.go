package main

import (
	"encoding/binary"
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
	// server mode
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
	// client mode
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

func handleUDPConnection(conn *net.UDPConn) {
	input := make([]byte, 1024)
	size, source, err := conn.ReadFromUDP(input[0:])
	packet := input[0:size]

	if err != nil {
		fmt.Println(err)
	} else {
		// update packet from HomeDns client
		if string(packet[0:homeDNSPrefixLength]) == homeDNSPrefix {
			message := string(packet[homeDNSPrefixLength:size])
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
		// DNS query packet
		} else if isQueryPacket(packet) {
			query := parseQueryPacket(packet)
			fmt.Println("query", query)
			if record, has := db[query.subdomain]; has == true {
				fmt.Println(record)
			} else {
				fmt.Printf("record for '%s' not found\n", query.subdomain)
				fmt.Println(db)
			}
		}
	}
}

type query struct {
	qid int16
	subdomain string
}

func isQueryPacket(packet []byte) bool {
	flag := "\x01\x00"
	tail := "\x00\x00\x01\x00\x01"
	return string(packet[2:4]) == flag && string(packet[len(packet)-len(tail):]) == tail
}

func parseQueryPacket(packet []byte) *query {
	q := new(query)
	flagInt, _ := binary.Uvarint(packet[0:2])
	q.qid = int16(flagInt)
	hostname := packet[12:len(packet) - 5]
	subdomainLength := int(hostname[0])
	q.subdomain = string(hostname[1:subdomainLength + 1])
	return q
}
