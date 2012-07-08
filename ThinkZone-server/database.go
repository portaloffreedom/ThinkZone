// database
package main

import (
//"fmt"
)

type User struct {
	id       int
	username string
}

type databaseRegistration struct {
	UserNameToId map[string]*User
	contatore    int //sostituire con lista degli id non usati
}

var data databaseRegistration = databaseRegistration{make(map[string]*User), 0}

//this function assign an id to an username
//return the id assigned
//return false if the id already existed, true if a new one was assigned
func (datab *databaseRegistration) connectUser(s string) (id int, newuser bool) {
	user := datab.UserNameToId[s]

	if user == nil {
		datab.contatore++
		user = new(UserStatus)
		user.id = datab.contatore
		datab.UserNameToId[s] = user
		//return id, true
		newuser = true
	} else {
		//return id, false
		user = testuser
		//if ...//TODO ripensare tutta questa parte...
		newuser = false
	}

	return

}

//deprecated
func (datab *databaseRegistration) AddUserId(s string) (int, bool) {
	id := datab.UserNameToId[s]

	if id == 0 {
		datab.contatore++
		id = datab.contatore
		datab.UserNameToId[s] = id
		return id, true
	} else {
		//TODO error already exists
		return 0, false
	}
	return 0, false
}

func (datab *databaseRegistration) GetUserId(s string) int {
	id := datab.UserNameToId[s]
	return id
}
