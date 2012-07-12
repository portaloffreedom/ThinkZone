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

func (client *Client) gestisciDisconnessione(conv *database.Conversation) {
	conv.UserDisconnection(client.user) //TODO conv prendila direttamente dal client
}

func gestisciClient(conn net.Conn) (*Client, func(chan *Client)) {
	fmt.Print("Nuova connessione: ")
	fmt.Println(conn.RemoteAddr())

	client := NewClient(&conn)
	if client == nil {
		conn.Close()
		return nil, nil
	}

	return client, func(readiness chan *Client) {

		for {
			_, err := client.stream.ReadByte()
			if err != nil {
				//TODO gestisci errore
				//TODO fare un pacchetto per la raccolta degli errori
				fmt.Print("connessione interrotta: ")
				fmt.Println(conn.RemoteAddr())
				client.gestisciDisconnessione(mainConv)
				return
			}
			err = client.stream.UnreadByte()
			if err != nil {
				//TODO gestisci errore
				//TODO fare un pacchetto per la raccolta degli errori
				fmt.Print("connessione interrotta: ")
				fmt.Println(conn.RemoteAddr())
				client.gestisciDisconnessione(mainConv)
				return
			}

			readiness <- client
			<-client.blocco
		}
	}
}

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

func (client *Client) handshake() bool {

	conn := client.conn

	//LETTURA USERNAME
	s, err := client.stream.ReadString('\\')
	if err != nil {
		logs.Error("Errore nel leggere lo username di: ", (*conn).RemoteAddr().String(), "\n motivazione: ", err.Error())
	}
	username := strings.Trim(s, "\\")
	logs.Log("IP:", (*conn).RemoteAddr().String(), " USERNAME:", username)

	//CONTROLLO SE L'UTENTE È GIÀ REGISTRATO
	var newuser bool
	client.user, newuser = database.Data.ConnectUser(username)
	if !newuser {
		logs.Log("connessione di un utente già registrato")
		//TODO gestisci se l'utente è già connesso alla conversazione - \
		//due utenti con lo stesso nome non possono essere connessi contemporaneamente

		//TODO RICHIESTA PASSWORD

	} else {
		//TODO registrazione
		logs.Log("TODO")
	}

	//REGISTRAZIONE DELL'UTENTE ALLA CONVERSAZIONE
	//TODO scegliere la conversazione a cui connettersi -.-
	err2 := mainConv.NewUserConnection(client.user)
	if err2 != nil {
		//CONTROLLO SE L'UTENTE È GIÀ CONNESSO
		logs.Error("impossibile connettere: ", err2.Error())
		return false
	}

	client.stream.WriteString(strconv.Itoa(client.user.ID))
	client.stream.WriteRune('\\')

	//TODO spedisci lo stato attuale della conversazione
	client.stream.WriteString(mainConv.TotalConversation())

	client.stream.Flush()
	return true
}
