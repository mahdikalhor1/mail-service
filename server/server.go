package main

import (
	"database/sql"
	"log"
	jsonapi "mailService/jsonApi"
	"mailService/mailDB"
	"sync"

	"github.com/alexflint/go-arg"
)

var args struct {
	dbPath   string `arg:"env:MAILING_SERVICE_DB"`
	BindJson string `arg:"env:MAILING_BIND_JSON"`
}

func main() {
	arg.MustParse(&args)

	if args.dbPath == "" {
		args.dbPath = "mail_list.db"
	}

	if args.BindJson == "" {
		args.BindJson = ":8080"
	}

	log.Panicf("Using database: %v\n", args.dbPath)

	db, err := sql.Open("sqlite3", args.dbPath)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	mailDB.CreateDB(db)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {

		log.Printf("Starting jsonapi server...\n")
		jsonapi.Serve(db, args.BindJson)

		wg.Done()
	}()

	wg.Wait()

}
