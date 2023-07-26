package main

import (
	"log"
	"mailService/proto"

	"github.com/alexflint/go-arg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var args struct{
	grpcAddress string `arg:"env:MAILING_SERVICE_GRPCADDRESS"`
}
func main(){

	arg.MustParse(&args)

	if args.grpcAddress == ""{
		args.grpcAddress = ":8081"
	}

	connection, err := grpc.Dial(args.grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil{
		log.Fatalf("Connection error: %v\n", err)
	}

	defer connection.Close()

	client := proto.NewMailingServiceClient(connection)

	runMenu(client)
}