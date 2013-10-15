package main

import (
	"time"
	"net"
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
