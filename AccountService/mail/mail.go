package mail

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendVerificationEmail(to, code string) error {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM")

	auth := smtp.PlainAuth("", user, password, host)
	subject := "Subject: Verification Code \r \n"
	body := fmt.Sprintf("\r \nYour verification code is: %s\nIt expires in 5 minutes.", code)

	msg := []byte(subject + body)
	addr := host + ":" + port
	return smtp.SendMail(addr, auth, from, []string{to}, msg)
}
