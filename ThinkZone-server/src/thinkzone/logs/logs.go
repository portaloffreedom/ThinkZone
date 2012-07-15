// logs
package logs

import (
	"fmt"
	"os"
	//	"os/signal"
	"strings"
	//	"strconv"
)

var (
	stampa_su_terminale bool   = false
	logFilename         string = "thinkzone-server.log"
	logFile             *os.File
)

func init() {
	fmt.Println("Inizializzo il file di log")
	var err error
	logFile, err = os.Create(logFilename)
	if err != nil {
		fmt.Println("Errore: impossibile aprire il file di log\n\tmotivo:", err)
	}

	logFile.WriteString("File di log del server di ThinkZone\n")

	c := make(chan os.Signal, 1)
	//	signal.Notify(c)	

	go func(c chan os.Signal) {
		s := <-c
		fmt.Println("catturato un segnale")
		//		if s == os.Kill {
		fmt.Println("segnale catturato:", s)
		fmt.Println("chiusura del file di log in corso")
		logFile.Close()
		//		}
		return
	}(c)
}

//TODO var tutti i log
//type LogFile interface {
//}

//TODO lock sulla scrittura sul log

func StampaSuTerminale(stampa bool) {
	stampa_su_terminale = stampa
}

func Log(messageArray ...string) {
	message := strings.Join(messageArray, "")
	//	message = strings.Join([]string{message}, " : ") //TODO stampare anche l'orario del log'

	if stampa_su_terminale {
		fmt.Println(message)
	}

	stampaSuFile(message)
}

func stampaSuFile(message string) {
	logFile.WriteString(message)
	logFile.WriteString("\n")
}
