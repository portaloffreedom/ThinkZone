package main

import (
	"fmt"
	"net"
	"math"
	"bufio"
	"database/sql"
	_ "github.com/jbarham/gopgsqldriver"
)

func createTestErr() func(err error) {
	contatore := 0
	
	return func(err error) {
		contatore++
		fmt.Printf("passaggio numero %v: ",contatore)
		if err != nil {
			fmt.Println("errore")
			fmt.Println(err)
		} else {
			fmt.Println("passato")
		}
	}
}

type rec struct {
	tf  bool
	i32 int
	i64 int64
	s   string
	b   []byte
}

var testTuples = []rec{
	{false, math.MinInt32, math.MinInt64, "hello world", []byte{0xDE, 0xAD}},
	{true, math.MaxInt32, math.MaxInt64, "Γεια σας κόσμο", []byte{0xBE, 0xEF}},
}

func testDB() {
	testErr := createTestErr()
	
	db,err := sql.Open("postgres", "dbname=testlibreoffice")
	testErr(err)
	
	_,err = db.Exec("CREATE TABLE gopq_test (tf bool, i32 int, i64 bigint, s text)")
	testErr(err)
	
	defer db.Exec("DROP TABLE gopq_test")

	// Insert test rows.
	stmt, err := db.Prepare("INSERT INTO gopq_test VALUES ($1, $2, $3, $4)")
	testErr(err)
	defer stmt.Close()
	for _, row := range testTuples {
		_, err = stmt.Exec(row.tf, row.i32, row.i64, row.s)
		testErr(err)
	}

	// Verify that all test rows were inserted.
	rows, err := db.Query("SELECT COUNT(*) FROM gopq_test")
	testErr(err)
	if !rows.Next() {
		fmt.Println("Result.Next failed")
	}
	var count int
	err = rows.Scan(&count)
	testErr(err)
	if count != len(testTuples) {
		fmt.Println("invalid row count %d, expected %d", count, len(testTuples))
	}
	rows.Close()

	// Retrieve inserted rows and verify inserted values.
	rows, err = db.Query("SELECT * FROM gopq_test")
	testErr(err)
	for i := 0; rows.Next(); i++ {
		var tf bool
		var i32 int
		var i64 int64
		var s string

		err := rows.Scan(&tf, &i32, &i64, &s)
		testErr(err)
		if err != nil {
			fmt.Println("scan error:", err)
		}
		if tf != testTuples[i].tf {
			fmt.Println("bad bool")
		}
		if i32 != testTuples[i].i32 {
			fmt.Println("bad int32")
		}
		if i64 != testTuples[i].i64 {
			fmt.Println("bad int64")
		}
		if s != testTuples[i].s {
			fmt.Println("bad string")
		}
	}
	rows.Close()
}

func cacca(cosa string) func () {
	contatore := 0
        azione := cosa
	return func (){
		contatore++
		contatorelocale := contatore
		//for i:=10; i<1000; i++ {
		//for {
			//a := 1
			//b := a
			fmt.Printf("connessione numero %v, %v\n",contatorelocale,azione)
		//}
	}
}

func main() {
	fmt.Println("ciao mondo")
	ln,err := net.Listen("tcp",":4000")
	if err != nil {
		fmt.Println("Errore nell'aprire la connessione")
		// handle error
	}
	connessione_aperta := cacca("aperta { ----")
    connessione_chiusa := cacca("chiusa }")
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("### connessione persa ###")
			continue
		}
                connessione_aperta()
                input := bufio.NewReader(conn)
                
		//go handleConnection(conn)
		fmt.Fprintf(conn,"culo\n")
		var pirla string
                //c := make (chan string)
		pirla,_ = input.ReadString('\n')
		//fmt.Fprintf(conn,pirla)
                fmt.Print(pirla)
		//fmt.Fprintf(conn,"culo2\n")
		conn.Close()
		connessione_chiusa()
		//testDB()
	}
}