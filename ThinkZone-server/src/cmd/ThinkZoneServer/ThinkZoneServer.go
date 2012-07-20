// ThinkZone-server project ThinkZoneServer.go
package main

import (
	"fmt"
	"thinkzone/logs"
	"thinkzone/network"

	//	"database/sql"	
	//	_ "github.com/jbarham/gopgsqldriver"
	//	"math"
)

func main() {
	fmt.Println("Benvenuto su ThinkZone!")
	logs.StampaSuTerminale(true)

	network.StartServer(":4242")
	//network.StartServer(":80")

}
