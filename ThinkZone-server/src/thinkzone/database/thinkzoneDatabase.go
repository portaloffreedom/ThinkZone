// database
package database

import (
	//"fmt"
	"bytes"
	"crypto/sha256"
	"hash"
)

type User struct {
	ID       int
	Username string
	password hash.Hash
}

type databaseRegistration struct {
	UserNameToId map[string]*User
	UserIDtoUser map[int]*User
	contatore    int //sostituire con lista degli id non usati
}

var (
	ServerFakeUser User                 = User{42, "server", sha256.New()}
	MainConv       *Conversation        = NewConversation(&ServerFakeUser)
	Data           databaseRegistration = databaseRegistration{make(map[string]*User), make(map[int]*User), 0}
)

//this function assign an id to an username
//return the id assigned
//return false if the id already existed, true if a new one was assigned
//TODO sostituire questa funzione con 2 diverse: login e createUser
func (datab *databaseRegistration) ConnectUser(s string) (user *User, newuser bool) {
	//user = datab.UserNameToId[s]
	var present bool

	if user, present = datab.UserNameToId[s]; !present {
		newuser = true
	} else {
		newuser = false
	}

	return

}

func (datab *databaseRegistration) RegisterNewUser(username, password string) (user *User, success bool) {

	var present bool

	if user, present = datab.UserNameToId[username]; !present {
		datab.contatore++
		user = new(User)
		user.ID = datab.contatore
		user.Username = username
		hashpassword := sha256.New()
		hashpassword.Write([]byte(password))
		user.password = hashpassword

		datab.UserNameToId[username] = user
		datab.UserIDtoUser[user.ID] = user

		success = true
	} else {
		success = false
		user = nil
	}

	return

}

func (datab *databaseRegistration) GetUserByName(s string) *User {
	return datab.UserNameToId[s]
}

func (datab *databaseRegistration) GetUserByID(id int) *User {
	return datab.UserIDtoUser[id]
}

func (user *User) VerifyPassword(passwordInput string) bool {
	hashinput := sha256.New()
	hashinput.Write([]byte(passwordInput))
	return bytes.Equal(hashinput.Sum([]byte{}), user.password.Sum([]byte{}))
}
