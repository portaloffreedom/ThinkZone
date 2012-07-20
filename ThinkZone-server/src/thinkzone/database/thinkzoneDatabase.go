// database
package database

import (
	//"fmt"
	"bytes"
	"crypto/sha256"
	"hash"
)

// Struttura dati che memorizza i dati dell'utente
type User struct {
	ID       int       // Codice intero univoco identificativo dell'utente
	Username string    // Username
	password hash.Hash // Hash della password vera
}

// Struttura dati core del dababase: tiene memorizzati tutti gli utenti
type DatabaseRegistration struct {
	userNameToId map[string]*User // mappa di tutti gli utente per username
	userIDtoUser map[int]*User    // mappa di tutti gli utenti per ID
	contatore    int              // sostituire con lista degli id non usati
}

// variabili globali per il funzionamento del Database
var (
	//Utente fake che rappresenta il server
	ServerFakeUser User = User{0, "server", sha256.New()}
	// Conversazione principale attiva (da eliminare quando vengono implementate
	// più conversazioni per server
	MainConv *Conversation = NewConversation(&ServerFakeUser)
	// Database principale in cui vengono resgistrati tutti gli utenti
	Data DatabaseRegistration = DatabaseRegistration{make(map[string]*User), make(map[int]*User), 0}
)

// Questa funzione connette un nuovo utente al database (in pratica ne verifica 
// solamente la già avvenuta registrazione) e ritorna il puntatore alla struttura
// dati dell'utente con tutti i dati necessari (id e hash della password)
func (datab *DatabaseRegistration) ConnectUser(s string) (user *User, newuser bool) {
	//user = datab.userNameToId[s]
	var present bool

	user, present = datab.userNameToId[s]
	newuser = !present

	return

}

// Questa funzione registra un nuovo utente nel database
func (datab *DatabaseRegistration) RegisterNewUser(username, password string) (user *User, success bool) {

	var present bool

	if user, present = datab.userNameToId[username]; !present {
		datab.contatore++
		user = new(User)
		user.ID = datab.contatore
		user.Username = username
		hashpassword := sha256.New()
		hashpassword.Write([]byte(password))
		user.password = hashpassword

		datab.userNameToId[username] = user
		datab.userIDtoUser[user.ID] = user

		success = true
	} else {
		success = false
		user = nil
	}

	return

}

// Trasforma il nome dell'utente nel puntatore alla sua struttura dati
// (da per scontato che la struttura esista)
func (datab *DatabaseRegistration) GetUserByName(s string) *User {
	return datab.userNameToId[s]
}

// Trasforma l'id dell'utente nel puntatore alla sua struttura dati
// (da per scontato che la struttura esista)
func (datab *DatabaseRegistration) GetUserByID(id int) *User {
	return datab.userIDtoUser[id]
}

// Verifica che la password dell'utente sia corretta rispetto a quella passata
// come parametro (internamente le password sono trasformate, verificate
// memorizzate come checksum)
//
// Ritorna true se la password è corretta
func (user *User) VerifyPassword(passwordInput string) bool {
	hashinput := sha256.New()
	hashinput.Write([]byte(passwordInput))
	return bytes.Equal(hashinput.Sum([]byte{}), user.password.Sum([]byte{}))
}
