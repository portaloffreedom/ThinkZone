// superstring
package main

import (
	"fmt"
	"strconv"
	"strings"
)

type elemSuperString struct {
	elemento   []rune
	size       int
	succ, prec *elemSuperString
}

type SuperString struct {
	testa *elemSuperString
	dim   int
}

func NewElemSuperString(prec *elemSuperString, succ *elemSuperString) *elemSuperString {

	elem := new(elemSuperString)

	elem.elemento = []rune{' '}
	elem.size = 0
	elem.succ = succ
	elem.prec = prec

	return elem
}

func NewSuperString() *SuperString {

	// crea lista
	lista := new(SuperString)
	lista.testa = NewElemSuperString(nil, nil)

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

// Attenzione! non ricalcola la nuova size per una questione di performance
// @param nuova
func (elem *elemSuperString) sostituisciStringa(nuova []rune) {
	elem.elemento = nuova
}

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

func (lista *SuperString) insSingleElem(appendRune rune, pos int) {
	lista.insElem([]rune{appendRune}, pos)
}

func (lista *SuperString) insStringElem(s string, pos int) {
	lista.insElem([]rune(s), pos)
}

func (lista *SuperString) insElem(appendRunes []rune, pos int) {

	if lista.dim == 1 && lista.testa.size == 0 {
		pos = -1
	}

	if lista.testa == nil {
		lista.testa = NewElemSuperString(nil, nil)
	}

	dimStr := len(appendRunes)
	if dimStr == 0 {
		return
	}

	//roRebuildTotal = true

	//TODO eliminare questa parte
	if pos == 0 {
		lista.testa = NewElemSuperString(nil, lista.testa)
		lista.testa.succ.prec = lista.testa
		lista.testa.elemento = appendRunes
		lista.testa.size = dimStr
		lista.dim += 1

		return
	}

	posAttuale := lista.testa
	cont := 0

	for posAttuale != nil && pos >= posAttuale.size {
		cont++
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
		tmp := NewElemSuperString(posAttuale, posAttuale.succ)
		if posAttuale.succ != nil {
			posAttuale.succ.prec = tmp
		}
		posAttuale.succ = tmp
		lista.dim += 1

		//split := strings.Fields(posAttuale.elemento)

		vecchio := posAttuale.elemento
		tmp.sostituisciStringa(vecchio[pos+1 : posAttuale.size]) //TODO rotto
		posAttuale.sostituisciStringa(vecchio[0 : pos+1])

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

// remove dalla stringa corrente
func removeFromString(s string, s_size int, pos int, howmany int) string {

	return strings.Join([]string{s[0:pos], s[pos+howmany : s_size]}, "")

}

func removeFromRunes(origin []rune, s_size int, pos int, howmany int) []rune {

	return append(origin[0:pos], origin[pos+howmany:s_size]...)
}

func (lista *SuperString) delElem(pos int, howmany int) {
	if pos < 0 {
		fmt.Println("che cazzo stai cercando di eliminare???")
		return
	}
	//roRebuildTotal = true

	posAttuale := lista.testa

	for pos >= posAttuale.size {
		pos -= posAttuale.size
		posAttuale = posAttuale.succ
	}

	//elimina prima stringa
	quanti_eliminare := howmany
	if pos+howmany <= posAttuale.size { //elimina solo posizione attuale
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

	for posAttuale.succ != nil && quanti_eliminare >= posAttuale.size {
		quanti_eliminare -= posAttuale.size

		tmp := posAttuale.succ
		lista.delSingleElem(posAttuale)
		posAttuale = tmp
	}

	//elimina ultima stringa
	if posAttuale.size == quanti_eliminare {
		lista.delSingleElem(posAttuale)
	} else {
		posAttuale.elemento = removeFromRunes(posAttuale.elemento, posAttuale.size, 0, quanti_eliminare)
		posAttuale.size -= quanti_eliminare
	}
}

/*
void SuperString::delElem(int pos, const int howmany){
.....
    //elimina ultima stringa
    if (posAttuale->size == quanti_eliminare){
        delSingleElem(posAttuale);
    }
    else {
        posAttuale->elem->remove(0,quanti_eliminare);
        posAttuale->size -= quanti_eliminare;
    }


}*/

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
		lista.testa = NewElemSuperString(nil, nil)
	} else {
		lista.dim--
	}

	//delete elemento
	elemento = nil
}
