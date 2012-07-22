// errors
package logs

import (
	"fmt"
	"os"
	"strings"
)

// Scrivi il messaggio di errore sul file di log (il messaggio verrà stampato
// anche a terminale)
func Error(messageArray ...string) {
	message := strings.Join(messageArray, "")
	message = strings.Join([]string{"#ERROR#", message}, " ")

	fmt.Fprintln(os.Stderr, message)
	stampaSuFile(message)
}
