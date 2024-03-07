package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/rm4n0s/collabeek/hr/cases/server"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error: loading .env file:", err)
	}

	errs := []error{}
	dbHost := os.Getenv("HR_COLLABEEK_DB")
	if len(dbHost) == 0 {
		errs = append(errs, errors.New("HR_COLLABEEK_DB is missing"))
	}
	portServerStr := os.Getenv("HR_COLLABEEK_PORT")
	if len(portServerStr) == 0 {
		errs = append(errs, errors.New("HR_COLLABEEK_PORT is missing"))
	}
	portServer, err := strconv.Atoi(portServerStr)
	if err != nil {
		log.Fatal("HR_COLLABEEK_PORT is not correct")
	}
	tokenSecret := os.Getenv("HR_COLLABEEK_SECRET")
	if len(tokenSecret) == 0 {
		errs = append(errs, errors.New("HR_COLLABEEK_SECRET is missing"))
	}
	usernameSmtp := os.Getenv("HR_COLLABEEK_SMTP_USERNAME")
	if len(usernameSmtp) == 0 {
		errs = append(errs, errors.New("HR_COLLABEEK_SMTP_USERNAME is missing"))
	}
	passwordSmtp := os.Getenv("HR_COLLABEEK_SMTP_PASSWORD")
	if len(passwordSmtp) == 0 {
		errs = append(errs, errors.New("HR_COLLABEEK_SMTP_PASSWORD is missing"))
	}
	hostSmtp := os.Getenv("HR_COLLABEEK_SMTP_HOST")
	if len(hostSmtp) == 0 {
		errs = append(errs, errors.New("HR_COLLABEEK_SMTP_HOST is missing"))
	}
	portSmtpStr := os.Getenv("HR_COLLABEEK_SMTP_PORT")
	if len(portSmtpStr) == 0 {
		errs = append(errs, errors.New("HR_COLLABEEK_SMTP_PORT is missing"))
	}
	portSmtp, err := strconv.Atoi(portSmtpStr)
	if err != nil {
		log.Fatal("HR_COLLABEEK_SMTP_PORT is not correct")
	}
	senderEmail := os.Getenv("HR_COLLABEEK_SENDER_EMAIL")
	if len(senderEmail) == 0 {
		errs = append(errs, errors.New("HR_COLLABEEK_SENDER_EMAIL is missing"))
	}
	adminEmail := os.Getenv("HR_COLLABEEK_ADMIN_EMAIL")
	if len(adminEmail) == 0 {
		errs = append(errs, errors.New("HR_COLLABEEK_ADMIN_EMAIL is missing"))
	}

	if len(errs) > 0 {
		fmt.Println("Errors")
		for _, v := range errs {
			fmt.Println(v.Error())
		}
		os.Exit(1)
	}

	srv := server.ServerInput{
		DBHost:       dbHost,
		TokenSecret:  tokenSecret,
		SmtpUsername: usernameSmtp,
		SmtpPassword: passwordSmtp,
		SmtpHost:     hostSmtp,
		SmtpPort:     portSmtp,
		SenderEmail:  senderEmail,
		AdminEmail:   adminEmail,
	}
	hrEcho, err := server.NewHrCollabeekEchoServer(srv)
	if err != nil {
		log.Fatal("Error: failed to initialize server", err)
	}
	hrEcho.Logger.Fatal(hrEcho.Start(fmt.Sprintf(":%d", portServer)))
}
