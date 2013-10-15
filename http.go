/**
API handlers for http endpoints
*/

package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
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
	id, property, err := endpoint(req.URL)
	switch req.Method {
	case "GET":
		if record, found := db.findById(id); err != nil || found == false {
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
		res.WriteHeader(404)
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
