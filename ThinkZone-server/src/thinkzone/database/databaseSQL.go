// database
package database

import (
	"container/list"
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
	updatePost   string    = "UPDATE post SET text = $5,  father = $2, first_response = $3, second_response =$4 WHERE id = $1;"
	updatePostOp *sql.Stmt = nil

	//Operazioni per gestire le conversazioni
	insertConv   string    = "INSERT INTO conversation VALUES ($1)"
	insertConvOp *sql.Stmt = nil
	updateConv   string    = "UPDATE conversation SET id = $1"
	updateConvOp *sql.Stmt = nil

	// Script in SQL per creare le tabelle nel database 
	createDbSqlScript []string = []string{
		"CREATE TABLE t_user ( id INT PRIMARY KEY, username CHAR(64) NOT NULL UNIQUE, password BYTEA NOT NULL)",
		"CREATE TABLE conversation ( id INT PRIMARY KEY )",
		"CREATE TABLE archive ( t_user INT, conversation INT, PRIMARY KEY (t_user,conversation), FOREIGN KEY (t_user) REFERENCES t_user(id), FOREIGN KEY (conversation) REFERENCES conversation(id))",
		"CREATE TABLE post ( id INT NOT NULL, conversation INT NOT NULL, text TEXT, father INT, first_response INT, second_response INT, PRIMARY KEY (id, conversation), FOREIGN KEY (conversation) REFERENCES conversation(id))",                           //, FOREIGN KEY (pather) REFERENCES post(id), FOREIGN KEY (first_response) REFERENCES post(id), FOREIGN KEY (second_response) REFERENCES post(id))",
		"CREATE TABLE author ( t_user INT NOT NULL, post INT NOT NULL, conversation INT NOT NULL, PRIMARY KEY (t_user,conversation, post), FOREIGN KEY (t_user) REFERENCES t_user(id), FOREIGN KEY (post, conversation) REFERENCES post(id, conversation))"} //, FOREIGN KEY (conversation) REFERENCES post(conversation))"}

)

func CreateDataBaseSQL() error {
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
		//		logs.Error("Impossibile creare i post nel database\nmotivo: ", err.Error())
		return err
	}

	updatePostOp, err = db.Prepare(updatePost)
	if err != nil {
		//		logs.Error("Impossibile salvare i post nel database\nmotivo: ", err.Error())
		return err
	}

	insertConvOp, err = db.Prepare(insertConv)
	if err != nil {
		//		logs.Error("Impossibile salvare le conversazioni nel database\nmotivo: ", err.Error())
		return err
	}

	updateConvOp, err = db.Prepare(updateConv)
	if err != nil {
		return err
	}

	logs.AggiungiAzioneDiChiusura(func() {
		if insertUserOp != nil {
			insertUserOp.Close()
		}
		if insertConvOp != nil {
			err := MainConv.salvaTutteLeConversazioniSQL()
			if err != nil {
				logs.Error("Impossibile salvare tutti i post\nmotivo: ", err.Error())
			}
			insertConvOp.Close()
		}
		if insertPostOp != nil {
			insertPostOp.Close()
		}
		if updatePostOp != nil {
			updatePostOp.Close()
		}
		if updateConvOp != nil {
			updateConvOp.Close()
		}
		return
	})

	return nil

}

func salvaUtenteSQL(user *User) error {

	_, err := insertUserOp.Exec(user.ID, user.Username, user.password)
	if err != nil {
		logs.Error("Impossibile salvare ", user.Username, " nel database")
		return err
	}

	return nil
}

func (datab *DatabaseRegistration) CaricaUtentiSQL() error {

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

func (datab *DatabaseRegistration) creaPostSQL(conv *Conversation, post *Post) error {
	logs.Log("creo post ", strconv.Itoa(post.idPost))

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
		logs.Error("Impossibile creare il post ", strconv.Itoa(post.idPost), " nel database")
		return err
	}

	return nil

}

