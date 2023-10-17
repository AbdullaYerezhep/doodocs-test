package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"mime/multipart"
	"strconv"
	"test/config"
	"gopkg.in/gomail.v2"
)

type EmailService interface {
	SendFileByEmail(file *multipart.FileHeader) (bool, error)
	Address() string
}

func NewEmailService()	*Email {
	return &Email{}
}

type Email struct {}

func (e *Email) SendFileByEmail(emails []string, file *multipart.FileHeader, cfg *config.Config) (bool, error) {
	from := cfg.Gmail.Login
	password := cfg.Gmail.Password    
	to := append([]string{}, emails...)   
	host := cfg.SmtpServer.Host
	port, _ := strconv.Atoi(cfg.SmtpServer.Port)

	uploadedFile, err := file.Open()

	if err != nil {
		return false, err
	}

	defer uploadedFile.Close()
	if err != nil {
		return false, err
	}

	// Read the file data.
	fileData, err := io.ReadAll(uploadedFile)
	if err != nil {
		return false, err
	}
	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to...)
	msg.SetHeader("Subject", "Secret file")
	msg.SetBody("text/plain", "Do not respond! this is automated generated email")

	msg.Attach(file.Filename, gomail.SetCopyFunc(func(w io.Writer) error {
		_, err := io.Copy(w, bytes.NewReader(fileData))
		return err
	}))

	dialer := gomail.NewDialer(host, port, from, password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := dialer.DialAndSend(msg); err != nil {
		return false, err
	}

	fmt.Println("Email Sent!")
	return true, nil
}	
		


