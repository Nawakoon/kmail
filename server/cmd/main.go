package main

import (
	"log"
	"net/http"
	handler "passwordless-mail-server/pkg/api"
	"passwordless-mail-server/pkg/auth"
	"passwordless-mail-server/pkg/mail"
	"passwordless-mail-server/pkg/util"

	_ "github.com/lib/pq"
)

func main() {
	const PORT = ":8080"

	database := util.ConnectDatabase()
	defer database.Close()

	// service factory
	mailStore := mail.NewStore(database)
	uuidStore := auth.NewUUIDStore(database)
	mailService := mail.NewService(mailStore, uuidStore)
	mailHandler := handler.NewHandler(mailService)

	// routes
	http.HandleFunc("/health", mailHandler.HealthCheck)
	http.HandleFunc("/mail/inbox", mailHandler.GetInbox)
	http.HandleFunc("/mail/inbox/", mailHandler.GetMail)
	http.HandleFunc("/mail/send", mailHandler.SendMail)

	log.Printf("Server is running on port %s\n", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}
