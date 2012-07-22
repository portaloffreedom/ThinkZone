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

// Errore di traduzione dallo stream a un azione sensata
type TranslationError struct {
	s string
}

// Trasforma il TranslationError in una stringa
func (err TranslationError) Error() string {
	return err.s
}

// Struttura dati contente tutti i dati relativi ad un utente e allo
// stato della sua connessione
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

// Mangia un carattere dal canale e controlla che sia uguale al carattere c 
func mangiaCarattereDiControllo(c rune, input chan rune) bool {
	d := <-input
	if c == d {
		return true
	}
	//else
	return false
}

// Tenta di leggere una sequenza di caratteri dal canale trasformandoli
// in un intero. L'ultimo carattere letto che non è una cifra (0-9) 
// viene ritornato anch'esso per un controllo sulla coerenza della
// comunicazione
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

// Questa funzione entra in un ciclo e non vi esce finché il canale di input non
// viene chiuso. Quello che fa è un parsing dei caratteri dal canale di input
// controllando le azioni da svolgere e riversare (intelligentemente) sul canale
// di output quello che deve essere spedito ai vari client. 
//
// This function should run as goroutine
func gestisciTestoConversazione(input chan rune, output chan string) {
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
				//				for i := range versoClient {
				//					output <- versoClient[i]
				//				}
				output <- string(versoClient)
			}

		default:
			//database.MainConv.TestaPost.Text(activeUser).InsSingleElem(c, cursor)
			errore = database.MainConv.InsElem(activeUser, []rune{c})
			output <- string(c)
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

// Funzione che innesta tutta la procedura per il parsing di quello ricevuto
// dai client.
//TODO migliora la documentazione
func flasher(codaCiclica *list.List, readiness chan *Client) {
	tempoDaAspettare := 30 * time.Millisecond
	var lastActiveUser int = -1
	input := make(chan rune, 256)
	output := make(chan string, 256)
	logs.AggiungiAzioneDiChiusura(func() {
		close(input)
		close(output)
	})

	go gestisciTestoConversazione(input, output)

	// semplice funzione che rispedisce tutto quello che passa
	// dal canale di output sui socket delle connessioni
	go func() {
		for buffer := range output {
			for e := codaCiclica.Front(); e != nil; e = e.Next() {
				client := e.Value.(*Client)
				client.stream.WriteString(buffer)
				client.stream.Flush()
			}
		}
	}()

	var spegniti bool = false

	logs.AggiungiAzioneDiChiusura(func() {
		spegniti = true
	})

	for !spegniti {
		start := time.Now()

		quanti := len(readiness)

		for i := 0; i < quanti; i++ {
			clientAttivo := <-readiness
			var chiSonoString string

			if lastActiveUser != clientAttivo.user.ID {
				chiSonoString = strings.Join([]string{"\\U", strconv.Itoa(clientAttivo.user.ID), "\\"}, "")
				lastActiveUser = clientAttivo.user.ID
				database.MainConv.UtenteAttivo = lastActiveUser
			} else {
				chiSonoString = ""
			}
			chiSonoSSize := len(chiSonoString)

			//leggi cosa spedire
			var daLeggere int = clientAttivo.stream.Reader.Buffered()
			buffer := make([]rune, chiSonoSSize+daLeggere)
			var err error
			var size int
			for i, j := chiSonoSSize, 0; i < chiSonoSSize+daLeggere && j < daLeggere; i++ {
				buffer[i], size, err = clientAttivo.stream.ReadRune()
				j += size
				//fmt.Printf("#####size:_%v_ carattere:_%v_%v_\n",strconv.Itoa(size),string(buffer[i]),buffer[i])
				if err != nil {
					//TODO gestisci errore
					logs.Error("Errore nel leggere dalla rete")
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

// Funzione che fa partire la routine per la gestione del parsing dei "comandi"
// in arrivo dai client, svolge le azioni, rispedisce le risposte corrette.
//
// Inoltre spedisce i nuovi client nei client da gestire
func spedisci(codaNewConn chan *Client, readiness chan *Client) {
	codaCiclica := list.New()

	go flasher(codaCiclica, readiness)

	for nuovaConnessione := range codaNewConn {
		codaCiclica.PushFront(nuovaConnessione)
		//TODO sincronizzare questa parte visto che è eseguita contemporaneamte 
		// da più goroutine
	}
}

// Inizializzare il server in ascolto per le Sincronizzazione delle conversazioni
//
// "laddress string" indica su quali indirizza ascoltare e su quale porta 
func StartServer(laddress string) {
	ln, err := net.Listen("tcp", laddress)
	logs.Log("Server in ascolto su: \"", laddress, "\"")
	if err != nil {
		logs.Error("Errore nell'aprire la connessione: ", err.Error())
		return
		//TODO handle error
	}

	//database.MainConv = database.NewConversation(database.ServerFakeUser) //duplicato!

	//canale := make(chan byte, 256)
	codaReadiness := make(chan *Client, 64)
	codaAccettazioni := make(chan *Client, 64)
	logs.AggiungiAzioneDiChiusura(func() {
		close(codaReadiness)
		close(codaAccettazioni)
	})

	go spedisci(codaAccettazioni, codaReadiness)

	var spegniti bool = false

	logs.AggiungiAzioneDiChiusura(func() {
		spegniti = true
		ln.Close()
	})

	for !spegniti {
		conn, err := ln.Accept()
		if spegniti {
			return
		}
		if err != nil {
			//TODO fare un pacchetto per la raccolta degli errori
			logs.Error("Tentativo di connessione non andato a buon fine: ", err.Error())
		}

		go func() {
			client, gestore := GestisciClient(conn)
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
