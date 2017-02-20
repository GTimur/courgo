//Реализует отправку почтовых сообщений
package courgo

import (
	"html/template"
	"net/smtp"
	"log"
	"net/mail"
)

//Структура учетных данных sendmail
type EmailCredentials struct {
	Username, Password, Server, From string
	Port                             int
	UseTLS                           bool
}

const emailTemplate = `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}

{{.Body}}
`

var t *template.Template

func init() {
	t = template.New("email")
	t.Parse(emailTemplate)
}

func SendEmailMsg(authCreds EmailCredentials, msg *Message) error {

	// compose the message
	msg.Subject = "Тестовое сообщение."
	msg.Body = "Это сообщение было написано в пятницу."
	msg.From = mail.Address{Name: "COURIER GO", Address: authCreds.From}
	msg.To = []string{"to-timur@yandex.ru"}

	// add attachments
	if err := msg.Attach("Вложение1.jpg"); err != nil {
		log.Fatal(err)
	}
	if err := msg.Attach("Вложение2.pdf"); err != nil {
		log.Fatal(err)
	}
	if err := msg.Attach("Вложение3.docx"); err != nil {
		log.Fatal(err)
	}


	//Отправляем почту
	auth := smtp.PlainAuth("", authCreds.Username, authCreds.Password, authCreds.Server)

	//Отправка без TLS
	if !authCreds.UseTLS {
		if err := SendMail(authCreds.Server, uint(authCreds.Port), auth, msg); err != nil {
			log.Println(err)
			return err
		}
		return nil
	}
	//Отправка с TLS
	if err := SendMailSSL(authCreds.Server, uint(authCreds.Port), auth, msg); err != nil {
		log.Println(err)
		return err
	}
	return nil
}