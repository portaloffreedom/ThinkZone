// database
package main

import (
//"fmt"
)

type UserStatus struct {
	id        int
	connected bool
}

type databaseConversation struct {
	UserNameToId map[string]*UserStatus
	contatore    int
}

var data databaseConversation = databaseConversation{make(map[string]*UserStatus), 0}

//this function assign an id to an username
//return the id assigned
//return false if the id already existed, true if a new one was assigned
func (datab *databaseConversation) connectUser(s string) (id int, newuser bool) {
	user := datab.UserNameToId[s]

	if user == nil {
		datab.contatore++
		user = new(UserStatus)
		user.id = datab.contatore
		user.connected = true
		datab.UserNameToId[s] = user
		//return id, true
		newuser = true
	} else {
		//return id, false
		user = testuser
		if ...//TODO ripensare tutta questa parte...
		newuser = false
	}

	return

}

//deprecated
func (datab *databaseConversation) AddUserId(s string) (int, bool) {
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

func (datab *databaseConversation) GetUserId(s string) int {
	id := datab.UserNameToId[s]
	return id
}
