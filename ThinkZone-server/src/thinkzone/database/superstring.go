// superstring
package database

import (
	"strconv"
	"strings"
	"thinkzone/logs"
)

// elemento della coda della superstring
type elemSuperString struct {
	elemento   []rune
	size       int
	succ, prec *elemSuperString
}

// Struttura dati della superstring
type SuperString struct {
	testa *elemSuperString
	dim   int
}

// crea un nuovo elemento della coda Superstring: non viene inserito elemento, vengono
// solo inizializzate variabili all'interno dell'elemento in creazione
func newElemSuperString(prec *elemSuperString, succ *elemSuperString) *elemSuperString {

	elem := new(elemSuperString)

	elem.elemento = make([]rune, 16, 16)
	elem.elemento = elem.elemento[0:0]
	elem.size = 0
	elem.succ = succ
	elem.prec = prec

	return elem
}

// Crea una nuova superstring vuota
func NewSuperString() *SuperString {

	// crea lista
	lista := new(SuperString)
	lista.testa = newElemSuperString(nil, nil)

	lista.dim = 1
	return lista
}

/* decostruttore serve?
SuperString::~SuperString()
{
    delete this->total;

    elemListaString *posAttuale;

    for (int i=1; true; i++) {
        posAttuale = testa->succ;
        delete testa;

        if (posAttuale == NULL)
            return;

        testa = posAttuale;
    }
}
*/

// Sostituisce la vecchia stringa memorizzata nel nodo con quella nuova
// (se ci fosse da gestire la memoria la gestirebbe)
//
// Attenzione! non ricalcola la nuova size per una questione di performance
func (elem *elemSuperString) sostituisciStringa(nuova []rune) {
	//	elem.elemento = nuova
	elem.elemento = make([]rune, len(nuova))
	copy(elem.elemento, nuova)
}

// Trasforma la supestring nella stringa vera e propria. Se separators == true
// allora ogni elemento della supestringa viene intervallato da parentesi
// quadre contenenti la dimensione di ogni singolo elemento.
//
// Esempio:  "ciao mondo!" ←→ "ciao [5]mondo![6]"
func (lista *SuperString) GetComplete(separators bool) string {

	tmp := lista.testa
	dim := lista.dim

	var elemento []string

	//	var elemento []string
	if separators == true {
		elemento = []string{"", "[", "", "]"}
	}

	var totalVector []string = make([]string, dim)

	tmp = lista.testa
	for j := 0; j < dim; j++ {
		if separators == true {
			elemento[0] = string(tmp.elemento)
			elemento[2] = strconv.Itoa(tmp.size)
			totalVector[j] = strings.Join(elemento, "")
		} else {
			totalVector[j] = string(tmp.elemento)
		}
		tmp = tmp.succ
	}

	return strings.Join(totalVector, "")

}

// Inserisce un singolo rune (carattere) all'interno della superstring
// nella data posizione
func (lista *SuperString) InsSingleElem(appendRune rune, pos int) {
	lista.InsElem([]rune{appendRune}, pos)
}

// Inserisce una stringa dentro la supestringa nella data posizione
func (lista *SuperString) InsStringElem(s string, pos int) {
	lista.InsElem([]rune(s), pos)
}

// Inserisce uno slice di rune dentro la supestringa nella data posizione
func (lista *SuperString) InsElem(appendRunes []rune, pos int) {

	if lista.dim == 1 && lista.testa.size == 0 {
		pos = -1
	}

	if lista.testa == nil {
		lista.testa = newElemSuperString(nil, nil)
	}

	dimStr := len(appendRunes)
	if dimStr == 0 {
		return
	}

	//roRebuildTotal = true

	//TODO eliminare questa parte
	if pos == 0 {
		lista.testa = newElemSuperString(nil, lista.testa)
		lista.testa.succ.prec = lista.testa
		lista.testa.elemento = appendRunes
		lista.testa.size = dimStr
		lista.dim += 1

		return
	}

	posAttuale := lista.testa
	//	cont := 0

	for posAttuale != nil && pos >= posAttuale.size {
		//		cont++
		pos -= posAttuale.size
		posAttuale = posAttuale.succ
	}

	if posAttuale == nil || pos == -1 {
		//fmt.Println("SuperString: inserisco in fondo all'ultima stringa")
		coda := lista.testa
		for coda.succ != nil {
			coda = coda.succ
		}

		//		coda.elemento = strings.Join([]string{coda.elemento, appendStr}, "")
		coda.elemento = append(coda.elemento, appendRunes...)
		coda.size += dimStr

		return
	}

	if pos == 0 {
		posAttuale = posAttuale.prec
	} else {
		//fmt.Println("SuperString: scindo due strighe")
		tmp := newElemSuperString(posAttuale, posAttuale.succ)
		if posAttuale.succ != nil {
			posAttuale.succ.prec = tmp
		}
		posAttuale.succ = tmp
		lista.dim += 1

		//split := strings.Fields(posAttuale.elemento)

		vecchio := posAttuale.elemento
		tmp.sostituisciStringa(vecchio[pos:posAttuale.size]) //TODO rotto
		posAttuale.sostituisciStringa(vecchio[0:pos])

		tmp.size = posAttuale.size - (pos)
		posAttuale.size = pos
	}

	//	posAttuale.elemento = strings.Join([]string{posAttuale.elemento, appendStr}, "")
	posAttuale.elemento = append(posAttuale.elemento, appendRunes...)
	posAttuale.size += dimStr

	return

}

