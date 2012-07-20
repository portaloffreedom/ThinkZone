package network

import (
	"container/list"
	"fmt"
	"os"
	"os/signal"
	"thinkzone/logs"
)

var (
	azioniDiChiusura *list.List // lista delle azioni da eseguire
)

func init() {
	azioniDiChiusura = list.New()
	c := make(chan os.Signal, 1)
	signal.Notify(c)

	AggiungiAzioneDiChiusura(func() {
		fmt.Println("porco diooooooooooooooooooooooo")
	})

	go func(c chan os.Signal) {
		sig := <-c

		logs.Log("\nsegnale catturato dal modulo server: ", sig.String())

		for fun := azioniDiChiusura.Front(); fun != nil; fun = fun.Next() {
			//svolgi tutte le azioni di chiusura
			fun.Value.(func())()
		}
		close(c)
		return
	}(c)
}

func AggiungiAzioneDiChiusura(azione func()) {
	azioniDiChiusura.PushFront(azione)
}
