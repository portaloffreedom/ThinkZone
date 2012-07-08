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
		fmt.Println("#", i)

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
		fmt.Println("#", i, "eseguita su:", nome)
		i++
	}
}

func TestSuperString(test *testing.T) {

	prova := NewSuperString()
	testEqual := getTestEqual("primo", test, prova)

	//#1
	testEqual("", "[0]")

	//#2 ---insert---
	prova.insElem("99", 0)
	testEqual("99", "99[2]")

	//#3
	prova.insElem("11!", 0)
	testEqual("11!99", "11![3]99[2]")

	//#4
	prova.insElem("###", 3)
	testEqual("11!###99", "11!###[6]99[2]")

	//#5
	prova.insElem("tro", 2)
	testEqual("11tro!###99", "11tro[5]!###[4]99[2]")

	//#6 ----delete----
	prova.delElem(1, 1)
	testEqual("1tro!###99", "1tro[4]!###[4]99[2]")

	//#7
	prova.delElem(4, 1)
	testEqual("1tr!###99", "1tr[3]!###[4]99[2]")

	//#8
	prova.delElem(5, 3)
	testEqual("1tr!99", "1tr[3]![1]99[2]")

	//#9
	prova.delElem(4, 1)
	testEqual("1tr99", "1tr[3]99[2]")

	//#10
	prova.delElem(1, 3)
	testEqual("99", "99[2]")

	//#11
	prova.insElem("porcoDioZoccolo", 0)
	prova.insElem(" ", 8)
	prova.insElem(" ", 5)
	testEqual("porco Dio Zoccolo99", "porco [6]Dio [4]Zoccolo[7]99[2]")

	//#12
	prova.delElem(4, 15)
	testEqual("por9", "por[3]9[1]")

	fmt.Println()
}

func TestSuperString2(test *testing.T) {
	prova := NewSuperString()
	testEqual := getTestEqual("secondo", test, prova)

	//#1
	testEqual("", "[0]")

	prova.insElem("Australopitecus---", 0)
	testEqual("Australopitecus---", "Australopitecus---[18]")

	prova.insElem("0", 16)
	testEqual("Australopitecus-0--", "Australopitecus-0[17]--[2]")

}
