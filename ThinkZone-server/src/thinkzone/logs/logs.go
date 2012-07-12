// logs
package logs

import (
	"fmt"
	//	"strconv"
	"strings"
)

//TODO var tutti i log
//type LogFile interface {
//}

var (
	stampa_su_terminale bool   = false
	logFilename         string = "thinkzone-server.log"
)

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

}
