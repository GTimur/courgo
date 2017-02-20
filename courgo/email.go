//Реализует отправку почтовых сообщений с вложениями
//часть кода позаимствована у https://github.com/scorredoira
package courgo

import (
	"net/mail"
	"path/filepath"
	"io/ioutil"
	"fmt"
	"mime"
	"encoding/base64"
	"bytes"
	"time"
	"strings"
	"net/smtp"
	"log"
	"crypto/tls"
	"strconv"
)

// Attachment represents an email attachment.
type Attachment struct {
	Filename string
	Data     []byte
	Inline   bool
}

// Message represents a smtp message.
type Message struct {
	From            mail.Address
	To              []string
	Cc              []string
	Bcc             []string
	ReplyTo         string
	Subject         string
	Body            string
	BodyContentType string
	Attachments     map[string]*Attachment
}

func (m *Message) attach(file string, inline bool) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	_, filename := filepath.Split(file)

	m.Attachments[filename] = &Attachment{
		Filename: filename,
		Data:     data,
		Inline:   inline,
	}

	return nil
}

// AttachBuffer attaches a binary attachment.
func (m *Message) AttachBuffer(filename string, buf []byte, inline bool) error {
	m.Attachments[filename] = &Attachment{
		Filename: filename,
		Data:     buf,
		Inline:   inline,
	}
	return nil
}

// Attach attaches a file.
func (m *Message) Attach(file string) error {
	return m.attach(file, false)
}

// Inline includes a file as an inline attachment.
func (m *Message) Inline(file string) error {
	return m.attach(file, true)
}

func newMessage(subject string, body string, bodyContentType string) *Message {
	m := &Message{Subject: subject, Body: body, BodyContentType: bodyContentType}

	m.Attachments = make(map[string]*Attachment)

	return m
}

// NewMessage returns a new Message that can compose an email with attachments
func NewMessage(subject string, body string) *Message {
	return newMessage(subject, body, "text/plain")
}

// NewHTMLMessage returns a new Message that can compose an HTML email with attachments
func NewHTMLMessage(subject string, body string) *Message {
	return newMessage(subject, body, "text/html")
}

// Bytes returns the mail data
func (m *Message) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	buf.WriteString("From: " + m.From.String() + "\r\n")

	t := time.Now()
	buf.WriteString("Date: " + t.Format(time.RFC822) + "\r\n")

	buf.WriteString("To: " + strings.Join(m.To, ",") + "\r\n")
	if len(m.Cc) > 0 {
		buf.WriteString("Cc: " + strings.Join(m.Cc, ",") + "\r\n")
	}

	//fix  Encode
	var coder = base64.StdEncoding
	var subject = "=?UTF-8?B?" + coder.EncodeToString([]byte(m.Subject)) + "?="
	buf.WriteString("Subject: " + subject + "\r\n")

	if len(m.ReplyTo) > 0 {
		buf.WriteString("Reply-To: " + m.ReplyTo + "\r\n")
	}

	buf.WriteString("MIME-Version: 1.0\r\n")

	boundary := "8e1f7336599caf92131b646ea6b8f5b9"  //разделитель вложений в письме (может быть любой, я использовал md5)

	if len(m.Attachments) > 0 {
		buf.WriteString("Content-Type: multipart/mixed; boundary=" + boundary + "\r\n")
		buf.WriteString("\r\n--" + boundary + "\r\n")
	}

	buf.WriteString(fmt.Sprintf("Content-Type: %s; charset=utf-8\r\n\r\n", m.BodyContentType))
	buf.WriteString(m.Body)
	buf.WriteString("\r\n")

	if len(m.Attachments) > 0 {
		for _, attachment := range m.Attachments {
			buf.WriteString("\r\n\r\n--" + boundary + "\r\n")

			if attachment.Inline {
				buf.WriteString("Content-Type: message/rfc822\r\n")
				buf.WriteString("Content-Disposition: inline; filename=\"" + attachment.Filename + "\"\r\n\r\n")

				buf.Write(attachment.Data)
			} else {
				ext := filepath.Ext(attachment.Filename)
				mimetype := mime.TypeByExtension(ext)
				if mimetype != "" {
					mime := fmt.Sprintf("Content-Type: %s;\r\n", mimetype)
					buf.WriteString(mime)
					//gtg:
					buf.WriteString("\tname=\"" + attachment.Filename + "\"\r\n")
				} else {
					buf.WriteString("Content-Type: application/octet-stream;\r\n")
					//gtg:
					buf.WriteString("\tname=\"" + attachment.Filename + "\"\r\n")
				}
				buf.WriteString("Content-Transfer-Encoding: base64\r\n")
				buf.WriteString("Content-Disposition: attachment;\r\n\tfilename=\""+attachment.Filename+"\"\r\n\r\n")

				//buf.WriteString("Content-Disposition: attachment;\r\nfilename=\"=?UTF-8?B?")
				//buf.WriteString(coder.EncodeToString([]byte(attachment.Filename)))
				//buf.WriteString("?=\"\r\n\r\n")

				b := make([]byte, base64.StdEncoding.EncodedLen(len(attachment.Data)))
				base64.StdEncoding.Encode(b, attachment.Data)

				// write base64 content in lines of up to 76 chars
				for i, l := 0, len(b); i < l; i++ {
					buf.WriteByte(b[i])
					if (i + 1) % 76 == 0 {
						buf.WriteString("\r\n")
					}
				}
			}
			//buf.WriteString("\r\n--" + boundary)
		}
		buf.WriteString("\r\n--" + boundary)
		buf.WriteString("--")
	}

	return buf.Bytes()
}

// Tolist returns all the recipients of the email
func (m *Message) Tolist() []string {
	tolist := m.To

	for _, cc := range m.Cc {
		tolist = append(tolist, cc)
	}

	for _, bcc := range m.Bcc {
		tolist = append(tolist, bcc)
	}

	return tolist
}

// Send sends the message.
func SendMail(addr string, port uint, auth smtp.Auth, m *Message) error {
	return smtp.SendMail(addr+":"+strconv.Itoa(int(port)), auth, m.From.Address, m.Tolist(), m.Bytes())
}

// Send sends the message.
func SendMailSSL(addr string, port uint, auth smtp.Auth, m *Message) error {
	return SendMailTLS(addr, port, auth, m.From.Address, m.Tolist(), m.Bytes())
}

/*Реализация аналога smtp.SendMail для работы с SSL (yandex, google smtp)*/
func SendMailTLS(host string, port uint, auth smtp.Auth, from string,to []string, msg []byte) (err error) {
	serverAddr := fmt.Sprintf("%s:%d", host, port)

	conn, err := tls.Dial("tcp", serverAddr, nil)
	if err != nil {
		log.Println("Error Dialing", serverAddr)
		log.Println("Error Dialing", err)
		return err
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Println("Error SMTP connection", err)
		return err
	}

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				log.Println("Error during AUTH", err)
				return err
			}
		}
	}

	if err = c.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return c.Quit()
}