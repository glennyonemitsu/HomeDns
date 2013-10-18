package main

import (
	"fmt"
	"net/http"
	"flag"
	"net"
	"strings"
	"strconv"
)

var db map[string]*record
var flagBindHttp *string
var flagBindDns *string
var flagPassword *string

type record struct {
	name	string
	ttl		int64
	ipv4	net.IP
}

func init() {
	db = make(map[string]*record)

	flagBindHttp = flag.String("bind_http", "0.0.0.0:8053", "IP:PORT to bind HTTP server to")
	flagBindDns = flag.String("bind_dns", "0.0.0.0:53", "IP:PORT to bind DNS server to")
	flagPassword = flag.String("password", "", "password to be used for HTTP client authentication")
	flag.Parse()

}

func main() {
	http.HandleFunc("/", router)
	err := http.ListenAndServe(*flagBindHttp, nil)
	if err != nil {
		fmt.Println(err)
	}
}

func router(res http.ResponseWriter, req *http.Request) {
	if password := req.FormValue("password"); password == *flagPassword {
		rec := new(record)
		rec.name = req.FormValue("name")
		rec.ttl, _ = strconv.ParseInt(req.FormValue("ttl"), 10, 64)
		rec.ipv4 = net.ParseIP(ipFromAddress(req.RemoteAddr))
		db[rec.name] = rec
	}
}

func ipFromAddress(address string) string {
	parts := strings.Split(address, ":")
	return parts[0]
}
