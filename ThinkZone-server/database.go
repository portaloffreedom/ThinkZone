// database
package main

import (
//"fmt"
)

type database struct {
	UserNameToId map[string]int
	contatore    int
}

var data database = database{make(map[string]int), 0}

func AddUserId(s string) (int, bool) {
	id := data.UserNameToId[s]

	if id == 0 {
		data.contatore++
		id = data.contatore
		data.UserNameToId[s] = id
		return id, true
	} else {
		//TODO error already exists
		return 0, false
	}
	return 0, false
}

func GetUserId(s string) int {
	id := data.UserNameToId[s]
	return id
}
