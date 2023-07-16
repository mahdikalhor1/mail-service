package mailDB

import (
	"database/sql"
	"log"
	"time"

	"github.com/mattn/go-sqlite3"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
)


type EmailEntry struct{
	Id string
	Email string
	ConfirmedAt *time.Time
	OptOut bool
}


func createDB(db *sql.DB){

	_, err := db.Exec(`
	CREATE TABLE emails(
		id TEXT PRIMARY KEY,
		email TEXT UNIQUE
		confirmed_at INTEGER,
		opt_out INTEGER
	);
	`)
	
	sqlErr, ok := err.(sqlite3.Error)

	if ok {
		if sqlErr.code != 1{
			
			//code one indicates that the table is exists already
			log.Fatal(sqlErr)
		}
	}else{
		log.Fatal(err)
	}
}

func getEntryFromRow(row *sql.Rows)(*EmailEntry, error){

	var id string
	var email string
	var confirmedatInt int64
	var optOut bool

	err := row.Scan(&id, &email, &confirmedatInt, &optOut)

	if err != nil{
		log.Println(err)

		return nil, err
	}

	confirmedat := time.Unix(confirmedatInt, 0)

	return &EmailEntry{id, email, &confirmedat, optOut}, nil
}