// errors
package logs

import (
	"fmt"
	"strings"
)

func Error(messageArray ...string) {
	message := strings.Join(messageArray, "")
	message = strings.Join([]string{"#ERROR#", message}, " ") //TODO stampare anche l'orario del log'

	fmt.Println(message)

	stampaSuFile(message)
}
