// superstring_test
package main

import (
	"fmt"
	"strings"
	"testing"
)

func errore(test *testing.T, s string) {
	fmt.Printf("######Errore:|%v|\n", s)
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
	prova.insStringElem("99", 0)
	testEqual("99", "99[2]")

	//#3
	prova.insStringElem("12!", 0)
	testEqual("12!99", "12![3]99[2]")

	//#4
	prova.insStringElem("###", 3)
	testEqual("12!###99", "12!###[6]99[2]")

	//#5
	prova.insStringElem("tro", 2)
	testEqual("12tro!###99", "12tro[5]!###[4]99[2]")

	//#6 ----delete----
	prova.delElem(1, 1)
	testEqual("1tro!###99", "1tro[4]!###[4]99[2]")

	//#7
	prova.delElem(4, 1)
	testEqual("1tro###99", "1tro[4]###[3]99[2]")

	//#8
	prova.delElem(4, 3)
	testEqual("1tro99", "1tro[4]99[2]")

	//#9
	prova.delElem(3, 1)
	testEqual("1tr99", "1tr[3]99[2]")

	//#10
	prova.delElem(0, 3)
	testEqual("99", "99[2]")

	//#11
	prova.insStringElem("porcoDioZoccolo", 0)
	prova.insStringElem(" ", 8)
	prova.insStringElem(" ", 5)
	testEqual("porco Dio Zoccolo99", "porco [6]Dio [4]Zoccolo[7]99[2]")

	//#12
	prova.delElem(3, 15)
	testEqual("por9", "por[3]9[1]")

	fmt.Println()
}

func TestSuperString2(test *testing.T) {
	prova := NewSuperString()
	testEqual := getTestEqual("secondo", test, prova)

	//#1
	testEqual("", "[0]")

	//#2
	prova.insStringElem("Australopitecus---", 0)
	testEqual("Australopitecus---", "Australopitecus---[18]")

	//#3
	prova.insStringElem("0", 16)
	testEqual("Australopitecus-0--", "Australopitecus-0[17]--[2]")

}
