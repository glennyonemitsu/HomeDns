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
	fmt.Println(req.URL)
	id, property, err := endpoint(req.URL)
	if record, findErr := db.findById(id); err != nil || findErr != nil || record == nil {
		res.WriteHeader(404)
		fmt.Fprint(res, err)
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
		default:
			res.WriteHeader(404)
		}
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
