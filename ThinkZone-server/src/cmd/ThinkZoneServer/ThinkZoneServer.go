// ThinkZone-server project ThinkZoneServer.go
package main

import (
	"fmt"
	"thinkzone/logs"
	"thinkzone/network"
)

var Version string = "0.2a"

func main() {
	fmt.Println("Benvenuto su ThinkZone!")
	logs.StampaSuTerminale(true)

	/*logs.AggiungiAzioneDiChiusura(func() {
		logs.Log("inizio a salvare gli utenti")
		err := database.SalvaUtenti()
		if err != nil {
			logs.Error(err.Error())
			return
		}
		logs.Log("finito di salvare gli utenti")
	})*/

	network.StartServer(":4242")
	//network.StartServer(":80")

	<-logs.ChiusuraCompletata
}
