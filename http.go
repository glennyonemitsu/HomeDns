/**
API handlers for http endpoints
*/

package main

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"strconv"
	"errors"
)

/**
 These are parts of the URL split by "/"
 */
const (
	_ = iota // 0 is skipped since split always has first item empty
	apiRecordId
	apiRecordProperty
)

func router(res http.ResponseWriter, req *http.Request) {
	id, property, badURL := endpoint(req.URL)
	record, found := db.findById(id)
	switch req.Method {
	case "GET":
		if badURL != nil || found == false {
			res.WriteHeader(404)
		} else {
			switch property {
			case "ipv4":
				fmt.Fprint(res, record.ipv4)
			case "ipv6":
				fmt.Fprint(res, record.ipv6)
			case "id":
				fmt.Fprint(res, record.id)
			case "ttl":
				fmt.Fprint(res, record.timeToLive)
			case "":
				fmt.Fprintln(res, record.outputText())
			default:
				res.WriteHeader(404)
			}
		}
	case "PUT":
		var ipv4ip, ipv6ip net.IP

		password := req.PostFormValue("password")
		ipv4 := req.PostFormValue("ipv4")
		ipv6 := req.PostFormValue("ipv6")
		ttl, _ := strconv.ParseInt(req.PostFormValue("ttl"), 10, 32)
		timeToLive := int32(ttl)
		if timeToLive > 3600 {
			timeToLive = 3600
		}
		clientIp := getClientIP(req.RemoteAddr)

		ipv4ip = net.ParseIP(ipv4)
		if ipv4ip == nil && isIpv4(clientIp) {
			ipv4ip = net.ParseIP(clientIp)
		}

		ipv6ip = net.ParseIP(ipv6)
		if ipv6ip == nil && isIpv6(clientIp) {
			ipv6ip = net.ParseIP(clientIp)
		}

		if found {
		} else {
			db.setRecord(id, ipv4ip, ipv6ip, timeToLive, password)
		}
		fmt.Fprint(res, req.Header)
		fmt.Fprint(res, req.RemoteAddr)
	}
}

func endpoint(url *url.URL) (id string, property string, err error) {
	id = ""
	property = ""
	err = nil

	parts := strings.Split(url.Path, "/")
	if l := len(parts); l == 2 {
		err = errors.New("Invalid URL")
	} else if l == 3 {
		id = parts[apiRecordId]
	} else if l == 4 {
		id = parts[apiRecordId]
		property = parts[apiRecordProperty]
	}
	return id, property, err

}

func getClientIP(addr string) (ip string) {
	splitIndex := strings.LastIndex(addr, ":")
	if splitIndex == -1 {
		ip = addr
	} else {
		ip = addr[:splitIndex]
	}
	return
}

func isIpv4(ip string) bool {
	if strings.Count(ip, ".") > 0 {
		return true
	} else {
		return false
	}
}

func isIpv6(ip string) bool {
	if strings.Count(ip, ":") > 0 {
		return true
	} else {
		return false
	}
}
