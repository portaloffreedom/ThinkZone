// logs
package logs

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	// vero se i messaggi a livello di loggin devono essere stampati anche a terminale
	stampa_su_terminale bool = false

	// nome base del file di log del server
	logFilename string = "thinkzone-server"

	// stream del file di log 
	logFile *os.File

	// lock sul logFile
	logFileLock chan int = make(chan int, 1)

	// variabile necessaria per stabilire se scrivere che è cambiata data
	lastDate string = ""
)

// Trasforma la variabile time passata in ingresso in una stringa ordinata e leggibile
func getTimeString(t time.Time) (date, clock string) {
	hour, min, sec := t.Clock()
	year, month, day := t.Date()
	clock = strings.Join([]string{strconv.Itoa(hour), strconv.Itoa(min), strconv.Itoa(sec)}, ".")
	date = strings.Join([]string{strconv.Itoa(day), month.String(), strconv.Itoa(year)}, "-")
	return
}

// Init del pacchetto che apre il file di log in scrittura e gestisce cosa deve accadere in un 
// eventuale chiusura
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

	logFileLock <- 0 //abilita la prima scrittura sul log
}

// Questa funzione chiude il file di log. Da chiamare come procedura per chiudere il 
// server
func ChiudiLog() {
	fmt.Println("chiusura del file di log in corso")
	if logFile != nil {
		logFile.Close()
	} else {

		//TODO lock sulla scrittura sul log
		fmt.Println("impossibile chiudere il file di log")
	}
	return
}

// Imposta se i messaggi di log devono essere stampati anche a terminale
// (i messaggi di errore saranno stampati a terminale comunque)
func StampaSuTerminale(stampa bool) {
	<-logFileLock

	stampa_su_terminale = stampa

	logFileLock <- 0
}

// Scrivi il messaggio sul file di Log
func Log(messageArray ...string) {
	message := strings.Join(messageArray, "")

	stampaSuFile(message)
	if stampa_su_terminale {
		fmt.Println(message)
	}
}

// Scrivi sul file il messaggio di log
func stampaSuFile(message string) {
	<-logFileLock

	date, clock := getTimeString(time.Now())

	if lastDate != date {
		logFile.WriteString("\n" + date + "\n")
		lastDate = date
	}

	message = strings.Replace(message, "\n", "\n\t", -1)

	logFile.WriteString(clock + ": " + message)
	logFile.WriteString("\n")

	logFileLock <- 0
}
