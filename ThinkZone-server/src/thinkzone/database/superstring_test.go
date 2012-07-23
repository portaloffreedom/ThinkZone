// superstring_test
package database

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
	prova.InsStringElem("99", 0)
	testEqual("99", "99[2]")

	//#3
	prova.InsStringElem("12!", 0)
	testEqual("12!99", "12![3]99[2]")

	//#4
	prova.InsStringElem("###", 3)
	testEqual("12!###99", "12!###[6]99[2]")

	//#5
	prova.InsStringElem("tro", 2)
	testEqual("12tro!###99", "12tro[5]!###[4]99[2]")

	//#6 ----delete----
	prova.DelElem(1, 1)
	testEqual("1tro!###99", "1tro[4]!###[4]99[2]")

	//#7
	prova.DelElem(4, 1)
	testEqual("1tro###99", "1tro[4]###[3]99[2]")

	//#8
	prova.DelElem(4, 3)
	testEqual("1tro99", "1tro[4]99[2]")

	//#9
	prova.DelElem(3, 1)
	testEqual("1tr99", "1tr[3]99[2]")

	//#10
	prova.DelElem(0, 3)
	testEqual("99", "99[2]")

	//#11
	prova.InsStringElem("porcoDioZoccolo", 0)
	prova.InsStringElem(" ", 8)
	prova.InsStringElem(" ", 5)
	testEqual("porco Dio Zoccolo99", "porco [6]Dio [4]Zoccolo[7]99[2]")

	//#12
	prova.DelElem(3, 15)
	testEqual("por9", "por[3]9[1]")

	fmt.Println()
}

func TestSuperString2(test *testing.T) {
	prova := NewSuperString()
	testEqual := getTestEqual("secondo", test, prova)

	//#1
	testEqual("", "[0]")

	//#2
	prova.InsStringElem("Australopitecus---", 0)
	testEqual("Australopitecus---", "Australopitecus---[18]")

	//#3
	prova.InsStringElem("0", 16)
	testEqual("Australopitecus-0--", "Australopitecus-0[17]--[2]")

}
