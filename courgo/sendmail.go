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
	Username, Password, Server string
	Port                       int
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

func SendEmailMsg() {
	/*authCreds := &EmailCredentials{
		Username: "noti@ymkbank.ru",
		Password: "Bank9991",
		Server: "10.20.20.6",
		Port: 25,
	}*/

	authCreds := &EmailCredentials{
		Username: "to-timur@yandex.ru",
		Password: "blank",
		Server: "smtp.yandex.ru",
		Port: 465,
	}

	// compose the message
	m := NewHTMLMessage("Тестовое сообщение.", "Это сообщение было написано в пятницу.")
	m.From = mail.Address{Name: "COURIER GO", Address: "to-timur@yandex.ru"}
	//m.To = []string{"gtg@ymkbank.ru"}
	m.To = []string{"to-timur@yandex.ru"}

	// add attachments
	if err := m.Attach("Вложение1.jpg"); err != nil {
		log.Fatal(err)
	}
	if err := m.Attach("Вложение2.pdf"); err != nil {
		log.Fatal(err)
	}
	if err := m.Attach("Вложение3.docx"); err != nil {
		log.Fatal(err)
	}

	// send it
	auth := smtp.PlainAuth("",authCreds.Username,authCreds.Password,authCreds.Server)

	//if err := SendMail(authCreds.Server+":"+strconv.Itoa(authCreds.Port), auth, m); err != nil {
	//	log.Fatal(err)
	//}
	if err := SendMailSSL(authCreds.Server,uint(authCreds.Port), auth, m); err != nil {
		log.Fatal(err)
	}
}