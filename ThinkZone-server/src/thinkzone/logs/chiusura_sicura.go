package logs

import (
	"container/list"
	"os"
	"os/signal"
)

var (
	azioniDiChiusura *list.List // lista delle azioni da eseguire
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

		for fun := azioniDiChiusura.Front(); fun != nil; fun = fun.Next() {
			//svolgi tutte le azioni di chiusura
			fun.Value.(func())()
		}
		close(c)
		
		Log("#### Operazioni di chiusura completate ####")
		ChiudiLog() //chiude il file di log
		return
	}(c)
}

func AggiungiAzioneDiChiusura(azione func()) {
	azioniDiChiusura.PushFront(azione)
}
