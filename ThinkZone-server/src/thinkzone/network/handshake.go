// handshake
package network

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"thinkzone/database"
	"thinkzone/logs"
)

// Errore nel login
type ErrorLogin struct {
	message       string
	clientAddress net.Addr
}

// Crea l'errore del login
func NewErrorLogin(client *Client, message string) *ErrorLogin {
	err := new(ErrorLogin)
	err.message = message
	err.clientAddress = (*client.conn).RemoteAddr()
	return err
}

// Trasforma l'errore in stringa
func (err ErrorLogin) Error() string {
	return "Login Error from: " + err.clientAddress.String() + "\nMessaggio: " + err.message
}

// gestisce la disconnessione del client
func (client *Client) gestisciDisconnessione(conv *database.Conversation) {
	conv.UserDisconnection(client.user) //TODO conv prendila direttamente dal client
}

// Inizializza la struttura dati del client e crea la funzione per gestire il client
// (da partire come goroutine)
func GestisciClient(conn net.Conn) (*Client, func(chan *Client)) {
	fmt.Print("Nuova connessione: ")
	fmt.Println(conn.RemoteAddr())

	client := NewClient(&conn)
	if client == nil {
		conn.Close()
		return nil, nil
	}

	return client, func(readiness chan *Client) {

		for {
			_, _, err := client.stream.ReadRune()
			if err != nil {
				logs.Error("connessione interrotta: ", conn.RemoteAddr().String(), "\n\tmotivazione: ", err.Error())
				client.gestisciDisconnessione(database.MainConv)
				return
			}
			err = client.stream.UnreadRune()
			if err != nil {
				logs.Error("impossibile fare UnreadByte: ", conn.RemoteAddr().String(), "\n\tmotivazione: ", err.Error())
				client.gestisciDisconnessione(database.MainConv)
				return
			}

			readiness <- client
			<-client.blocco
		}
	}
}

// Inizializza la struttura dai del client e ne gestisce l'handshake prima di ritornare
//
// Ritorna nil se l'handshake non ha avuto successo
//TODO ritornare anche un errore 
func NewClient(conn *net.Conn) *Client {
	var client *Client = new(Client)
	client.conn = conn
	client.stream = bufio.NewReadWriter(bufio.NewReader(*conn), bufio.NewWriter(*conn))

	success := client.handshake()
	if !success {
		return nil
	}

	client.blocco = make(chan int, 1)
	return client
}

// Svolge l'handshake con il client. Ritorna true se l'handshake ha avuto successo
func (client *Client) handshake() bool {

	keepAlive, err := client.leggiEseguiComando()
	if err != nil {
		logs.Error(err.Error())
		return false
	}

	if !keepAlive {
		return false
	}

	//REGISTRAZIONE DELL'UTENTE ALLA CONVERSAZIONE
	//TODO scegliere la conversazione a cui connettersi -.-
	err2 := database.MainConv.NewUserConnection(client.user)
	if err2 != nil {
		//CONTROLLO SE L'UTENTE È GIÀ CONNESSO
		logs.Error("impossibile connettere: ", err2.Error())
		return false
	}

	client.stream.WriteString(strconv.Itoa(client.user.ID))
	client.stream.WriteRune('\\')

	//spedisci lo stato attuale della conversazione
	client.stream.WriteString("\\U0\\")
	client.stream.WriteString(database.MainConv.GetComplete(false))

	//Invia le posizioni di tutte quante le persone connesse
	client.stream.WriteString(database.MainConv.GetAllPositionString())
	//invia l'utente attivo
	client.stream.WriteString("\\U" + strconv.Itoa(database.MainConv.UtenteAttivo) + "\\")

	client.stream.Flush()
	return true
}

// Da chiamare subito successivamente all'handshake (propriamente viene chiamato
// all'interno della handshake() stessa) e svolge l'azione richiesta dal client.
//
// Azioni disponibili: Registrazione e Login
//
// Ritorna true se la connessione deve rimanere aperta
//TODO ritorna invece l'azione da eseguire sul client come funzione
func (client *Client) leggiEseguiComando() (bool, error) {
	s, err := client.stream.ReadString('\\')
	if err != nil {
		return false, NewErrorLogin(client, "Errore nel leggere\n\tMotivazione: "+err.Error())
	}

	if s != "\\" {
		return false, NewErrorLogin(client, "Errore di comunicazione, mi aspettavo un comando, ho invece ricevuto: "+s)
	}

	command, err := client.stream.ReadString('\\')
	if err != nil {
		return false, NewErrorLogin(client, "Impossibile leggere comando \n\tmotivazione: "+err.Error())
	}

	switch command {
	case "L0\\":
		err := client.login()
		return true, err
	case "L1\\":
		err := client.registrati()
		return false, err
	default:
		return false, NewErrorLogin(client, "Comando non supportato dal server. Invocato da "+(*client.conn).RemoteAddr().String())
	}

	return false, NewErrorLogin(client, "boh...")
}

// Registra l'utente nel server come nuovo utente
func (client *Client) registrati() error {
	logs.Log("registrazione nuovo utente")
	username, err := client.stream.ReadString('\\')
	if err != nil {
		return NewErrorLogin(client, "Errore nel leggere lo username\n\tmotivo: "+err.Error())
	}
	username = strings.Trim(username, "\\")

	password, err := client.stream.ReadString('\\')
	if err != nil {
		return NewErrorLogin(client, "Errore nel leggere la password\n\tmotivo: "+err.Error())
	}
	password = strings.Trim(password, "\\")

	var success bool
	client.user, success = database.Data.RegisterNewUser(username, password)
	if !success {
		return NewErrorLogin(client, "impossibile registrare nuovo utente")
	}

	logs.Log("Registrato nuovo utente: ", username, " password: ", password)
	client.stream.WriteString("\\R0\\")
	client.stream.Flush()

	return nil
}

// Esegue il login del client (registra lo stato online
func (client *Client) login() error {
	logs.Log("login nuovo utente")
	username, err := client.stream.ReadString('\\')
	if err != nil {
		return NewErrorLogin(client, "Errore nel leggere lo username\n\tmotivo: "+err.Error())
	}
	username = strings.Trim(username, "\\")

	//CONTROLLO SE L'UTENTE È GIÀ REGISTRATO
	var newuser bool
	client.user, newuser = database.Data.ConnectUser(username)
	if newuser {
		err := NewErrorLogin(client, "connessione impossibile: Utente non registrato!")
		client.stream.WriteString("\\R1\\")
		client.stream.Flush()
		return err
	}

	logs.Log("IP:", (*client.conn).RemoteAddr().String(), " USERNAME:", username, " UserID:", strconv.Itoa(client.user.ID))

	//RICHIESTA PASSWORD
	password, err := client.stream.ReadString('\\')
	if err != nil {
		return NewErrorLogin(client, "Errore nel leggere la password\n\tmotivo: "+err.Error())
	}

	password = strings.Trim(password, "\\")
	if !client.user.VerifyPassword(password) {
		return NewErrorLogin(client, username+" Password errata")
	}

	logs.Log("Login riuscito di ", username, " da ", (*client.conn).RemoteAddr().String())

	client.stream.WriteString("\\R0\\")
	client.stream.Flush()
	return nil
}
