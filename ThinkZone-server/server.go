// server
package main

import (
	//	"ThinkZoneDatabase"
	"bufio"
	"container/list"
	"fmt"
	"net"
	"strconv" //to convert integer into a string
	"strings"
	"time"

//	"unsafe"
)

type Client struct {
	conn   *net.Conn
	stream *bufio.ReadWriter
	blocco chan int
	user   *User
	//	username string //duplicated value
	//  userid   int //duplicated value
}

var serverFakeUser User = User{42, "server"}
var mainConv *Conversation /*{	"prova",
1,
connected     map[int]*User
postMap       map[int]*Post
contatorePost int
testaPost     *Post	 }*/

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
	var newuser bool
	client.user, newuser = data.ConnectUser(username)
	mainConv.NewUserConnection(client.user)
	if !newuser {
		fmt.Println("impossibile connettere di nuovo lo stesso userid")
		return nil
	}

	client.stream.WriteString(strconv.Itoa(client.user.id))
	client.stream.WriteRune('\\')
	client.stream.Flush()

	client.blocco = make(chan int, 1)
	return client
}

func gestisciClient(conn net.Conn) (*Client, func(chan *Client)) {
	fmt.Print("Nuova connessione: ")
	fmt.Println(conn.RemoteAddr())

	//client := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	client := NewClient(&conn)
	if client == nil {
		conn.Close()
		return nil, nil
	}

	return client, func(readiness chan *Client) {

		for {
			//buf, err := client.stream.ReadByte()
			_, err := client.stream.ReadByte()
			if err != nil {
				//TODO gestisci errore
				//TODO fare un pacchetto per la raccolta degli errori
				fmt.Print("connessione interrotta: ")
				fmt.Println(conn.RemoteAddr())
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
func mangiaCarattereDiControllo(c byte, input chan byte) bool {
	d := <-input
	if c == d {
		return true
	}
	//else
	return false
}

func mangiaIntero(input chan byte) (valore int, lastRead byte) {
	buffer := make([]byte, 32, 32)
	for i := 0; i < 32; i++ {
		//		buffer = buffer[i+1]
		b := <-input
		fmt.Print(string(b))
		buffer[i] = b
		switch b {
		case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
		default:
			lastRead = buffer[i]
			buffer = buffer[:i]

			valore64, err := strconv.ParseInt(string(buffer), 10, 8)
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
func gestisciTestoConversazione(input chan byte) {
	activeUser := data.GetUserByID(0)
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
				mainConv.testaPost.Text(activeUser).delElem(cursor, howmany)
			case 'U':
				newUserID, ccc := mangiaIntero(input)
				if ccc != '\\' {
					fmt.Println("ERRORE lettura stream: carattere di controllo mangiato non Intero")
				}
				activeUser = data.GetUserByID(newUserID)

			case '\\':
				mainConv.testaPost.Text(activeUser).insElem("\\", cursor)

			default:
				fmt.Println("ERRORE azione", cc, "non disponibile")
			}

		default:
			mainConv.testaPost.Text(activeUser).insElem(string(c), cursor)
			cursor++

			fmt.Println("---", mainConv.testaPost.Text(activeUser).GetComplete(true), "---") //DEBUG

			//TODO anche qui si può pensare ad un flasher a tempo per minimizzare il lavoro su superstring
		}
	}
}

func flasher(codaCiclica *list.List, readiness chan *Client) {
	tempoDaAspettare := 20 * time.Millisecond
	var lastActiveUser int = -1
	toSuperString := make(chan byte, 256)

	go gestisciTestoConversazione(toSuperString)

	for {
		start := time.Now()

		quanti := len(readiness)

		for i := 0; i < quanti; i++ {
			clientAttivo := <-readiness
			var chiSonoString string

			if lastActiveUser != clientAttivo.user.id {
				chiSonoString = strings.Join([]string{"\\U", strconv.Itoa(clientAttivo.user.id), "\\"}, "")
				lastActiveUser = clientAttivo.user.id
			} else {
				chiSonoString = ""
			}
			chiSonoSSize := len(chiSonoString)
			//fmt.Printf("dimensione di %s: %d\n", chiSonoString, chiSonoSSize)

			//leggi cosa spedire
			var daLeggere int = clientAttivo.stream.Reader.Buffered()
			buffer := make([]byte, chiSonoSSize+daLeggere)
			var err error
			for i := chiSonoSSize; i < chiSonoSSize+daLeggere; i++ {
				buffer[i], err = clientAttivo.stream.ReadByte()
				toSuperString <- buffer[i]
				if err != nil {
					//TODO gestisci errore
					fmt.Println("Errore nel leggere dalla rete")
					clientAttivo.blocco <- 1 //TODO dovresti chiudere il canale e tutto quanto
					break
				}
			}
			clientAttivo.blocco <- 0

			//prepara il buffer da spedire
			for i := 0; i < chiSonoSSize; i++ {
				buffer[i] = []byte(chiSonoString)[i]
			}

			//spedisci
			for e := codaCiclica.Front(); e != nil; e = e.Next() {
				client := e.Value.(*Client)
				client.stream.Write(buffer)
				client.stream.Flush()
			}
		}

		duration := time.Since(start)
		if duration <= tempoDaAspettare {
			//			fmt.Printf("ho aspettato: %v ", tempoDaAspettare-duration)
			//			fmt.Printf("### dati in coda: %d ###\n", len(readiness))
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
	if err != nil {
		fmt.Println("Errore nell'aprire la connessione")
		//TODO handle error
	}

	mainConv = NewConversation(&serverFakeUser)

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
		go gestore(codaReadiness)
		codaAccettazioni <- client
	}
}
