// cominication_grammar.y
%{
	
package main

import (
	"fmt"
)
%}

%token DIGIT

%%
line	: expr '\n'		{ fmt.Println($1) }
		;
expr	: expr '+' term	{ $$ = $1 + $3 }
		| term
		;
term	: term '*' factor { $$ = $1 * $3 }
		| factor
		;
factor	: '(' expr ')'	{ $$ = $2 }
		| DIGIT
		;
%%

//type yySymType int

type Lexer struct {
}

func (lex *Lexer) Lex(lval *yySymType) int {
	return int(*lval)
}

func (lex *Lexer) Error(s string) {
	fmt.Println("Errore!")
	fmt.Println(s)
}

/*
func comm() {
	fi := bufio.NewReader(os.NewFile(0, "stdin"))

	for {
		var eqn string
		var ok bool

		fmt.Printf("equation: ")
		if eqn, ok = readline(fi); ok {
			CalcParse(&CalcLex{s: eqn})
		} else {
			break
		}
	}
}
*/