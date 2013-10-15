package main

import (
	"errors"
)

type database struct {
	records		map[string]record
}

func (db *database) findById(id string) (*record, error) {
	if rec, ok := db.records[id]; ok {
		return &rec, nil
	} else {
		return nil, errors.New("Not found")
	}

}