func (datab *DatabaseRegistration) salvaPostSQL(conv *Conversation, post *Post) error {
	logs.Log("salvo post ", strconv.Itoa(post.idPost))
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

	//campi:					id		padre	first_response 		second_response 	text
	_, err := updatePostOp.Exec(post.idPost, idPadre, idRPri, idRSec, post.testo.GetComplete(false))
	if err != nil {
		logs.Error("Impossibile salvare il post ", strconv.Itoa(post.idPost), " nel database")
		return err
	}

	return nil
}

func (conv *Conversation) salvaTuttiIPostSQL(data *DatabaseRegistration) error {
	messaggio := "salvataggio di tutti i post sul database"
	logs.Log(messaggio)
	err := conv.testaPost.salvaPostRicSQL(conv, data)

	if err != nil {
		logs.Error(messaggio + " fallito")
	} else {
		logs.Log(messaggio + " riuscito")
	}

	return err
}

func (post *Post) salvaPostRicSQL(conv *Conversation, data *DatabaseRegistration) error {

	err := data.salvaPostSQL(conv, post)
	if err != nil {
		return err
	}

	if post.rispostaPrincipale != nil {
		err := post.rispostaPrincipale.salvaPostRicSQL(conv, data)
		if err != nil {
			return err
		}
	}
	if post.rispostaSecondaria != nil {
		err := post.rispostaSecondaria.salvaPostRicSQL(conv, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (conv *Conversation) salvaTutteLeConversazioniSQL() error {
	_, err := updateConvOp.Exec(0)
	if err != nil {
		logs.Error("Impossibile salvare la conversazione ", strconv.Itoa(conv.ID), " nel database")
		return err
	}
	return conv.salvaTuttiIPostSQL(&Data)
}

func (datab *DatabaseRegistration) caricaConversazioni() error {

	rows, err := db.Query("SELECT * FROM conversation")
	if err != nil {
		return err
	}

	logs.Log("caricamento conversazioni")
	for ; rows.Next(); datab.contatore++ {
		var conversationID int

		err := rows.Scan(&conversationID)
		if err != nil {
			return err
		}

		MainConv = new(Conversation)
		MainConv.ID = conversationID
		MainConv.connected = make(map[int]*convUser)
		MainConv.postMap = make(map[int]*Post)
		MainConv.contatorePost = 0

		err = MainConv.caricaPost()
		if err != nil {
			return err
		}

		fmt.Println("caricata conversazione:", conversationID)

	}
	return nil
}

func (conv *Conversation) caricaPost() error {

	query, err := db.Prepare("SELECT id,father,text,first_response,second_response FROM post WHERE conversation = $1")
	if err != nil {
		return err
	}
	rows, err := query.Query(conv.ID)
	if err != nil {
		return err
	}

	parentPostMap := make(map[int]*Post)
	firSonPostMap := make(map[int]*Post)
	secSonPostMap := make(map[int]*Post)
	logs.Log("caricamento post della conversazione ", strconv.Itoa(conv.ID))
	for ; rows.Next(); conv.contatorePost++ {
		var id, parent, first_response, second_response int
		var text string

		err := rows.Scan(&id, &parent, &text, &first_response, &second_response)
		if err != nil {
			return err
		}

		post := new(Post)
		post.idPost = id
		post.testo = NewSuperString()
		post.testo.InsStringElem(text, 0)
		conv.postMap[id] = post
		post.writers = list.New()

		if id == 0 {
			conv.testaPost = post
		}

		parentPostMap[parent] = post
		firSonPostMap[first_response] = post
		secSonPostMap[second_response] = post

		fmt.Println("caricato vecchio post: ", strconv.Itoa(id), "\ntesto: ", text)

	}
	for k, post := range parentPostMap {
		post.padre = conv.postMap[k]
	}
	for k, post := range firSonPostMap {
		post.rispostaPrincipale = conv.postMap[k]
	}
	for k, post := range secSonPostMap {
		post.rispostaSecondaria = conv.postMap[k]
	}

	return nil
}
