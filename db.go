package main


type database struct {
	records		map[string]record
}

func (db *database) findById(id string) (rec record, found bool) {
	rec, found = db.records[id]
	return rec, found
}

