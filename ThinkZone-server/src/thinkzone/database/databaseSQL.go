// database
package database

import (
	"database/sql"
	"fmt"
	_ "github.com/jbarham/gopgsqldriver"
	"strconv"
	"strings"
	"thinkzone/logs"
)

// variabili per gestire il database
var (
	//Operazioni per inserire utente
	insertUser   string    = "INSERT INTO t_user VALUES ($1, $2, $3)"
	insertUserOp *sql.Stmt = nil

	//Operazioni per gestire post
	insertPost   string    = "INSERT INTO post VALUES ($1, $2, $3, $4, $5, $6)"
	insertPostOp *sql.Stmt = nil

	//Operazioni per gestire le conversazioni
	insertConv   string    = "INSERT INTO conversation VALUES ($1)"
	insertConvOp *sql.Stmt = nil

	// Script in SQL per creare le tabelle nel database 
	createDbSqlScript []string = []string{
		"CREATE TABLE t_user ( id INT PRIMARY KEY, username CHAR(64) NOT NULL UNIQUE, password BYTEA NOT NULL)",
		"CREATE TABLE conversation ( id INT PRIMARY KEY )",
		"CREATE TABLE archive ( t_user INT, conversation INT, PRIMARY KEY (t_user,conversation), FOREIGN KEY (t_user) REFERENCES t_user(id), FOREIGN KEY (conversation) REFERENCES conversation(id))",
		"CREATE TABLE post ( id INT NOT NULL, conversation INT NOT NULL, text TEXT, pather INT, first_response INT, second_response INT, PRIMARY KEY (id, conversation), FOREIGN KEY (conversation) REFERENCES conversation(id))",                           //, FOREIGN KEY (pather) REFERENCES post(id), FOREIGN KEY (first_response) REFERENCES post(id), FOREIGN KEY (second_response) REFERENCES post(id))",
		"CREATE TABLE author ( t_user INT NOT NULL, post INT NOT NULL, conversation INT NOT NULL, PRIMARY KEY (t_user,conversation, post), FOREIGN KEY (t_user) REFERENCES t_user(id), FOREIGN KEY (post, conversation) REFERENCES post(id, conversation))"} //, FOREIGN KEY (conversation) REFERENCES post(conversation))"}

)

func CreateDataBase() error {
	var err error

	db, err = sql.Open("postgres", "dbname=thinkzoneDB user=thinkzone")
	if err != nil {
		logs.Error("Impossibile aprire il database")
		return err
	}

	for i := range createDbSqlScript {
		_, err = db.Exec(createDbSqlScript[i])
		if err != nil {
			if "result error: ERROR:  relation \"t_user\" already exists\n" != err.Error() {
				logs.Error("Impossibile creare le tabelle del database\nmotivo: _", err.Error(), "_")
				return err
			}
			break
		}

	}

	insertUserOp, err = db.Prepare(insertUser)
	if err != nil {
		//		logs.Error("Impossibile salvare gli utenti nel database\nmotivo: ", err.Error())
		return err
	}

	insertPostOp, err = db.Prepare(insertPost)
	if err != nil {
		//		logs.Error("Impossibile salvare i post nel database\nmotivo: ", err.Error())
		return err
	}

	insertConvOp, err = db.Prepare(insertConv)
	if err != nil {
		//		logs.Error("Impossibile salvare le conversazioni nel database\nmotivo: ", err.Error())
		return err
	}

	return nil

}

func salvaUtente(user *User) error {

	_, err := insertUserOp.Exec(user.ID, user.Username, user.password)
	if err != nil {
		logs.Error("Impossibile salvare ", user.Username, " nel database")
		return err
	}

	return nil
}

/*func SalvaUtenti() error {
	logs.Log("salvo1")
	operazione, err := db.Prepare(insertUser)
	if err != nil {
		logs.Error("Impossibile salvare gli utenti nel database")
		return err
	}
	defer operazione.Close()

	logs.Log("salvo2")
	for _, user := range Data.userIDtoUser {	
		_, err = operazione.Exec(user.ID, user.Username, user.password)
		if err != nil {
			logs.Error("Impossibile salvare ", user.Username, " nel database")
			return err
		}
	}

	logs.Log("salvo3")

	return nil
}*/

func (datab *DatabaseRegistration) CaricaUtenti() error {

	rows, err := db.Query("SELECT * FROM t_user")
	if err != nil {
		return err
	}

	logs.Log("caricamento utenti")
	for ; rows.Next(); datab.contatore++ {
		var username string
		var userID int
		var password string

		err := rows.Scan(&userID, &username, &password)
		if err != nil {
			return err
		}

		username = strings.TrimSpace(username)
		password = strings.TrimSpace(password)

		user := new(User)
		user.Username = username
		user.ID = userID
		user.password = []byte(password)

		fmt.Println("caricato vecchio utente: ", username, " ID: ", strconv.Itoa(userID), "\npassword: ", password)

		datab.userNameToId[username] = user
		datab.userIDtoUser[userID] = user

	}
	return nil
}

func (datab *DatabaseRegistration) salvaPost(conv *Conversation, post *Post) error {
	var idPadre, idRPri, idRSec int = -1, -1, -1
	if post.padre != nil {
		idPadre = post.padre.idPost
	}
	if post.rispostaPrincipale != nil {
		idRPri = post.rispostaPrincipale.idPost
	}
	if post.rispostaSecondaria != nil {
		idRSec = post.rispostaSecondaria.idPost
	}

	_, err := insertPostOp.Exec(post.idPost, conv.ID, post.testo.GetComplete(false), idPadre, idRPri, idRSec)
	if err != nil {
		logs.Error("Impossibile salvare il post ", strconv.Itoa(post.idPost), " nel database")
		return err
	}

	return nil
}

func (conv *Conversation) salvaTuttiIPost(data *DatabaseRegistration) error {
	messaggio := "salvataggio di tutti i post sul database"
	logs.Log(messaggio)
	err := conv.testaPost.salvaPostRic(conv, data)

	if err != nil {
		logs.Error(messaggio + " fallito")
	} else {
		logs.Log(messaggio + " riuscito")
	}

	return err
}

func (post *Post) salvaPostRic(conv *Conversation, data *DatabaseRegistration) error {

	err := data.salvaPost(conv, post)
	if err != nil {
		return err
	}

	if post.rispostaPrincipale != nil {
		err := post.rispostaPrincipale.salvaPostRic(conv, data)
		if err != nil {
			return err
		}
	}
	if post.rispostaSecondaria != nil {
		err := post.rispostaSecondaria.salvaPostRic(conv, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (conv *Conversation) salvaTutteLeConversazioni() error {
	_, err := insertConvOp.Exec(0)
	if err != nil {
		logs.Error("Impossibile salvare il post", strconv.Itoa(conv.ID), " nel database")
		return err
	}
	return nil

}
