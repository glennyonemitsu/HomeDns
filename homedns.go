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
	db = database{records: make(map[string]*record)}
	db.setRecord("glenn", net.ParseIP("127.0.0.1"), net.ParseIP(""), 3600, "password")
	db.setRecord("platters", net.ParseIP("192.168.1.124"), net.ParseIP(""), 3600, "password")
}
