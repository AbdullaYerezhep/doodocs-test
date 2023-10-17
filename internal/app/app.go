package app

import (
	"fmt"
	"log"
	"net/http"
	"test/config"
	"test/internal/handler"
	"test/internal/service"
	"test/internal/service/archive"
	"test/internal/service/email"
)

func Run(addr string) {
	config := config.MustLoad()
	//int archive service 
	archive := archive.NewArchiveService()
	//init email service 
	email := email.NewEmailService()
	//init service
	service := service.NewService(archive, email, config)
	//init handler
	handler := handler.NewHandler(service)
	//init router
	router := handler.Router()
	//start server
	fmt.Println("Server is running on port :8080!")
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Something went wrong, could not start server: %s", err)
	}

}
