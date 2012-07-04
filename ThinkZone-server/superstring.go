// superstring
package main

import (
	"fmt"
	"strings"
)

type elemSuperString struct {
	elemento   string
	size       int
	succ, prec *elemSuperString
}

type SuperString struct {
	testa *elemSuperString
	dim   int
}

func NewElemSuperString(prec *elemSuperString, succ *elemSuperString) *elemSuperString {

	elem := new(elemSuperString)

	elem.elemento = ""
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
func (elem *elemSuperString) sostituisciStringa(nuova string) {
	elem.elemento = nuova
}

func (lista *SuperString) GetCompleteWithSeparators(separator string) string {

	tmp := lista.testa
	dim := lista.dim

	var totalVector []string = make([]string, dim)

	tmp = lista.testa
	for j := 0; j < dim; j++ {
		totalVector[j] = tmp.elemento
		tmp = tmp.succ
	}

	return strings.Join(totalVector, separator)

}

func (lista *SuperString) GetComplete() string {

	return lista.GetCompleteWithSeparators("")

}

func (lista *SuperString) insElem(appendStr string, pos int) {

	if lista.testa == nil {
		lista.testa = NewElemSuperString(nil, nil)
	}

	dimStr := len(appendStr)
	if dimStr == 0 {
		return
	}

	//TODO eliminare questa parte
	if pos == 0 {
		lista.testa = NewElemSuperString(nil, lista.testa)
		lista.testa.succ.prec = lista.testa
		lista.testa.elemento = appendStr
		lista.dim = dimStr

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
		fmt.Println("SuperString: inserisco in fondo all'ultima stringa")
		coda := lista.testa
		for coda.succ != nil {
			coda = coda.succ
		}

		coda.elemento = strings.Join([]string{coda.elemento, appendStr}, "")
		coda.size += dimStr

		return
	}

	if pos == 0 {
		posAttuale = posAttuale.prec
	} else {
		fmt.Println("SuperString: scindo due strighe")
		tmp := NewElemSuperString(posAttuale, posAttuale.succ)
		posAttuale.succ = tmp

		//split := strings.Fields(posAttuale.elemento)

		vecchio := posAttuale.elemento
		tmp.sostituisciStringa(vecchio[0:pos])
		tmp.sostituisciStringa(vecchio[pos+1 : dimStr])

		tmp.size = posAttuale.size - pos
		posAttuale.size = pos
	}

	posAttuale.elemento = strings.Join([]string{posAttuale.elemento, appendStr}, "")
	posAttuale.size += dimStr

	return

}

//func (lista *SuperString) insElemChar(appendChar byte, pos int) {
//	lista.insElem(append("", appendChar), pos)
//}

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
	}

	//delete elemento
	elemento = nil
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
	//quanti_eliminare := howmany
}

/*
void SuperString::delElem(int pos, const int howmany){
.....
    //elimina prima stringa
    int quanti_eliminare = howmany;
    if (pos+howmany <= posAttuale->size) {   //elimina solo posizione attuale
        posAttuale->elem->remove(pos,quanti_eliminare);
        posAttuale->size -= quanti_eliminare;
        if (posAttuale->size == 0)
            delSingleElem(posAttuale);
        return;
    }
    quanti_eliminare = posAttuale->size-pos;
    posAttuale->elem->remove(pos,quanti_eliminare);
    posAttuale->size -= (quanti_eliminare);

    if (posAttuale->size == 0) {
        elemListaString* temp = posAttuale->succ;
        delSingleElem(posAttuale);
        posAttuale = temp;
        if (posAttuale == NULL)
            posAttuale = testa;
    }
    else
        posAttuale = posAttuale->succ;

    //elimina stringhe di mezzo
    quanti_eliminare = howmany - quanti_eliminare;

    while (posAttuale->succ != NULL && quanti_eliminare >= posAttuale->size) {
        quanti_eliminare -= posAttuale->size;

        elemListaString* temp = posAttuale->succ;
        delSingleElem(posAttuale);
        posAttuale = temp;
    }

    //elimina ultima stringa

    if (posAttuale->size == quanti_eliminare){
        delSingleElem(posAttuale);
    }
    else {
        posAttuale->elem->remove(0,quanti_eliminare);
        posAttuale->size -= quanti_eliminare;
    }


}*/
