// logs
package logs

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

var (
	stampa_su_terminale bool   = false
	logFilename         string = "thinkzone-server"
	logFile             *os.File
	lastDate            string = ""
)

func getTimeString(t time.Time) (date, clock string) {
	hour, min, sec := t.Clock()
	year, month, day := t.Date()
	clock = strings.Join([]string{strconv.Itoa(hour), strconv.Itoa(min), strconv.Itoa(sec)}, ".")
	date = strings.Join([]string{strconv.Itoa(day), month.String(), strconv.Itoa(year)}, "-")
	return
}

func init() {
	fmt.Println("Inizializzo il file di log")

	now := time.Now()
	date, clock := getTimeString(now)
	logFilename = strings.Join([]string{logFilename, date, clock}, "_") + ".log"

	var err error
	logFile, err = os.Create(logFilename)
	if err != nil {
		fmt.Println("Errore: impossibile aprire il file di log\n\tmotivo:", err)
	}

	logFile.WriteString("File di log del server di ThinkZone\n")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func(c chan os.Signal) {
		sig := <-c
		//		if s == os.Kill {
		fmt.Println("\nsegnale catturato:", sig)
		fmt.Println("chiusura del file di log in corso")
		if logFile != nil {
			logFile.Close()
		} else {
			fmt.Println("impossibile chiudere il file di log")
		}
		//		}
		myself, err := os.FindProcess(os.Getpid())
		if err != nil {
			fmt.Println("ma che cazz? non riesco a trovare me stesso?\n\tmotivo:", err)
			return
		}
		myself.Signal(os.Kill)

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

	date, clock := getTimeString(time.Now())

	if lastDate != date {
		logFile.WriteString("\n" + date + "\n")
		lastDate = date
	}

	logFile.WriteString(clock + ": " + message)
	logFile.WriteString("\n")
}
