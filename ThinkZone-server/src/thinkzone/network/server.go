// server
package network

import (
	"bufio"
	"container/list"
	"fmt"
	"net"
	"strconv" //to convert integer into a string
	"strings"
	"thinkzone/database"
	"thinkzone/logs"
	"time"

//	"unsafe"
)

type TranslationError struct {
	s string
}

func (err TranslationError) Error() string {
	return err.s
}

type Client struct {
	//socket TCP
	conn *net.Conn

	//stream del socket TCP
	stream *bufio.ReadWriter

	//se il blocco è attivo il server sta leggendo lo stream di questo Client
	blocco chan int

	//utente associato al client
	user *database.User

	//conversazione attiva
	//TODO *database.Conversation
}

func mangiaCarattereDiControllo(c rune, input chan rune) bool {
	d := <-input
	if c == d {
		return true
	}
	//else
	return false
}

func mangiaIntero(input chan rune) (valore int, lastRead rune) {
	buffer := make([]rune, 256, 256)
	for i := 0; i < 32; i++ {
		//		buffer = buffer[i+1]
		b := <-input
		//		fmt.Print(string(b)) //DEBUG
		buffer[i] = b
		switch b {
		case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
		default:
			lastRead = buffer[i]
			buffer = buffer[:i]

			valore64, err := strconv.ParseInt(string(buffer), 10, 0)
			if err != nil {
				fmt.Println("ERRORE nel convertire string in int")
			}
			valore = int(valore64)
			return
		}
	}

	return -1, 0
}

//this function should run as goroutine
func gestisciTestoConversazione(input chan rune, output chan rune) {
	activeUser := database.Data.GetUserByID(0)
	var errore error
	var versoClient []rune
	//	cursor := 0
	for c := range input {
		//		c := <-input
		//		fmt.Println("input",string(c))
		switch c {
		case '\\': //caso di carattere di controllo
			cc := <-input
			//			fmt.Println("input",string(cc))
			switch cc {
			case 'K':
				parent, cu := mangiaIntero(input)
				if cu != '\\' {
					errore = TranslationError{"ERRORE lettura stream: carattere di controllo mangiato != '\\'"}
					break
				}
				var postRID int
				postRID, errore = database.MainConv.Respond(parent, activeUser)
				versoClient = []rune{'\\', 'K'}
				versoClient = append(versoClient, ([]rune(strconv.Itoa(parent)))...)
				versoClient = append(versoClient, '\\')
				versoClient = append(versoClient, ([]rune(strconv.Itoa(postRID)))...)
				versoClient = append(versoClient, '\\')
			case 'P':
				postCursor, cu := mangiaIntero(input)
				if cu != '\\' {
					errore = TranslationError{"ERRORE lettura stream: carattere di controllo mangiato != '\\'"}
					break
				}
				errore = database.MainConv.ChangePost(activeUser, postCursor)
				versoClient = []rune{'\\', cc}
				versoClient = append(versoClient, ([]rune(strconv.Itoa(postCursor)))...)
				versoClient = append(versoClient, '\\')
			case 'C':
				cursor, cu := mangiaIntero(input)
				if cu != '\\' {
					errore = TranslationError{"ERRORE lettura stream: carattere di controllo mangiato != '\\'"}
					break
				}
				errore = database.MainConv.ChangePos(activeUser, cursor)
				versoClient = []rune{'\\', cc}
				versoClient = append(versoClient, ([]rune(strconv.Itoa(cursor)))...)
				versoClient = append(versoClient, '\\')
			case 'D':
				howmany, cu := mangiaIntero(input)
				if cu != '\\' {
					errore = TranslationError{"ERRORE lettura stream: carattere di controllo mangiato != '\\'"}
					break
				}
				versoClient = []rune{'\\', cc}
				versoClient = append(versoClient, ([]rune(strconv.Itoa(howmany)))...)
				versoClient = append(versoClient, '\\')
				//cursor -= howmany
				//database.MainConv.TestaPost.Text(activeUser).DelElem(cursor, howmany)
				errore = database.MainConv.DelElem(activeUser, howmany)
			case 'U':
				newUserID, cu := mangiaIntero(input)
				if cu != '\\' {
					errore = TranslationError{"ERRORE lettura stream: carattere di controllo mangiato != '\\'"}
					break
				}
				activeUser = database.Data.GetUserByID(newUserID)
				versoClient = []rune{'\\', cc}
				versoClient = append(versoClient, ([]rune(strconv.Itoa(newUserID)))...)
				versoClient = append(versoClient, '\\')

			case '\\':
				//database.MainConv.TestaPost.Text(activeUser).InsSingleElem('\\', cursor)
				errore = database.MainConv.InsElem(activeUser, []rune{'\\'})
				versoClient = []rune{'\\'}
				//cursor++

			default:
				fmt.Println("ERRORE azione", cc, "non disponibile")
			}

			if errore != nil {
				logs.Error("errore nel lavorare sulla superstringa\n\tultimo comando: ", string(cc), "\n\tmotivo: ", errore.Error())
				errore = nil
			} else {
				for i := range versoClient {
					output <- versoClient[i]
				}
			}

		default:
			//database.MainConv.TestaPost.Text(activeUser).InsSingleElem(c, cursor)
			errore = database.MainConv.InsElem(activeUser, []rune{c})
			output <- c
			//cursor++

			//TODO anche qui si può pensare ad un flasher a tempo per minimizzare il lavoro su superstring
		}

		if errore != nil {
			logs.Error("errore nel lavorare sulla superstringa\n\t ultimo carattere: ", string(c), "\n\tmotivo: ", errore.Error())
			errore = nil
		}

		//		fmt.Println("---", database.MainConv.TestaPost.Text(activeUser).GetComplete(true), "---") //DEBUG
		fmt.Println(database.MainConv.GetComplete(true)) //DEBUG

	}
}

