// database
package database

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"thinkzone/logs"
	"time"

//	"thinkzone/network"
)

// Struttura dati che memorizza i dati dell'utente
type User struct {
	ID       int    // Codice intero univoco identificativo dell'utente
	Username string // Username
	password []byte // Hash della password vera
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
	ServerFakeUser *User = nil
	// Conversazione principale attiva (da eliminare quando vengono implementate
	// più conversazioni per server
	MainConv *Conversation = nil
	// Database principale in cui vengono resgistrati tutti gli utenti
	Data DatabaseRegistration = DatabaseRegistration{make(map[string]*User), make(map[int]*User), 0}

	// Database in cui registrare i dati del server
	db *sql.DB
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
	success = false

	if user, present = datab.userNameToId[username]; !present {
		datab.contatore++
		user = new(User)
		user.ID = datab.contatore
		user.Username = username
		hashpassword := sha256.New()
		hashpassword.Write([]byte(password))
		user.password = hashpassword.Sum([]byte{})

		err := salvaUtenteSQL(user)
		if err != nil {
			logs.Error(err.Error())
			return nil, false
		}

		datab.userNameToId[username] = user
		datab.userIDtoUser[user.ID] = user

		success = true

	} else {
		return nil, false
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
	//	fmt.Println("paragone password:")
	//	fmt.Printf("input:_%v_\n", (hashinput.Sum([]byte{})))
	//	fmt.Printf("datab:_%v_\n", (user.password))
	return bytes.Equal(hashinput.Sum([]byte{}), user.password)
}

// Inizializza tutte le operazioni necessarie per sul database
func init() {
	logs.Log("init del database")

	ServerFakeUser = new(User)
	ServerFakeUser.ID = 0
	ServerFakeUser.Username = "server"
	serverPassword := sha256.New()
	serverPassword.Write([]byte("toor"))
	ServerFakeUser.password = serverPassword.Sum([]byte{})

	err := CreateDataBaseSQL()
	if err != nil {
		logs.Error(err.Error())
	}

	//	MainConv = NewConversation(ServerFakeUser)
	err = Data.caricaConversazioni()
	if err != nil {
		logs.Error("Impossibile caricare le conversazioni vecchie\nmotivo: ", err.Error())
	}

	//	err = salvaUtente(ServerFakeUser)
	//	if err != nil {
	//		logs.Error("Impossibile salvare l'utente server nel database\nmotivo: ", err.Error())
	//	}()

	Data.CaricaUtentiSQL()
	if err != nil {
		logs.Error("Impossibile caricare gli utenti dal database\nmotivo: ", err.Error())
	}

	go func() {
		var spegniti bool = false
		logs.AggiungiAzioneDiChiusura(func() { spegniti = true })

		tk := time.NewTicker(1 * time.Minute)
		for !spegniti {
			<-tk.C
			err := MainConv.salvaTutteLeConversazioniSQL()
			if err != nil {
				logs.Error("Impossibile salvare tutti i post\nmotivo: ", err.Error())
			}
			logs.Log("salvate tutte le conversazioni sul database")
			//err := MainConv.salvaTuttiIPostSQL(&Data)
			//if err != nil {
			//	logs.Error("Impossibile salvare tutti i post\nmotivo: ", err.Error())
			//}
		}
		tk.Stop()
	}()
}
