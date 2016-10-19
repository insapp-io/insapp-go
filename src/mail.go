package main

import (
	"net/smtp"
)


func SendEmail(to string, subject string, body string) {
	config, _ := Configuration()
  from := config.Email
	pass := config.Password
	cc := config.Email

	msg := "From: " + from + "\n" +
		"To: " + from + "\n" +
    "Cc: " + cc + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{from}, []byte(msg))
}
