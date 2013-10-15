package main

import (
	"fmt"
	"net"
	"net/http"
)

var db database

func main() {
	http.HandleFunc("/", router)
	err := http.ListenAndServe("0.0.0.0:8081", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func init() {
	db = database{records: make(map[string]record)}
	db.records["glenn"] = record{id: "glenn", ipv4: net.ParseIP("127.0.0.1"), timeToLive: 3600, password: "password"}
	db.records["platters"] = record{id: "platters", ipv4: net.ParseIP("192.168.1.124"), timeToLive: 3600, password: "password"}
}
