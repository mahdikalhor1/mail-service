package main

import (
	"database/sql"
	"log"
	grpcapi "mailService/grpcApi"
	jsonapi "mailService/jsonApi"
	"mailService/mailDB"
	"sync"

	"github.com/alexflint/go-arg"
)

var args struct {
	dbPath   string `arg:"env:MAILING_SERVICE_DB"`
	BindJson string `arg:"env:MAILING_BIND_JSON"`
	BindGrpc string `arg:"env:MAILING_BIND_GRPC"`
}

func main() {
	arg.MustParse(&args)

	if args.dbPath == "" {
		args.dbPath = "mail_list.db"
	}

	if args.BindJson == "" {
		args.BindJson = ":8080"
	}

	if args.BindGrpc == "" {
		args.BindJson = ":8081"
	}

	log.Printf("Using database: %v\n", args.dbPath)

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

	wg.Add(1)

	go func() {

		log.Printf("Starting grpcapi server...\n")
		grpcapi.Serve(db, args.BindJson)

		wg.Done()
	}()

	wg.Wait()

}
