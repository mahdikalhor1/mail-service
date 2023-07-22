package mailDB

import (
	"database/sql"
	"log"
	"time"

	"github.com/mattn/go-sqlite3"
)


type EmailEntry struct{
	Id string
	Email string
	ConfirmedAt *time.Time
	OptOut bool
}

type GetEmailBatchParams struct{
	Page int
	Count int
}


func CreateDB(db *sql.DB){

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

func InsertEmail(db *sql.DB, email string) error{

	_, err := db.Exec(`
	INSERT INTO emails(email, confirmed_at, opt_out)
	VALUES(?, 0, false)`, email)

	if err != nil{
		log.Println(err)
	}

	return nil
}

func GetEmail(db *sql.DB, email string)(*EmailEntry, error){
	
	rows, err := db.Query(`
	SELECT *
	FROM emails
	WHERE email = ?`, email)

	defer rows.Close()

	if err != nil{

		log.Println(err)
		return nil, err
	}

	for rows.Next(){
		return getEntryFromRow(rows)
	}

	return nil, nil
}

func UpdateEmail(db *sql.DB, emailEntry *EmailEntry) error{

	confirmTimeInt := emailEntry.ConfirmedAt.Unix()

	_, err := db.Query(`
	INSERT INTO emails(email, confirmed_at, opt_out)
	VALUES(?, ?, ?)
	ON CONFLICT(email) DO
	UPDATE SET
	confirmed_at=? opt_out=?`,emailEntry.Email, confirmTimeInt,
	emailEntry.OptOut, confirmTimeInt, emailEntry.OptOut)

	if err != nil{
		log.Println(err)
		return err
	}

	return nil
}

func DeleteEmail(db *sql.DB, email string) error{

	_, err := db.Query(`
	UPDATE emails
	SET opt_out = true
	WHERE email = ?`, email)

	if err != nil{
		log.Println(err)
	}

	return nil
}


func GetEmailBatch(db *sql.DB, gp GetEmailBatchParams)([]EmailEntry, error){

	emails := make([]EmailEntry, 0, gp.Count)


	rows, err := db.Query(`
	SELECT *
	FROM emails
	WHERE opt_out = false
	ORDER BY id ASC
	LIMIT ? OFFSET ?`, gp.Count, gp.Count * (gp.Page - 1))

	defer rows.Close()

	if err != nil{
		
		log.Println(err)
		return nil, err
	}

	for rows.Next(){

		entry, err := getEntryFromRow(rows)

		if err != nil{
			log.Println(err)
			
			return nil, err
		}

		emails = append(emails, *entry)
	}

	return emails, nil
}