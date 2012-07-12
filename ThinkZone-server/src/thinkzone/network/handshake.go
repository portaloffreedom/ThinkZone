// handshake
package network

import (
	"fmt"
	"net"
	"thinkzone/database"
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
