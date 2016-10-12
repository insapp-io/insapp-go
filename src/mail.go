package main

import (
	"net/smtp"
)


func SendEmail(to string, subject string, body string) {
	config, _ := Configuration()
  from := config.Email
	pass := config.Password
	cc := "insapp.contact@gmail.com"
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
    "Cc: " + cc + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))
}
