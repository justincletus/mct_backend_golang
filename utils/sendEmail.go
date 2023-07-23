package utils

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"strconv"

	"github.com/justincletus/cms/config"
	"github.com/k3a/html2text"
	gomail "gopkg.in/mail.v2"
)

var temp *template.Template

func init() {
	temp = template.Must(template.ParseGlob("templates/*.html"))
}

type EmailBody struct {
	Code     string
	Message  string
	Url      string
	Id       string
	Status   string
	Email    string
	Password string
	Fullname string
}

// SendEmail function receive from email address and html file
func (e_msg EmailBody) SendEmail(email string, subject string, e_temp string, options ...string) error {
	emailConfig, err := config.Config()
	ccEmail := ""
	if err != nil {
		return fmt.Errorf("email configuration values not set %v", err.Error())
	}

	if len(options) > 0 {
		ccEmail = options[0]
	}

	from := emailConfig["smtp_user"]
	password := emailConfig["smtp_password"]
	smtpHost := emailConfig["smtp_host"]
	smtpPort := emailConfig["smtp_port"]

	port, _ := strconv.Atoi(smtpPort)

	var body bytes.Buffer

	data := &e_msg

	if err := temp.ExecuteTemplate(&body, e_temp, &data); err != nil {
		fmt.Println(err)
		return err
	}

	// return nil
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", email)
	if ccEmail != "" {
		m.SetHeader("Cc", ccEmail)
	}
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))
	d := gomail.NewDialer(smtpHost, port, from, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}
