package email

import (
	"bytes"
	"emailnotifl3n/app/config"
	"emailnotifl3n/features/user"
	"fmt"
	"net/smtp"
	"strconv"
	"text/template"
)

type EmailData struct {
	URL     string
	Name    string
	Subject string
}

type EmailInterface interface {
	SendResetPasswordEmail(user *user.Core, token string) error
}

type emailService struct {
	url      string
	from     string
	password string
	host     string
	port     string
	user     string
}

func New() EmailInterface {
	cfg := config.InitConfig()
	emailCfg := cfg
	return &emailService{
		url:      emailCfg.PASSWD_URL,
		from:     emailCfg.EMAIL_FROM,
		user:     emailCfg.SMTP_USER,
		password: emailCfg.SMTP_PASS,
		host:     emailCfg.SMTP_HOST,
		port:     strconv.Itoa(emailCfg.SMTP_PORT),
	}
}

func (e *emailService) SendResetPasswordEmail(user *user.Core, token string) error {
	to := []string{user.Email}

	// Load the email template.
	t, err := template.ParseGlob("utils/templates/*.html")
	if err != nil {
		return err
	}

	// Prepare the email data.
	data := &EmailData{
		URL:     e.url + "?token=" + token,
		Name:    user.Name,
		Subject: "Reset Password",
	}

	// Prepare the email body by applying the data to the template.
	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return err
	}

	// Message.
	message := []byte(fmt.Sprintf("Subject: %s\n%s\n\n%s", data.Subject, "To: "+user.Name, body.String()))

	// Authentication.
	auth := smtp.PlainAuth("", e.user, e.password, e.host)

	// Sending email.
	err = smtp.SendMail(e.host+":"+e.port, auth, e.from, to, message)
	if err != nil {
		return err
	}

	return nil
}
