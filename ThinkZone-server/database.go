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
//TODO sostituire questa funzione con 2 diverse: login e createUser
func (datab *databaseRegistration) ConnectUser(s string) (user *User, newuser bool) {
	//user = datab.UserNameToId[s]
	var present bool

	if user, present = datab.UserNameToId[s]; !present {
		datab.contatore++
		user = new(User)
		user.id = datab.contatore
		user.username = s
		datab.UserNameToId[s] = user

		newuser = true
	} else {
		newuser = false
	}

	return

}

/*
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
}*/

func (datab *databaseRegistration) GetUserId(s string) int {
	user := datab.UserNameToId[s]
	return user.id
}
