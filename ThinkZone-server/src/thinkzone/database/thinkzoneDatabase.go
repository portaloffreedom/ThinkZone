// database
package database

import (
//"fmt"
)

type User struct {
	ID       int
	Username string
}

type databaseRegistration struct {
	UserNameToId map[string]*User
	UserIDtoUser map[int]*User
	contatore    int //sostituire con lista degli id non usati
}

var Data databaseRegistration = databaseRegistration{make(map[string]*User), make(map[int]*User), 0}

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
		user.ID = datab.contatore
		user.Username = s
		datab.UserNameToId[s] = user
		datab.UserIDtoUser[user.ID] = user

		newuser = true
	} else {
		newuser = false
	}

	return

}

/*
//deprecated
func (datab *databaseRegistration) AddUserId(s string) (int, bool) {
	ID := datab.UserNameToId[s]

	if ID == 0 {
		datab.contatore++
		ID = datab.contatore
		datab.UserNameToId[s] = ID
		return ID, true
	} else {
		//TODO error already exists
		return 0, false
	}
	return 0, false
}*/

func (datab *databaseRegistration) GetUserByName(s string) *User {
	return datab.UserNameToId[s]
}

func (datab *databaseRegistration) GetUserByID(id int) *User {
	return datab.UserIDtoUser[id]
}
