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

func CreateDataBase() error {
	var err error

	db, err = sql.Open("postgres", "dbname=thinkzoneDB user=thinkzone")
	if err != nil {
		logs.Error("Impossibile creare il database")
		return err
	}

	for i := range createDbSqlScript {
		_, err = db.Exec(createDbSqlScript[i])
		if err != nil {
			logs.Error("Impossibile creare le tabelle del database")
			return err
		}

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
