package main

import (
	"time"
	"net"
	"fmt"
)

type record struct {
	id             string
	timeToLive     int32
	ipv4           net.IP
	ipv6           net.IP
	password       string
	expirationDate time.Time
}

func (rec *record) property(name string) {
}


func (rec *record) outputText() string {
	var ipv4, ipv6 string
	format := "%s\n%s\n%s\n%d"
	if ipv4 = rec.ipv4.String(); ipv4 == "<nil>" {
		ipv4 = ""
	}
	if ipv6 = rec.ipv6.String(); ipv6 == "<nil>" {
		ipv6 = ""
	}
	output := fmt.Sprintf(format, rec.id, ipv4, ipv6, rec.timeToLive)
	return output
}


