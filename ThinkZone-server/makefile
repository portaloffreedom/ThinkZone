# -*- Makefile -*-

all: ThinkZone

ThinkZone: comm.go main.go comm.output
	go build
	
comm.go: comunication_grammar.y 
	go tool yacc -o="comm.go" -v="comm.output" cominication_grammar.y
	
comunication_grammar.y:

.PHONY: clean cleanall

cleanall: clean
	rm ThinkZone

clean:
	rm comm.go comm.output