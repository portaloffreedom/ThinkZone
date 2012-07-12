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
	"time"

//	"unsafe"
)

type Client struct {
	//socket TCP
	conn *net.Conn

	//stream del socket TCP
	stream *bufio.ReadWriter

	//se il blocco è attivo il server sta leggendo lo stream di questo Client
	blocco chan int

	//utente associato al client
	user *database.User
}

var serverFakeUser database.User = database.User{42, "server"}
var mainConv *database.Conversation = database.NewConversation(&serverFakeUser) /*{	"prova",
 * 1,
 * connected     map[int]*database.User
 * postMap       map[int]*Post
 * contatorePost int
 * TestaPost     *Post	 }
 */

func NewClient(conn *net.Conn) *Client {
	var client *Client = new(Client)
	client.stream = bufio.NewReadWriter(bufio.NewReader(*conn), bufio.NewWriter(*conn))

	//TODO get username from client TCP stream
	s, err := client.stream.ReadString('\\')
	if err != nil {
		fmt.Print("ERRORE NEL LEGGERE LO USERNAME DI: ")
		fmt.Println((*conn).RemoteAddr())
	}
	username := strings.Trim(s, "\\")
	fmt.Println("IP:", (*conn).RemoteAddr(), "USERNAME:", username)
	var newuser bool
	client.user, newuser = database.Data.ConnectUser(username)
	if !newuser {
		fmt.Println("connessione di un utente già registrato")
		//TODO gestisci se l'utente è già connesso alla conversazione - \
		//due utenti con lo stesso nome non possono essere connessi contemporaneamente
	}

	err2 := mainConv.NewUserConnection(client.user)
	if err2 != nil {
		fmt.Println("#01", err2)
		return nil
	}

	client.stream.WriteString(strconv.Itoa(client.user.ID))
	client.stream.WriteRune('\\')
	client.stream.Flush()

	client.blocco = make(chan int, 1)
	return client
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
func gestisciTestoConversazione(input chan rune) {
	activeUser := database.Data.GetUserByID(0)
	cursor := 0
	for {
		c := <-input
		switch c {
		case '\\': //caso di carattere di controllo
			cc := <-input
			switch cc {
			case 'P':
				cursor, cc = mangiaIntero(input)
				if cc != '\\' {
					fmt.Println("ERRORE lettura stream: carattere di controllo mangiato non Intero")
				}
			case 'D':
				howmany, ccc := mangiaIntero(input)
				if ccc != '\\' {
					fmt.Println("ERRORE lettura stream: carattere di controllo mangiato non Intero")
				}
				cursor -= howmany
				//TODO rimuovi testo
				mainConv.TestaPost.Text(activeUser).DelElem(cursor, howmany)
			case 'U':
				newUserID, ccc := mangiaIntero(input)
				if ccc != '\\' {
					fmt.Println("ERRORE lettura stream: carattere di controllo mangiato non Intero")
				}
				activeUser = database.Data.GetUserByID(newUserID)

			case '\\':
				mainConv.TestaPost.Text(activeUser).InsSingleElem('\\', cursor)
				cursor++

			default:
				fmt.Println("ERRORE azione", cc, "non disponibile")
			}

		default:
			mainConv.TestaPost.Text(activeUser).InsSingleElem(c, cursor)
			cursor++

			//TODO anche qui si può pensare ad un flasher a tempo per minimizzare il lavoro su superstring
		}

		//		fmt.Println("---", mainConv.TestaPost.Text(activeUser).GetComplete(true), "---") //DEBUG
		fmt.Println(mainConv.TestaPost.Text(activeUser).GetComplete(true)) //DEBUG

	}
}

func flasher(codaCiclica *list.List, readiness chan *Client) {
	tempoDaAspettare := 20 * time.Millisecond
	var lastActiveUser int = -1
	toSuperString := make(chan rune, 256)

	go gestisciTestoConversazione(toSuperString)

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
				toSuperString <- buffer[i]
				if err != nil {
					//TODO gestisci errore
					fmt.Println("Errore nel leggere dalla rete")
					clientAttivo.blocco <- 1 //TODO dovresti chiudere il canale e tutto quanto
					//client.
					break
				}
			}
			clientAttivo.blocco <- 0

			//prepara il buffer da spedire
			for i := 0; i < chiSonoSSize; i++ {
				buffer[i] = []rune(chiSonoString)[i]
			}

			//spedisci
			for e := codaCiclica.Front(); e != nil; e = e.Next() {
				client := e.Value.(*Client)
				client.stream.WriteString(string(buffer))
				client.stream.Flush()
			}
		}

		duration := time.Since(start)
		if duration <= tempoDaAspettare {
			time.Sleep(tempoDaAspettare - duration)
		}
	}
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
				return
			}

			readiness <- client
			<-client.blocco
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
	if err != nil {
		fmt.Println("Errore nell'aprire la connessione")
		//TODO handle error
	}

	mainConv = database.NewConversation(&serverFakeUser)

	//canale := make(chan byte, 256)
	codaReadiness := make(chan *Client, 64)
	codaAccettazioni := make(chan *Client, 64)

	go spedisci(codaAccettazioni, codaReadiness)

	for {
		conn, err := ln.Accept()
		if err != nil {
			//TODO fare un pacchetto per la raccolta degli errori
			fmt.Println("Tentativo di connessione non andato a buon fine")
		}

		client, gestore := gestisciClient(conn)
		if client != nil {
			go gestore(codaReadiness)
			codaAccettazioni <- client
		} else {
			conn.Close()
		}
	}
}
