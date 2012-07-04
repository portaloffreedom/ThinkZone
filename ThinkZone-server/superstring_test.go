// superstring_test
package main

import (
	"fmt"
	"strings"
	"testing"
)

func errore(test *testing.T, s string) {
	fmt.Println(":Errore:", s)
	test.Fail()
}

func equal(s string, t string) bool {
	return strings.EqualFold(s, t)
}

func getTestEqual(nome string, test *testing.T, lista *SuperString) func(neutra string, conSep string) {

	i := 1

	return func(neutra string, conSep string) {
		stringa := lista.GetComplete()
		if !equal(stringa, neutra) {
			errore(test, stringa)
		}

		stringa = lista.GetCompleteWithSeparators("|")
		if !equal(stringa, conSep) {
			errore(test, stringa)
		}
		fmt.Println("#", i, " eseguita su: ", nome)
		i++
	}
}

func TestSuperString(test *testing.T) {

	prova := NewSuperString()
	testEqual := getTestEqual("primo", test, prova)

	testEqual("", "")

	prova.insElem("99", 0)
	testEqual("99", "99|")

	//	prova.insElem("11!", 0)
	//	testEqual("11!99", "11!|99|")

	prova.insElem("trot", 0)
	//	testEqual("a11!99", "a|11!|99|")
	testEqual("trot99", "trot|99|")

}
