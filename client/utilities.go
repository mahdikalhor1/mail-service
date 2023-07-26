package main

import (
	"context"
	"log"
	"mailService/mailDB"
	"mailService/proto"
	"time"
)

func displayResponse(res *proto.GetEmailResponse, err error){
	if err != nil{
		log.Fatalf("\terror : %v\n", err)
	}

	if res.EmailEntry == nil{
		log.Printf("\temail response not found!\n")
	}else{
		log.Println(res.EmailEntry)
	}

	log.Println()
}

func displayEmailList(res *proto.GetEmailBatchResponse, err error){
	if err != nil{
		log.Fatalf("\terror : %v\n", err)
	}

	log.Println("emails:")

	for _, email := range res.EmailEntry{
		log.Println(email)
	}

	log.Println()
}

func insertEmail(client proto.MailingServiceClient, email string){
	log.Println("Insert email...")

	context, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	req := proto.InsertEmailRequest{Email: email}

	res, err := client.InsertEmail(context, &req)

	displayResponse(res, err)
}

func getEmail(client proto.MailingServiceClient, email string){
	log.Println("Get email...")

	context, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	req := proto.GetEmailRequest{Email: email}

	res, err := client.GetEmail(context, &req)

	displayResponse(res, err)
}

func deleteEmail(client proto.MailingServiceClient, email string){
	log.Println("Delete email...")

	context, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	req := proto.DeleteEmailRequest{Email: email}

	res, err := client.DeleteEmail(context, &req)

	displayResponse(res, err)
}

func updateEmail(client proto.MailingServiceClient, entry mailDB.EmailEntry){
	log.Println("Update email...")

	context, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()


	req := proto.UpdateEmailRequest{
		EmailEntry: &proto.EmailEntry{Email: entry.Email, Id: entry.Id,
			 ConfirmedAt: entry.ConfirmedAt.Unix(), Optout: entry.OptOut,}}

	res, err := client.UpdateEmail(context, &req)

	displayResponse(res, err)
}

func getEmailBatch(client proto.MailingServiceClient, page, count int){
	log.Println("Get email batch email...")

	context, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()


	req := proto.GetEmailBatchRequest{Page: int64(page), Count: int64(count)}
	res, err := client.GetEmailBatch(context, &req)

	displayEmailList(res, err)
}