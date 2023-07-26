package grpcapi

import (
	"context"
	"database/sql"
	"log"
	"mailService/mailDB"
	"mailService/proto"
	"net"
	"time"

	"google.golang.org/grpc"
)


type MailServer struct{
	proto.UnimplementedMailingServiceServer
	db *sql.DB
}


func pbEntryToMdbEntry(entry *proto.EmailEntry)(mailDB.EmailEntry){

	time := time.Unix(entry.ConfirmedAt, 0)

	return mailDB.EmailEntry{
		Id : entry.Id,
		Email: entry.Email,
		ConfirmedAt: &time,
		OptOut: entry.Optout,
	}
}

func mdbEntryToPb(entry *mailDB.EmailEntry)(proto.EmailEntry){

	return proto.EmailEntry{
		Id : entry.Id,
		Email: entry.Email,
		ConfirmedAt: entry.ConfirmedAt.Unix(),
		Optout: entry.OptOut,
	}
}

func returnEmailResponse(db *sql.DB, email string)(*proto.GetEmailResponse, error){
	
	emailEntry, err := mailDB.GetEmail(db,email)

	if err != nil{
		return &proto.GetEmailResponse{}, err
	}
	if emailEntry == nil {
		return &proto.GetEmailResponse{}, nil
	}

	pbEntry := mdbEntryToPb(emailEntry)

	return &proto.GetEmailResponse{EmailEntry: &pbEntry}, nil
}

func(server *MailServer) InsertEmail(context context.Context, req proto.InsertEmailRequest) (*proto.GetEmailResponse, error){
	log.Printf("Grpc insert email: %v\n", req)

	err := mailDB.InsertEmail(server.db, req.Email)

	if err != nil{

		return &proto.GetEmailResponse{}, err
	}

	return returnEmailResponse(server.db, req.Email)
}

func(server *MailServer) UpdateEmail(context context.Context, req proto.UpdateEmailRequest) (*proto.GetEmailResponse, error){
	log.Printf("Grpc update email: %v\n", req)

	emailEntry := pbEntryToMdbEntry(req.EmailEntry)
	err := mailDB.UpdateEmail(server.db, &emailEntry)

	if err != nil{

		return &proto.GetEmailResponse{}, err
	}

	return returnEmailResponse(server.db, emailEntry.Email)
}

func(server *MailServer) DeleteEmail(context context.Context, req proto.DeleteEmailRequest) (*proto.GetEmailResponse, error){
	log.Printf("Grpc delete email: %v\n", req)

	err := mailDB.DeleteEmail(server.db, req.Email)

	if err != nil{

		return &proto.GetEmailResponse{}, err
	}

	return returnEmailResponse(server.db, req.Email)
}

func(server *MailServer) GetEmail(context context.Context, req proto.GetEmailRequest) (*proto.GetEmailResponse, error){
	log.Printf("Grpc insert email: %v\n", req)

	return returnEmailResponse(server.db, req.Email)
}

func(server *MailServer) GetEmailBatch(context context.Context, req proto.GetEmailBatchRequest)(*proto.GetEmailBatchResponse, error){

	params := mailDB.GetEmailBatchParams{Count: int(req.Count), Page: int(req.Page)}

	emails, err := mailDB.GetEmailBatch(server.db, params)

	if err != nil{
		return &proto.GetEmailBatchResponse{}, err
	}

	emailEntries := make([]*proto.EmailEntry, 0)

	for _, email := range emails{
		entry := mdbEntryToPb(&email)
		emailEntries = append(emailEntries, &entry)
	}

	return &proto.GetEmailBatchResponse{EmailEntry: emailEntries}, nil
}

func Serve(db *sql.DB, bind string){
	listener, err := net.Listen("tcp", bind)

	if err != nil{
		log.Println("Grpc server error; failed to bind!", err)

	}

	grpcServer := grpc.NewServer()

	mailServer := MailServer{db: db}

	proto.RegisterMailingServiceServer(
		grpcServer, &mailServer)

	log.Printf("Grpc server is running on port: %v...\n", bind)

	if err:= grpcServer.Serve(listener); err != nil{
		log.Printf("Grpc server error: %v\n", err)
	}

}