//func (lista *SuperString) insElemChar(appendChar byte, pos int) {
//	lista.insElem(append("", appendChar), pos)
//}

// Rimuove dalla stringa corrente
func removeFromString(s string, s_size int, pos int, howmany int) string {

	return strings.Join([]string{s[0:pos], s[pos+howmany : s_size]}, "")

}

// Rimuove dallo slice di rune originale 
func removeFromRunes(origin []rune, s_size int, pos int, howmany int) []rune {
	//PER STAMPARE DEBUG
	//	var firstPart []rune
	//	if pos != 0 {  
	//		firstPart = origin[:pos]
	//	} else {
	//		firstPart = []rune{}
	//	}

	if pos+howmany < s_size {
		parteDaEliminare := origin[pos:]
		secondPart := origin[pos+howmany : s_size]
		//		fmt.Println("divisione stringa:", string(firstPart), ":", string(secondPart)) //DEBUG
		copy(parteDaEliminare, secondPart)
		//		origin = origin[:pos+1]
		//		origin = append(origin, origin[pos+howmany+1:s_size+1]...)
	}

	origin = origin[:s_size-howmany]
	return origin
}

// Elimina dalla Superstringa un tot caratteri nella posizione data.
// Elimina a partire dalla posizione data in poi
func (lista *SuperString) DelElem(pos int, howmany int) {
	//	fmt.Println("Cercando di eliminare: pos=", pos, " howmany=", howmany)
	if pos < 0 {
		logs.Error("che cazzo stai cercando di eliminare???")
		return
	}
	//roRebuildTotal = true

	posAttuale := lista.testa

	for pos > posAttuale.size {
		pos -= posAttuale.size
		posAttuale = posAttuale.succ
		if posAttuale == nil {
			logs.Error("Cercando di eliminare fuori dalla stringa, l'eliminazione verrà ignorata")
			return
		}
	}

	//elimina prima stringa
	quanti_eliminare := howmany
	if pos+howmany <= posAttuale.size { //elimina solo posizione attuale
		//		fmt.Println("Elimina solo posizione attuale") //DEBUG
		//vecchio := posAttuale.elemento
		//posAttuale.elemento = strings.Join([]string{vecchio[0:pos], vecchio[pos+howmany : posAttuale.size]}, "") //remove dalla stringa corrente
		if (posAttuale.size - quanti_eliminare) == 0 {
			lista.delSingleElem(posAttuale)
			return
		}
		posAttuale.elemento = removeFromRunes(posAttuale.elemento, posAttuale.size, pos, howmany)
		posAttuale.size -= quanti_eliminare
		return
	}

	quanti_eliminare = posAttuale.size - pos
	posAttuale.elemento = removeFromRunes(posAttuale.elemento, posAttuale.size, pos, quanti_eliminare)
	posAttuale.size -= quanti_eliminare

	if posAttuale.size == 0 {
		tmp := posAttuale.succ
		lista.delSingleElem(posAttuale)
		posAttuale = tmp
		if posAttuale == nil {
			posAttuale = lista.testa
		}
	} else {
		posAttuale = posAttuale.succ
	}

	//elimina le stringhe di mezzo
	quanti_eliminare = howmany - quanti_eliminare

	if posAttuale == nil {
		logs.Error("ma che cazzo succede? posAttuale == nil?")
		return
	}

	for posAttuale.succ != nil && quanti_eliminare >= posAttuale.size {
		//		fmt.Println("Elimino stringa di mezzo") //DEBUG
		quanti_eliminare -= posAttuale.size

		tmp := posAttuale.succ
		lista.delSingleElem(posAttuale)
		posAttuale = tmp
	}

	//elimina ultima stringa
	if posAttuale.size == quanti_eliminare {
		lista.delSingleElem(posAttuale)
	} else {
		//		fmt.Println("Elimino da ultima stringa") //DEBUG
		posAttuale.elemento = removeFromRunes(posAttuale.elemento, posAttuale.size, 0, quanti_eliminare)
		posAttuale.size -= quanti_eliminare
	}
}

// Elimina il nodo passato come parametro dalla Superstringa gestendo anche
// i casi limite in cui l'elemento non ha figli o parent e il caso in cui
// ci sia una sostituzione dell'elemento in testa
func (lista *SuperString) delSingleElem(elemento *elemSuperString) {
	if elemento.succ != nil {
		elemento.succ.prec = elemento.prec
	}

	if elemento.prec != nil {
		elemento.prec.succ = elemento.succ
	} else {
		lista.testa = elemento.succ
	}

	if lista.testa == nil {
		lista.testa = newElemSuperString(nil, nil)
	} else {
		lista.dim--
	}

	//delete elemento
	elemento = nil
}
