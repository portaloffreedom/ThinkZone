package logs

import (
	"container/list"
	"fmt"
	"os"
	"os/signal"
)

var (
	azioniDiChiusura   *list.List // lista delle azioni da eseguire
	ChiusuraCompletata = make(chan int, 1)
)

func init() {
	azioniDiChiusura = list.New()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Kill, os.Interrupt)

	//	AggiungiAzioneDiChiusura(func() {
	//		fmt.Println("porco diooooooooooooooooooooooo")
	//	})

	go func(c chan os.Signal) {
		sig := <-c

		Log("segnale catturato dal modulo server: ", sig.String())
		Log("#### Chisura in corso ####")

		i := 0
		for fun := azioniDiChiusura.Front(); fun != nil; fun = fun.Next() {
			//svolgi tutte le azioni di chiusura
			i++
			fun.Value.(func())()
			fmt.Println("svoltaAzione", i)
		}
		//		close(c)

		fmt.Println("#### Operazioni di chiusura completate ####")
		ChiudiLog() //chiude il file di log
		ChiusuraCompletata <- 0
		return
	}(c)
}

func AggiungiAzioneDiChiusura(azione func()) {
	azioniDiChiusura.PushFront(azione)
}
