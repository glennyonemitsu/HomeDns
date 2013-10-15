package main

import (
	"net"
)

type database struct {
	records		map[string]*record
}

func (db *database) findById(id string) (rec *record, found bool) {
	rec, found = db.records[id]
	return rec, found
}

func (db *database) setRecord(id string, ipv4 net.IP, ipv6 net.IP, timeToLive int32, password string) {
	record := new(record)
	record.id = id
	record.ipv4 = ipv4
	record.ipv6 = ipv6
	record.timeToLive = timeToLive
	db.records[id] = record
}
