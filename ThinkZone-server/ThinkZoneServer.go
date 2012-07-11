// ThinkZone-server project main.go
package main

import (
	"fmt"
	"thinkzone/network"

	//	"database/sql"	
	//	_ "github.com/jbarham/gopgsqldriver"
	//	"math"
)

func main() {
	fmt.Println("Hello## World!")

	network.StartServer(":4242")
	//	StartServer(":80")

	//lingua := NewSuperString()
	//fmt.Println(lingua.GetCompleteWithSeparators("]["))
}
