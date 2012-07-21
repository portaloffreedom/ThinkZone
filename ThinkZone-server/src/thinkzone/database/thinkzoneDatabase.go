// database
package database

import (
	//"fmt"
	"bytes"
	"crypto/sha256"
	"database/sql"
	"fmt"
	_ "github.com/jbarham/gopgsqldriver"
	"strconv"
	"strings"
	"thinkzone/logs"

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

	// Script in SQL per creare le tabelle nel database 
	createDbSqlScript []string = []string{
		"CREATE TABLE t_user ( id INT PRIMARY KEY, username CHAR(64) NOT NULL UNIQUE, password CHAR(256) NOT NULL)",
		"CREATE TABLE conversation ( id INT PRIMARY KEY )",
		"CREATE TABLE archive ( t_user INT, conversation INT, PRIMARY KEY (t_user,conversation), FOREIGN KEY (t_user) REFERENCES t_user(id), FOREIGN KEY (conversation) REFERENCES conversation(id))",
		"CREATE TABLE post ( id INT NOT NULL, conversation INT NOT NULL, text CHAR(1024), pather INT, first_response INT, second_response INT, PRIMARY KEY (id, conversation), FOREIGN KEY (conversation) REFERENCES conversation(id))",                     //, FOREIGN KEY (pather) REFERENCES post(id), FOREIGN KEY (first_response) REFERENCES post(id), FOREIGN KEY (second_response) REFERENCES post(id))",
		"CREATE TABLE author ( t_user INT NOT NULL, post INT NOT NULL, conversation INT NOT NULL, PRIMARY KEY (t_user,conversation, post), FOREIGN KEY (t_user) REFERENCES t_user(id), FOREIGN KEY (post, conversation) REFERENCES post(id, conversation))"} //, FOREIGN KEY (conversation) REFERENCES post(conversation))"}
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

		err := salvaUtente(user)
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
	return bytes.Equal(hashinput.Sum([]byte{}), user.password)
}

// variabili per gestire il database
var (
	insertUser   string    = "INSERT INTO t_user VALUES ($1, $2, $3)"
	insertUserOp *sql.Stmt = nil
)

func CreateDataBase() error {
	var err error

	db, err = sql.Open("postgres", "dbname=thinkzoneDB user=thinkzone")
	if err != nil {
		logs.Error("Impossibile creare il database")
		return err
	}

	for i := range createDbSqlScript {
		_, err = db.Exec(createDbSqlScript[i])
		if err != nil {
			logs.Error("Impossibile creare le tabelle del database")
			return err
		}

	}

	return nil

}

func salvaUtente(user *User) error {

	_, err := insertUserOp.Exec(user.ID, user.Username, user.password)
	if err != nil {
		logs.Error("Impossibile salvare ", user.Username, " nel database")
		return err
	}

	return nil
}

/*func SalvaUtenti() error {
	logs.Log("salvo1")
	operazione, err := db.Prepare(insertUser)
	if err != nil {
		logs.Error("Impossibile salvare gli utenti nel database")
		return err
	}
	defer operazione.Close()

	logs.Log("salvo2")
	for _, user := range Data.userIDtoUser {	
		_, err = operazione.Exec(user.ID, user.Username, user.password)
		if err != nil {
			logs.Error("Impossibile salvare ", user.Username, " nel database")
			return err
		}
	}

	logs.Log("salvo3")

	return nil
}*/

func (datab *DatabaseRegistration) CaricaUtenti() error {

	rows, err := db.Query("SELECT * FROM t_user")
	if err != nil {
		return err
	}

	logs.Log("caricamento utenti")
	for ; rows.Next(); datab.contatore++ {
		var username string
		var userID int
		var password string

		err := rows.Scan(&userID, &username, &password)
		if err != nil {
			return err
		}

		username = strings.TrimSpace(username)
		password = strings.TrimSpace(password)

		user := new(User)
		user.Username = username
		user.ID = userID
		user.password = []byte(password)

		fmt.Println("caricato vecchio utente: ", username, " ID: ", strconv.Itoa(userID), "\npassword: ", password)

		datab.userNameToId[username] = user
		datab.userIDtoUser[userID] = user

	}
	return nil
}

func init() {
	logs.Log("init del database")

	ServerFakeUser = new(User)
	ServerFakeUser.ID = 0
	ServerFakeUser.Username = "server"
	serverPassword := sha256.New()
	serverPassword.Write([]byte("toor"))
	ServerFakeUser.password = serverPassword.Sum([]byte{})

	MainConv = NewConversation(ServerFakeUser)

	err := CreateDataBase()
	if err != nil {
		//logs.Error(err.Error())
	}

	insertUserOp, err = db.Prepare(insertUser)
	if err != nil {
		logs.Error("Impossibile salvare gli utenti nel database\nmotivo: ", err.Error())
	}

	//err = salvaUtente(ServerFakeUser)
	if err != nil {
		//logs.Error("Impossibile salvare l'utente server nel database\nmotivo: ", err.Error())
	}

	Data.CaricaUtenti()
	if err != nil {
		logs.Error("Impossibile caricare gli utenti dal database\nmotivo: ", err.Error())
	}

	logs.AggiungiAzioneDiChiusura(func() {
		insertUserOp.Close()
	})
}
