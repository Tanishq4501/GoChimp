package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"sync"
	"time"
)

type EmailData struct {
	Name    string
	Email   string
	Subject string
	Message string
}

func executeTemplate(templatePath string, data EmailData) (string, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		return "", err
	}

	return body.String(), nil
}

func emailWorker(id int, ch chan Recipient, wg *sync.WaitGroup) {

	defer wg.Done()

	for recipient := range ch {

		time.Sleep(200 * time.Millisecond)

		smptHost := "localhost"
		smptPort := "1025"

		// Prepare email data
		emailData := EmailData{
			Name:    recipient.Name,
			Email:   recipient.Email,
			Subject: "Welcome to Our Email Campaign",
			Message: "We're excited to have you join our community!",
		}

		// Execute template
		body, err := executeTemplate("./email_template.html", emailData)
		if err != nil {
			log.Printf("Worker %d: Failed to execute template: %v\n", id, err)
			continue
		}

		// Format email message with headers
		formattedMsg := fmt.Sprintf(
			"From: %s\r\n"+
				"To: %s\r\n"+
				"Subject: %s\r\n"+
				"MIME-Version: 1.0\r\n"+
				"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
				"\r\n"+
				"%s",
			"tanishqchoudhary5689@gmail.com",
			recipient.Email,
			emailData.Subject,
			body,
		)
		msg := []byte(formattedMsg)

		fmt.Printf("Worker %d: Sending email to %s \n", id, recipient.Email)

		err = smtp.SendMail(smptHost+":"+smptPort, nil, "tanishqchoudhary5689@gmail.com", []string{recipient.Email}, msg)

		if err != nil {
			log.Printf("Worker %d: Failed to send email to %s: %v\n", id, recipient.Email, err)
			continue
		}

		time.Sleep(50 * time.Millisecond)

		fmt.Printf("Worker %d: Sent email to %s \n", id, recipient.Email)

	}
}
