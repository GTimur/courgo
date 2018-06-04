/*
 	Реализует отправку почтовых сообщений с вложениями.
 */

package courgo

import (
	"log"
	"net/mail"
        "github.com/Gtimur/courgo/smtp"
)

//Структура учетных данных sendmail
type EmailCredentials struct {
	Username, Password, Server, From, FromName string
	Port                             int
	UseTLS                           bool
}

func SendEmailMsg(authCreds EmailCredentials, msg *Message) error {

	//Зафиксируем сведения об отправителе
	msg.From = mail.Address{Name: authCreds.FromName, Address: authCreds.From}

	//Отправляем почту
	auth := smtp.PlainAuth("", authCreds.Username, authCreds.Password, authCreds.Server)

	//Отправка без TLS
	if !authCreds.UseTLS {
		if err := SendMail(authCreds.Server, uint(authCreds.Port), auth, msg); err != nil {
			log.Println("SendEmailMsg error:",err)
			return err
		}
		return nil
	}
	//Отправка с TLS
	if err := SendMailSSL(authCreds.Server, uint(authCreds.Port), auth, msg); err != nil {
		log.Println("SendEmailMsgSSL error:",err)
		return err
	}
	return nil
}