func flasher(codaCiclica *list.List, readiness chan *Client) {
	tempoDaAspettare := 30 * time.Millisecond
	var lastActiveUser int = -1
	input := make(chan rune, 256)
	output := make(chan rune, 256)

	go gestisciTestoConversazione(input, output)

	go func() {
		for buffer := range output {
			for e := codaCiclica.Front(); e != nil; e = e.Next() {
				client := e.Value.(*Client)
				client.stream.WriteString(string(buffer))
				client.stream.Flush()
			}
		}
	}()

	for {
		start := time.Now()

		quanti := len(readiness)

		for i := 0; i < quanti; i++ {
			clientAttivo := <-readiness
			var chiSonoString string

			if lastActiveUser != clientAttivo.user.ID {
				chiSonoString = strings.Join([]string{"\\U", strconv.Itoa(clientAttivo.user.ID), "\\"}, "")
				lastActiveUser = clientAttivo.user.ID
			} else {
				chiSonoString = ""
			}
			chiSonoSSize := len(chiSonoString)
			//fmt.Printf("dimensione di %s: %d\n", chiSonoString, chiSonoSSize)

			//leggi cosa spedire
			var daLeggere int = clientAttivo.stream.Reader.Buffered()
			buffer := make([]rune, chiSonoSSize+daLeggere)
			var err error
			for i := chiSonoSSize; i < chiSonoSSize+daLeggere; i++ {
				buffer[i], _, err = clientAttivo.stream.ReadRune()
				//				toSuperString <- buffer[i]
				if err != nil {
					//TODO gestisci errore
					fmt.Println("Errore nel leggere dalla rete")
					clientAttivo.blocco <- 1 //TODO dovresti chiudere il canale e tutto quanto
					clientAttivo.gestisciDisconnessione(database.MainConv)
					break
				}
			}
			clientAttivo.blocco <- 0

			//prepara il buffer da spedire
			for i := 0; i < chiSonoSSize; i++ {
				buffer[i] = []rune(chiSonoString)[i]
			}
			//spedisci
			for i := 0; i < len(buffer); i++ {
				input <- buffer[i]
			}
		}

		duration := time.Since(start)
		if duration <= tempoDaAspettare {
			time.Sleep(tempoDaAspettare - duration)
		}
	}
}

func spedisci(codaNewConn chan *Client, readiness chan *Client) {
	codaCiclica := list.New()

	go flasher(codaCiclica, readiness)

	for nuovaConnessione := range codaNewConn {
		codaCiclica.PushFront(nuovaConnessione)
	}
}

func StartServer(laddr string) {
	ln, err := net.Listen("tcp", laddr)
	logs.Log("Server in ascolto su: \"", laddr, "\"")
	if err != nil {
		logs.Error("Errore nell'aprire la connessione: ", err.Error())
		return
		//TODO handle error
	}

	database.MainConv = database.NewConversation(&database.ServerFakeUser)

	//canale := make(chan byte, 256)
	codaReadiness := make(chan *Client, 64)
	codaAccettazioni := make(chan *Client, 64)

	go spedisci(codaAccettazioni, codaReadiness)

	for {
		conn, err := ln.Accept()
		if err != nil {
			//TODO fare un pacchetto per la raccolta degli errori
			logs.Error("Tentativo di connessione non andato a buon fine: ", err.Error())
		}

		go func() {
			client, gestore := gestisciClient(conn)
			if client != nil {
				go gestore(codaReadiness)
				codaAccettazioni <- client
			} else {
				//L'handshaking non è andato a buon fine
				conn.Close()
			}
		}()

	}
}
