// superstring_test
package main

import (
	"fmt"
	"strings"
	"testing"
)

func errore(test *testing.T, s string) {
	fmt.Println("######Errore:", s)
	test.Fail()
}

func equal(s string, t string) bool {
	return strings.EqualFold(s, t)
}

func getTestEqual(nome string, test *testing.T, lista *SuperString) func(neutra string, conSep string) {

	i := 1

	return func(neutra string, conSep string) {
		stringa := lista.GetComplete(false)
		if !equal(stringa, neutra) {
			errore(test, stringa)
		} else {
			fmt.Println("Corretto:", stringa)
		}

		stringa = lista.GetComplete(true)
		if !equal(stringa, conSep) {
			errore(test, stringa)
		} else {
			fmt.Println("Corretto:", stringa)
		}
		fmt.Println("#", i, " eseguita su: ", nome)
		i++
	}
}

func TestSuperString(test *testing.T) {

	prova := NewSuperString()
	testEqual := getTestEqual("primo", test, prova)

	testEqual("", "[0]")

	prova.insElem("99", -1)
	testEqual("99", "99[2]")

	prova.insElem("11!", 0)
	testEqual("11!99", "11![3]99[2]")

	prova.insElem("###", 3)
	testEqual("11!###99", "11!###[6]99[2]")

	prova.insElem("tro", 2)
	testEqual("11tro!###99", "11tro[5]!###[4]99[2]")

	prova.delElem(2, 1)
	testEqual("1tro!###99", "1tro[5]!###[4]99[2]")

}
