package main

import (
	"bufio"
	"fmt"
	"mailService/mailDB"
	"mailService/proto"
	"os"
	"strconv"
)


func runMenu(client proto.MailingServiceClient){

	fmt.Println("Enter a command(enter help).")

	scanner := bufio.NewScanner(os.Stdin)
	 
	for scanner.Scan(){

		switch(scanner.Text()){
		case "help":
			help()
		case "insert":
			fmt.Println("Eneter email Address")
			scanner.Scan()
			insertEmail(client, scanner.Text())
		case "delete":
			fmt.Println("Eneter email Address")
			scanner.Scan()
			deleteEmail(client, scanner.Text())
		case "get":
			fmt.Println("Eneter email Address")
			scanner.Scan()
			getEmail(client, scanner.Text())
		case "update":
			fmt.Println("Eneter email Address:")
			scanner.Scan()
			email := scanner.Text()
			fmt.Println("OptOut(true,false):")
			scanner.Scan()
			optOut, err := strconv.ParseBool(scanner.Text())
			if err != nil{
				fmt.Println("Invalid input!")
				break
			}
			
			entry := mailDB.EmailEntry{Email: email, OptOut: optOut}
			updateEmail(client, entry)
		
		case "getbatch":
			fmt.Println("Enter page number:")
			scanner.Scan()
			page, err := strconv.Atoi(scanner.Text())

			if err != nil{
				fmt.Println("Invalid input!")
				break
			}
			fmt.Println("Enter number of emails in each page:")
			scanner.Scan()
			count, err := strconv.Atoi(scanner.Text())
			if err != nil{
				fmt.Println("Invalid input!")
				break
			}
			
			getEmailBatch(client, page, count)
		case "exit":
			fmt.Println("Finished.")
			return
		default:
			fmt.Println("Invalid input!")
		}

	}
}

func help(){
	fmt.Printf("\tCommands:\n\tinsert\tdelete\tget\tupdate\tgetbatch\texit\n")
}